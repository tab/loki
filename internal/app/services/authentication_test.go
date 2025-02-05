package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/config"
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

	sessionsService := NewMockSessions(ctrl)
	tokensService := NewMockTokens(ctrl)
	log := logger.NewLogger()

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		sessionsService,
		tokensService,
		log)

	id := uuid.MustParse("8fdb516d-1a82-43ba-b82d-be63df569b86")
	sessionId := id.String()

	tests := []struct {
		name     string
		before   func()
		params   dto.CreateSmartIdSessionRequest
		expected *models.Session
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

				sessionsService.EXPECT().Create(ctx, &models.CreateSessionParams{
					SessionId: sessionId,
					Code:      "1234",
				}).Return(&models.Session{
					ID:   id,
					Code: "1234",
				}, nil)
			},
			params: dto.CreateSmartIdSessionRequest{
				Country:      "EE",
				PersonalCode: "30303039914",
			},
			expected: &models.Session{
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

				sessionsService.EXPECT().Create(ctx, &models.CreateSessionParams{
					SessionId: sessionId,
					Code:      "1234",
				}).Return(nil, assert.AnError)
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

	sessionsService := NewMockSessions(ctrl)
	tokensService := NewMockTokens(ctrl)
	log := logger.NewLogger()

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		sessionsService,
		tokensService,
		log)

	id := uuid.MustParse("8fdb516d-1a82-43ba-b82d-be63df569b86")

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

	sessionsService := NewMockSessions(ctrl)
	tokensService := NewMockTokens(ctrl)
	log := logger.NewLogger()

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		sessionsService,
		tokensService,
		log)

	id := uuid.MustParse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")
	sessionId := id.String()

	tests := []struct {
		name     string
		before   func()
		params   dto.CreateMobileIdSessionRequest
		expected *models.Session
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

				sessionsService.EXPECT().Create(ctx, &models.CreateSessionParams{
					SessionId: sessionId,
					Code:      "1234",
				}).Return(&models.Session{
					ID:   id,
					Code: "1234",
				}, nil)
			},
			params: dto.CreateMobileIdSessionRequest{
				Locale:       "ENG",
				PhoneNumber:  "+37268000769",
				PersonalCode: "60001017869",
			},
			expected: &models.Session{
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

				sessionsService.EXPECT().Create(ctx, &models.CreateSessionParams{
					SessionId: sessionId,
					Code:      "1234",
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

	sessionsService := NewMockSessions(ctrl)
	tokensService := NewMockTokens(ctrl)
	log := logger.NewLogger()

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		sessionsService,
		tokensService,
		log)

	id := uuid.MustParse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")

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

	sessionsService := NewMockSessions(ctrl)
	tokensService := NewMockTokens(ctrl)
	log := logger.NewLogger()

	id := uuid.MustParse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")
	sessionId := id.String()
	userId := uuid.MustParse("320284a1-8c96-4984-b492-b060310cfdac")

	service := NewAuthentication(
		cfg,
		smartIdMock,
		smartIdQueue,
		mobileIdMock,
		mobileIdQueue,
		sessionsService,
		tokensService,
		log)

	tests := []struct {
		name     string
		before   func()
		expected *models.User
		error    error
	}{
		{
			name: "Success (smart-id)",
			before: func() {
				sessionsService.EXPECT().FindById(ctx, sessionId).Return(&models.Session{
					ID:     id,
					UserId: userId,
					Status: AuthenticationSuccess,
				}, nil)

				tokensService.EXPECT().Create(ctx, gomock.Any()).Return(&models.User{
					ID:             userId,
					IdentityNumber: "PNOEE-30303039914",
					PersonalCode:   "303039914",
					FirstName:      "TESTNUMBER",
					LastName:       "OK",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)

				sessionsService.EXPECT().Delete(ctx, sessionId).Return(nil)
			},
			expected: &models.User{
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
				sessionsService.EXPECT().FindById(ctx, sessionId).Return(&models.Session{
					ID:     id,
					UserId: userId,
					Status: AuthenticationSuccess,
				}, nil)

				tokensService.EXPECT().Create(ctx, gomock.Any()).Return(&models.User{
					ID:             userId,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)

				sessionsService.EXPECT().Delete(ctx, sessionId).Return(nil)
			},
			expected: &models.User{
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
				sessionsService.EXPECT().FindById(ctx, sessionId).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
		{
			name: "Error: failed to delete session",
			before: func() {
				sessionsService.EXPECT().FindById(ctx, sessionId).Return(&models.Session{
					ID:     id,
					UserId: userId,
					Status: AuthenticationSuccess,
				}, nil)

				tokensService.EXPECT().Create(ctx, gomock.Any()).Return(&models.User{
					ID:             userId,
					IdentityNumber: "PNOEE-30303039914",
					PersonalCode:   "303039914",
					FirstName:      "TESTNUMBER",
					LastName:       "OK",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)

				sessionsService.EXPECT().Delete(ctx, sessionId).Return(assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Complete(ctx, sessionId)

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
