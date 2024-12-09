package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/app/repositories"
	"loki/internal/app/serializers"
	"loki/internal/config"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

func Test_Authentication_CreateSmartIdSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		SmartId: config.SmartId{
			BaseURL:          "https://sid.demo.sk.ee/smart-id-rp/v2",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
	}
	smartIdMock := NewMockSmartIdProvider(ctrl)
	smartIdQueue := make(chan *SmartIdQueue, 1)

	mobileIdMock := NewMockMobileIdProvider(ctrl)
	mobileIdQueue := make(chan *MobileIdQueue, 1)

	database := repositories.NewMockDatabase(ctrl)
	redis := repositories.NewMockRedis(ctrl)

	jwtMock := jwt.NewMockJwt(ctrl)
	log := logger.NewLogger()

	id, _ := uuid.Parse("8fdb516d-1a82-43ba-b82d-be63df569b86")
	sessionId := id.String()

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		database,
		redis,
		jwtMock,
		log)

	tests := []struct {
		name     string
		before   func()
		params   dto.CreateSmartIdSessionRequest
		expected *serializers.SessionSerializer
		error    error
	}{
		{
			name: "Success",
			before: func() {
				smartIdMock.EXPECT().CreateSession(ctx, dto.CreateSmartIdSessionRequest{
					Country:      "EE",
					PersonalCode: "30303039914",
				}).Return(&dto.SmartIdProviderSessionResponse{
					ID:   sessionId,
					Code: "1234",
				}, nil)

				redis.EXPECT().CreateSession(ctx, &models.Session{
					ID:           id,
					PersonalCode: "30303039914",
					Code:         "1234",
					Status:       "RUNNING",
				})
			},
			params: dto.CreateSmartIdSessionRequest{
				Country:      "EE",
				PersonalCode: "30303039914",
			},
			expected: &serializers.SessionSerializer{
				ID:   id,
				Code: "1234",
			},
			error: nil,
		},
		{
			name: "Error to create smart-id session",
			before: func() {
				smartIdMock.EXPECT().CreateSession(ctx, dto.CreateSmartIdSessionRequest{
					Country:      "EE",
					PersonalCode: "30303039914",
				}).Return(nil, assert.AnError)
			},
			params: dto.CreateSmartIdSessionRequest{
				Country:      "EE",
				PersonalCode: "30303039914",
			},
			expected: nil,
			error:    assert.AnError,
		},
		{
			name: "Error to save smart-id session",
			before: func() {
				smartIdMock.EXPECT().CreateSession(ctx, dto.CreateSmartIdSessionRequest{
					Country:      "EE",
					PersonalCode: "30303039914",
				}).Return(&dto.SmartIdProviderSessionResponse{
					ID:   sessionId,
					Code: "1234",
				}, nil)

				redis.EXPECT().CreateSession(ctx, &models.Session{
					ID:           id,
					PersonalCode: "30303039914",
					Code:         "1234",
					Status:       "RUNNING",
				}).Return(assert.AnError)
			},
			params: dto.CreateSmartIdSessionRequest{
				Country:      "EE",
				PersonalCode: "30303039914",
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.CreateSmartIdSession(ctx, tt.params)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Authentication_GetSmartIdSessionStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		SmartId: config.SmartId{
			BaseURL:          "https://sid.demo.sk.ee/smart-id-rp/v2",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
	}
	smartIdMock := NewMockSmartIdProvider(ctrl)
	smartIdQueue := make(chan *SmartIdQueue, 1)

	mobileIdMock := NewMockMobileIdProvider(ctrl)
	mobileIdQueue := make(chan *MobileIdQueue, 1)

	database := repositories.NewMockDatabase(ctrl)
	redis := repositories.NewMockRedis(ctrl)

	jwtMock := jwt.NewMockJwt(ctrl)
	log := logger.NewLogger()

	id, _ := uuid.Parse("8fdb516d-1a82-43ba-b82d-be63df569b86")

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		database,
		redis,
		jwtMock,
		log)

	tests := []struct {
		name     string
		before   func()
		id       uuid.UUID
		expected *dto.SmartIdProviderSessionStatusResponse
		error    error
	}{
		{
			name: "Success",
			before: func() {
				smartIdMock.EXPECT().GetSessionStatus(id).Return(&dto.SmartIdProviderSessionStatusResponse{
					State: "COMPLETE",
					Result: dto.SmartIdProviderResult{
						EndResult: "OK",
					},
					Signature: dto.SmartIdProviderSignature{
						Value:     "signature",
						Algorithm: "algorithm",
					},
					Cert: dto.SmartIdProviderCertificate{
						Value:            "certificate",
						CertificateLevel: "QUALIFIED",
					},
					InteractionFlowUsed: "displayTextAndPIN",
				}, nil)
			},
			id: id,
			expected: &dto.SmartIdProviderSessionStatusResponse{
				State: "COMPLETE",
				Result: dto.SmartIdProviderResult{
					EndResult: "OK",
				},
				Signature: dto.SmartIdProviderSignature{
					Value:     "signature",
					Algorithm: "algorithm",
				},
				Cert: dto.SmartIdProviderCertificate{
					Value:            "certificate",
					CertificateLevel: "QUALIFIED",
				},
				InteractionFlowUsed: "displayTextAndPIN",
			},
		},
		{
			name: "Error",
			before: func() {
				smartIdMock.EXPECT().GetSessionStatus(id).Return(nil, assert.AnError)
			},
			id:       id,
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.GetSmartIdSessionStatus(ctx, tt.id)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Authentication_CreateMobileIdSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		MobileId: config.MobileId{
			BaseURL:          "https://tsp.demo.sk.ee/mid-api",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
	}
	smartIdMock := NewMockSmartIdProvider(ctrl)
	smartIdQueue := make(chan *SmartIdQueue, 1)

	mobileIdMock := NewMockMobileIdProvider(ctrl)
	mobileIdQueue := make(chan *MobileIdQueue, 1)

	database := repositories.NewMockDatabase(ctrl)
	redis := repositories.NewMockRedis(ctrl)

	jwtMock := jwt.NewMockJwt(ctrl)
	log := logger.NewLogger()

	id, _ := uuid.Parse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")
	sessionId := id.String()

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		database,
		redis,
		jwtMock,
		log)

	tests := []struct {
		name     string
		before   func()
		params   dto.CreateMobileIdSessionRequest
		expected *serializers.SessionSerializer
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mobileIdMock.EXPECT().CreateSession(ctx, dto.CreateMobileIdSessionRequest{
					Locale:       "ENG",
					PhoneNumber:  "+37268000769",
					PersonalCode: "60001017869",
				}).Return(&dto.MobileIdProviderSessionResponse{
					ID:   sessionId,
					Code: "1234",
				}, nil)

				redis.EXPECT().CreateSession(ctx, &models.Session{
					ID:           id,
					PersonalCode: "60001017869",
					Code:         "1234",
					Status:       "RUNNING",
				})
			},
			params: dto.CreateMobileIdSessionRequest{
				Locale:       "ENG",
				PhoneNumber:  "+37268000769",
				PersonalCode: "60001017869",
			},
			expected: &serializers.SessionSerializer{
				ID:   id,
				Code: "1234",
			},
			error: nil,
		},
		{
			name: "Error to create mobile-id session",
			before: func() {
				mobileIdMock.EXPECT().CreateSession(ctx, dto.CreateMobileIdSessionRequest{
					Locale:       "ENG",
					PhoneNumber:  "+37268000769",
					PersonalCode: "60001017869",
				}).Return(nil, assert.AnError)
			},
			params: dto.CreateMobileIdSessionRequest{
				Locale:       "ENG",
				PhoneNumber:  "+37268000769",
				PersonalCode: "60001017869",
			},
			expected: nil,
			error:    assert.AnError,
		},
		{
			name: "Error to save mobile-id session",
			before: func() {
				mobileIdMock.EXPECT().CreateSession(ctx, dto.CreateMobileIdSessionRequest{
					Locale:       "ENG",
					PhoneNumber:  "+37268000769",
					PersonalCode: "60001017869",
				}).Return(&dto.MobileIdProviderSessionResponse{
					ID:   sessionId,
					Code: "1234",
				}, nil)

				redis.EXPECT().CreateSession(ctx, &models.Session{
					ID:           id,
					PersonalCode: "60001017869",
					Code:         "1234",
					Status:       "RUNNING",
				}).Return(assert.AnError)
			},
			params: dto.CreateMobileIdSessionRequest{
				Locale:       "ENG",
				PhoneNumber:  "+37268000769",
				PersonalCode: "60001017869",
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.CreateMobileIdSession(ctx, tt.params)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Authentication_GetMobileIdSessionStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		MobileId: config.MobileId{
			BaseURL:          "https://tsp.demo.sk.ee/mid-api",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
	}
	smartIdMock := NewMockSmartIdProvider(ctrl)
	smartIdQueue := make(chan *SmartIdQueue, 1)

	mobileIdMock := NewMockMobileIdProvider(ctrl)
	mobileIdQueue := make(chan *MobileIdQueue, 1)

	database := repositories.NewMockDatabase(ctrl)
	redis := repositories.NewMockRedis(ctrl)

	jwtMock := jwt.NewMockJwt(ctrl)
	log := logger.NewLogger()

	id, _ := uuid.Parse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		database,
		redis,
		jwtMock,
		log)

	tests := []struct {
		name     string
		before   func()
		id       uuid.UUID
		expected *dto.MobileIdProviderSessionStatusResponse
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mobileIdMock.EXPECT().GetSessionStatus(id).Return(&dto.MobileIdProviderSessionStatusResponse{
					State:  "COMPLETE",
					Result: "OK",
					Signature: dto.MobileIdProviderSignature{
						Value:     "signature",
						Algorithm: "algorithm",
					},
					Cert: "certificate",
				}, nil)
			},
			id: id,
			expected: &dto.MobileIdProviderSessionStatusResponse{
				State:  "COMPLETE",
				Result: "OK",
				Signature: dto.MobileIdProviderSignature{
					Value:     "signature",
					Algorithm: "algorithm",
				},
				Cert: "certificate",
			},
		},
		{
			name: "Error",
			before: func() {
				mobileIdMock.EXPECT().GetSessionStatus(id).Return(nil, assert.AnError)
			},
			id:       id,
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.GetMobileIdSessionStatus(ctx, tt.id)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Authentication_UpdateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		SmartId: config.SmartId{
			BaseURL:          "https://sid.demo.sk.ee/smart-id-rp/v2",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
		MobileId: config.MobileId{
			BaseURL:          "https://tsp.demo.sk.ee/mid-api",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
	}
	smartIdMock := NewMockSmartIdProvider(ctrl)
	smartIdQueue := make(chan *SmartIdQueue, 1)

	mobileIdMock := NewMockMobileIdProvider(ctrl)
	mobileIdQueue := make(chan *MobileIdQueue, 1)

	database := repositories.NewMockDatabase(ctrl)
	redis := repositories.NewMockRedis(ctrl)

	jwtMock := jwt.NewMockJwt(ctrl)
	log := logger.NewLogger()

	id, _ := uuid.Parse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		database,
		redis,
		jwtMock,
		log)

	tests := []struct {
		name     string
		before   func()
		params   models.Session
		expected *serializers.SessionSerializer
		error    error
	}{
		{
			name: "Success",
			before: func() {
				redis.EXPECT().UpdateSession(ctx, &models.Session{
					ID:     id,
					Status: "COMPLETE",
					Payload: models.SessionPayload{
						State:     "COMPLETE",
						Result:    "OK",
						Signature: "signature",
						Cert:      "certificate",
					},
				}).Return(nil)
			},
			params: models.Session{
				ID:     id,
				Status: "COMPLETE",
				Payload: models.SessionPayload{
					State:     "COMPLETE",
					Result:    "OK",
					Signature: "signature",
					Cert:      "certificate",
				},
			},
			expected: &serializers.SessionSerializer{
				ID:     id,
				Status: "COMPLETE",
			},
		},
		{
			name: "Error",
			before: func() {
				redis.EXPECT().UpdateSession(ctx, &models.Session{
					ID:     id,
					Status: "COMPLETE",
					Payload: models.SessionPayload{
						State:     "COMPLETE",
						Result:    "OK",
						Signature: "signature",
						Cert:      "certificate",
					},
				}).Return(assert.AnError)
			},
			params: models.Session{
				ID:     id,
				Status: "COMPLETE",
				Payload: models.SessionPayload{
					State:     "COMPLETE",
					Result:    "OK",
					Signature: "signature",
					Cert:      "certificate",
				},
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.UpdateSession(ctx, tt.params)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Authentication_FindSessionById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		SmartId: config.SmartId{
			BaseURL:          "https://sid.demo.sk.ee/smart-id-rp/v2",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
		MobileId: config.MobileId{
			BaseURL:          "https://tsp.demo.sk.ee/mid-api",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
	}
	smartIdMock := NewMockSmartIdProvider(ctrl)
	smartIdQueue := make(chan *SmartIdQueue, 1)

	mobileIdMock := NewMockMobileIdProvider(ctrl)
	mobileIdQueue := make(chan *MobileIdQueue, 1)

	database := repositories.NewMockDatabase(ctrl)
	redis := repositories.NewMockRedis(ctrl)

	jwtMock := jwt.NewMockJwt(ctrl)
	log := logger.NewLogger()

	id, _ := uuid.Parse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")
	sessionId := id.String()

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		database,
		redis,
		jwtMock,
		log)

	tests := []struct {
		name      string
		before    func()
		sessionId string
		expected  *serializers.SessionSerializer
		error     error
	}{
		{
			name: "Success",
			before: func() {
				redis.EXPECT().FindSessionById(ctx, id).Return(&models.Session{
					ID:     id,
					Status: "COMPLETE",
					Payload: models.SessionPayload{
						State:     "COMPLETE",
						Result:    "OK",
						Signature: "signature",
						Cert:      "certificate",
					},
				}, nil)
			},
			sessionId: sessionId,
			expected: &serializers.SessionSerializer{
				ID:     id,
				Status: "COMPLETE",
			},
		},
		{
			name: "Error",
			before: func() {
				redis.EXPECT().FindSessionById(ctx, id).Return(nil, assert.AnError)
			},
			sessionId: sessionId,
			expected:  nil,
			error:     assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindSessionById(ctx, tt.sessionId)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Authentication_Complete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		SmartId: config.SmartId{
			BaseURL:          "https://sid.demo.sk.ee/smart-id-rp/v2",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
		MobileId: config.MobileId{
			BaseURL:          "https://tsp.demo.sk.ee/mid-api",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
	}
	smartIdMock := NewMockSmartIdProvider(ctrl)
	smartIdQueue := make(chan *SmartIdQueue, 1)

	mobileIdMock := NewMockMobileIdProvider(ctrl)
	mobileIdQueue := make(chan *MobileIdQueue, 1)

	database := repositories.NewMockDatabase(ctrl)
	redis := repositories.NewMockRedis(ctrl)

	jwtMock := jwt.NewMockJwt(ctrl)
	log := logger.NewLogger()

	id, _ := uuid.Parse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")
	sessionId := id.String()

	userId, err := uuid.NewRandom()
	assert.NoError(t, err)

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		database,
		redis,
		jwtMock,
		log)

	tests := []struct {
		name      string
		before    func()
		sessionId string
		expected  *serializers.UserSerializer
		error     error
	}{
		{
			name: "Success (smart-id)",
			before: func() {
				redis.EXPECT().FindSessionById(ctx, id).Return(&models.Session{
					ID:     id,
					Status: "COMPLETE",
					Payload: models.SessionPayload{
						State:     "COMPLETE",
						Result:    "OK",
						Signature: "NUqfcuvZrQdJy9yAIpmDaOi22uWnCL7rbF+5Vx/th2pa7VK3+xRBM7S9CqEn9pRWsoOddqFDUfbz7w6Vd/AvaduLH3UEWs+LCTlQ+9liCGUcY4N97xMhlVwv1MnybBbDKKk7e+xXAdFGV7T+2lE5PwP9h4YyCl/1Jg1lXcuNWJEcu2E1bcJtOI6yDO+3PYEDuc/NNsj/1SZFvg+ffLhJOMKEOJe+Jxf6hsn6NoAFyBYBDvKGeAX92FHej5BFQbvk/sWJee9ENC+Mmjsr+rUiJI0iKh+WN0fiQYzdtv0TsowcGF0vqrRnlbDEc301xetowJBcefko8DcroqtvgzXQ3W0ruEeYKbehzEmB/iEI1iBjQi3hxrfaXD1cgZRzWcIurSzgv+rB5QE1xWV7GPRu9gV5b/yKRkfIbPclOa6OTpjlKTu+EG6qM7z1H9+UMp/Lx62Amin57W+oH0kiDm5zAMTETEDkRE0WXpKOlETvOUDS26hsXa9KlTnapisSpfSc0s2dCXjqYQ1Faw18gKKnZBhG5WkrqaaHpMGFqHVfPSVIe8uALVVaCBHzQ/Nly8dhE26YJXSY+BoIjrX4znXVCE38hpWPeYGMu+4Y/gm+STVkwQKVlXYIjG5nWpnkNI5ivvfHhusiLJf9bPKxMPSRtjme3g69vU4NHMpGnAJp2ER+i/S1DigULSEUQscrFNYFzu8Ha67a0SnRc1ozu2VCatyx+eFoVSuoriiVnOnb1/mXXsS3keBa2PygazHnMI2Hd7aKnhpM41fbKoa9awZbCy8Udw6gTfSyqLZNZ07fB1x9wEVnV7ZMZn8NcsBWYTR9v9DpZkYXO/7GgWTBCKtTwL0/TOWKmEkibPhmtGNqzA/+LeJoGCXgIvqBRZjVHLNZu90CtLCSddaSJR6MgMyfL4eGIWTybZULll8UO7G2XkRl46HsPVjv1CILmy80V6VHhRCBUEOpn9TWn6Q+hGelshbQozy4/hfidF9JkdXi0y/3fHBLyhsAcLsZ2+n6qoxs",
						Cert:      "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",
					},
				}, nil)

				jwtMock.EXPECT().Generate(jwt.Payload{ID: "PNOEE-30303039914"}, gomock.Any()).Return("access-token", nil)
				jwtMock.EXPECT().Generate(jwt.Payload{ID: "PNOEE-30303039914"}, gomock.Any()).Return("refresh-token", nil)

				database.EXPECT().CreateOrUpdateUserWithTokens(ctx, gomock.Any()).
					Return(&models.User{
						ID:             userId,
						IdentityNumber: "PNOEE-30303039914",
						PersonalCode:   "303039914",
						FirstName:      "TESTNUMBER",
						LastName:       "OK",
						AccessToken:    "access-token",
						RefreshToken:   "refresh-token",
					}, nil)

				redis.EXPECT().DeleteSessionByID(ctx, id).Return(nil)
			},
			sessionId: sessionId,
			expected: &serializers.UserSerializer{
				ID:             userId,
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "OK",
				AccessToken:    "access-token",
				RefreshToken:   "refresh-token",
			},
			error: nil,
		},
		{
			name: "Success (mobile-id)",
			before: func() {
				redis.EXPECT().FindSessionById(ctx, id).Return(&models.Session{
					ID:     id,
					Status: "COMPLETE",
					Payload: models.SessionPayload{
						State:     "COMPLETE",
						Result:    "OK",
						Signature: "tFib1+JlILD1ZBI3w4zUeyQLCh7aVcveelCNmgmYCzKn6ZdZ6keaZRiE0wlpBRjI+eHC1CHp4NKww5kp+M+Wsg==",
						Cert:      "MIIFXzCCA0egAwIBAgIQbZGUgkUm/YdiVSIaf13G+zANBgkqhkiG9w0BAQsFADBoMQswCQYDVQQGEwJFRTEiMCAGA1UECgwZQVMgU2VydGlmaXRzZWVyaW1pc2tlc2t1czEXMBUGA1UEYQwOTlRSRUUtMTA3NDcwMTMxHDAaBgNVBAMME1RFU1Qgb2YgRUlELVNLIDIwMTYwHhcNMjIwNDEyMDY1NDE4WhcNMjcwNDEyMjA1OTU5WjBtMQswCQYDVQQGEwJFRTEbMBkGA1UEAwwSRUlEMjAxNixURVNUTlVNQkVSMRMwEQYDVQQEDApURVNUTlVNQkVSMRAwDgYDVQQqDAdFSUQyMDE2MRowGAYDVQQFExFQTk9FRS02MDAwMTAxNzg2OTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABN5THR+i6stpLG0lq12yggjyLvyvu+0tUW2BF33CTrC019eG5oeeBYMVEm+ZBpdZlfwJUtlWXEuGw3XDdBgf14KjggHJMIIBxTAJBgNVHRMEAjAAMA4GA1UdDwEB/wQEAwIHgDB7BgNVHSAEdDByMGYGCSsGAQQBzh8SATBZMDcGCCsGAQUFBwIBFitodHRwczovL3NraWRzb2x1dGlvbnMuZXUvZW4vcmVwb3NpdG9yeS9DUFMvMB4GCCsGAQUFBwICMBIaEE9ubHkgZm9yIFRFU1RJTkcwCAYGBACPegECMB0GA1UdDgQWBBTTbr/pqvMoWlqZULY3pwCzp998ijAfBgNVHSMEGDAWgBSusOrhNvgmq6XMC2ZV/jodAr8StDB9BggrBgEFBQcBAQRxMG8wKQYIKwYBBQUHMAGGHWh0dHA6Ly9haWEuZGVtby5zay5lZS9laWQyMDE2MEIGCCsGAQUFBzAChjZodHRwczovL3NrLmVlL3VwbG9hZC9maWxlcy9URVNUX29mX0VJRC1TS18yMDE2LmRlci5jcnQwbAYIKwYBBQUHAQMEYDBeMFwGBgQAjkYBBTBSMFAWSmh0dHBzOi8vc2tpZHNvbHV0aW9ucy5ldS9lbi9yZXBvc2l0b3J5L2NvbmRpdGlvbnMtZm9yLXVzZS1vZi1jZXJ0aWZpY2F0ZXMvEwJFTjANBgkqhkiG9w0BAQsFAAOCAgEACfHa53mBsmnnnNlTa5DwXmI3R9tTjcNrjMa8alUdmvi50pipPpjvkYaCsiSJpUNNZ4EvfdRI1kKWbuCLc66MqQbZ2KNaOMZ8TODkx5uOhbOqGqWr0mBTJZJu7JboN3UB9/5lrpZlUuuhjAivpQRO1OqmQXMVfIRi/gy8sFc2l10gICZiBt2JBLmC7KiafEXK8WQJpEs8iKMHYLWubURcCrMvZXz2XdbSLfnjqS40P3uFhEJfeo6+aAsnD2M2AwBwmiF0P2/k8Vk/wc6miL8SZROJMMHqcvO9vQYlUihDSNn3Jz4fz0nzTZ4DtTSI6jGU0zy2SS2j7Srgdt+jDnIaANUEIHjsKGWqBWY77D3iNiSmN7OXLXhmsv10yjLxJSwBunFs+RvYQh1SiyvoM+/yq0SPrS/xzg0unRWAjmpGFMJ4Gw7Vn8+faJeTounOmrnZhG9LYrkLPsWxRUcxc+GfM11xUcwMExQvh7oRZ6iBibhKgLdYiP31Q4QEZpNZbjNIPVJtKEj4Z/zJlbySgb4TrXy+SgcmUilUPGowwo7wkSrw5wiG/lX9QxnQHbyUugY693BRH/twRk1jnYXie+p4wRUS4tNX9fb3hKSkIaWIMpryuHxl/WoG5/FmhDN6YpGqEaaf8/rQhUEAYae2dy5A4RgxkdtzQ2Q9uOz1qsw3T3g=",
					},
				}, nil)

				jwtMock.EXPECT().Generate(jwt.Payload{ID: "PNOEE-60001017869"}, gomock.Any()).Return("access-token", nil)
				jwtMock.EXPECT().Generate(jwt.Payload{ID: "PNOEE-60001017869"}, gomock.Any()).Return("refresh-token", nil)

				database.EXPECT().CreateOrUpdateUserWithTokens(ctx, gomock.Any()).
					Return(&models.User{
						ID:             userId,
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
						AccessToken:    "access-token",
						RefreshToken:   "refresh-token",
					}, nil)

				redis.EXPECT().DeleteSessionByID(ctx, id).Return(nil)
			},
			sessionId: sessionId,
			expected: &serializers.UserSerializer{
				ID:             userId,
				IdentityNumber: "PNOEE-60001017869",
				PersonalCode:   "60001017869",
				FirstName:      "EID2016",
				LastName:       "TESTNUMBER",
				AccessToken:    "access-token",
				RefreshToken:   "refresh-token",
			},
			error: nil,
		},
		{
			name: "Error: session not found",
			before: func() {
				redis.EXPECT().FindSessionById(ctx, id).Return(nil, assert.AnError)
			},
			sessionId: sessionId,
			expected:  nil,
			error:     assert.AnError,
		},
		{
			name: "Error: failed to delete session",
			before: func() {
				redis.EXPECT().FindSessionById(ctx, id).Return(&models.Session{
					ID:     id,
					Status: "COMPLETE",
					Payload: models.SessionPayload{
						State:     "COMPLETE",
						Result:    "OK",
						Signature: "NUqfcuvZrQdJy9yAIpmDaOi22uWnCL7rbF+5Vx/th2pa7VK3+xRBM7S9CqEn9pRWsoOddqFDUfbz7w6Vd/AvaduLH3UEWs+LCTlQ+9liCGUcY4N97xMhlVwv1MnybBbDKKk7e+xXAdFGV7T+2lE5PwP9h4YyCl/1Jg1lXcuNWJEcu2E1bcJtOI6yDO+3PYEDuc/NNsj/1SZFvg+ffLhJOMKEOJe+Jxf6hsn6NoAFyBYBDvKGeAX92FHej5BFQbvk/sWJee9ENC+Mmjsr+rUiJI0iKh+WN0fiQYzdtv0TsowcGF0vqrRnlbDEc301xetowJBcefko8DcroqtvgzXQ3W0ruEeYKbehzEmB/iEI1iBjQi3hxrfaXD1cgZRzWcIurSzgv+rB5QE1xWV7GPRu9gV5b/yKRkfIbPclOa6OTpjlKTu+EG6qM7z1H9+UMp/Lx62Amin57W+oH0kiDm5zAMTETEDkRE0WXpKOlETvOUDS26hsXa9KlTnapisSpfSc0s2dCXjqYQ1Faw18gKKnZBhG5WkrqaaHpMGFqHVfPSVIe8uALVVaCBHzQ/Nly8dhE26YJXSY+BoIjrX4znXVCE38hpWPeYGMu+4Y/gm+STVkwQKVlXYIjG5nWpnkNI5ivvfHhusiLJf9bPKxMPSRtjme3g69vU4NHMpGnAJp2ER+i/S1DigULSEUQscrFNYFzu8Ha67a0SnRc1ozu2VCatyx+eFoVSuoriiVnOnb1/mXXsS3keBa2PygazHnMI2Hd7aKnhpM41fbKoa9awZbCy8Udw6gTfSyqLZNZ07fB1x9wEVnV7ZMZn8NcsBWYTR9v9DpZkYXO/7GgWTBCKtTwL0/TOWKmEkibPhmtGNqzA/+LeJoGCXgIvqBRZjVHLNZu90CtLCSddaSJR6MgMyfL4eGIWTybZULll8UO7G2XkRl46HsPVjv1CILmy80V6VHhRCBUEOpn9TWn6Q+hGelshbQozy4/hfidF9JkdXi0y/3fHBLyhsAcLsZ2+n6qoxs",
						Cert:      "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",
					},
				}, nil)

				jwtMock.EXPECT().Generate(jwt.Payload{ID: "PNOEE-30303039914"}, gomock.Any()).Return("access-token", nil)
				jwtMock.EXPECT().Generate(jwt.Payload{ID: "PNOEE-30303039914"}, gomock.Any()).Return("refresh-token", nil)

				database.EXPECT().CreateOrUpdateUserWithTokens(ctx, gomock.Any()).
					Return(&models.User{
						ID:             userId,
						IdentityNumber: "PNOEE-30303039914",
						PersonalCode:   "303039914",
						FirstName:      "TESTNUMBER",
						LastName:       "OK",
						AccessToken:    "access-token",
						RefreshToken:   "refresh-token",
					}, nil)

				redis.EXPECT().DeleteSessionByID(ctx, id).Return(assert.AnError)
			},
			sessionId: sessionId,
			expected:  nil,
			error:     assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Complete(ctx, tt.sessionId)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
