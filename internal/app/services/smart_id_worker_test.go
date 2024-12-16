package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/app/serializers"
	"loki/internal/config"
	"loki/pkg/logger"
)

func Test_SmartIdWorker_Perform(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		SmartId: config.SmartId{
			BaseURL:          "https://sid.demo.sk.ee/smart-id-rp/v2",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
	}

	smartIdQueue := make(chan *SmartIdQueue, 1)
	authenticationService := NewMockAuthentication(ctrl)
	sessionService := NewMockSessions(ctrl)
	usersService := NewMockUsers(ctrl)
	log := logger.NewLogger()

	worker := NewSmartIdWorker(
		cfg,
		authenticationService,
		sessionService,
		usersService,
		smartIdQueue,
		log)

	certificate := "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"

	userId, err := uuid.NewRandom()
	assert.NoError(t, err)

	tests := []struct {
		name   string
		before func(sessionId uuid.UUID)
		error  error
	}{
		{
			name: "Success",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "OK",
						},
						Signature: dto.SmartIdProviderSignature{
							Value:     "signature",
							Algorithm: "algorithm",
						},
						Cert: dto.SmartIdProviderCertificate{
							Value:            certificate,
							CertificateLevel: "certificateLevel",
						},
						InteractionFlowUsed: "interactionFlowUsed",
					}, nil)

				usersService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.User{
					ID:             userId,
					IdentityNumber: "PNOEE-1234567890",
					PersonalCode:   "1234567890",
					FirstName:      "John",
					LastName:       "Doe",
				}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					UserId: userId,
					Status: AuthenticationSuccess,
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationSuccess,
				}, nil)
			},
			error: nil,
		},
		{
			name: "Failed to create user",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "OK",
						},
						Signature: dto.SmartIdProviderSignature{
							Value:     "signature",
							Algorithm: "algorithm",
						},
						Cert: dto.SmartIdProviderCertificate{
							Value:            certificate,
							CertificateLevel: "certificateLevel",
						},
						InteractionFlowUsed: "interactionFlowUsed",
					}, nil)

				usersService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
			error: assert.AnError,
		},
		{
			name: "Failed to update session",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "OK",
						},
						Signature: dto.SmartIdProviderSignature{
							Value:     "signature",
							Algorithm: "algorithm",
						},
						Cert: dto.SmartIdProviderCertificate{
							Value:            certificate,
							CertificateLevel: "certificateLevel",
						},
						InteractionFlowUsed: "interactionFlowUsed",
					}, nil)

				usersService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.User{
					ID:             userId,
					IdentityNumber: "PNOEE-1234567890",
					PersonalCode:   "1234567890",
					FirstName:      "John",
					LastName:       "Doe",
				}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					UserId: userId,
					Status: AuthenticationSuccess,
				}).Return(nil, assert.AnError)
			},
			error: assert.AnError,
		},
		{
			name: "SESSION_RESULT_USER_REFUSED",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "USER_REFUSED",
						},
					}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_USER_REFUSED_VC_CHOICE",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "USER_REFUSED_VC_CHOICE",
						},
					}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_VC_CHOICE",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_VC_CHOICE",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_USER_REFUSED_DISPLAYTEXTANDPIN",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "USER_REFUSED_DISPLAYTEXTANDPIN",
						},
					}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_DISPLAYTEXTANDPIN",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_DISPLAYTEXTANDPIN",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_USER_REFUSED_VC_CHOICE",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "USER_REFUSED_VC_CHOICE",
						},
					}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_VC_CHOICE",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_VC_CHOICE",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "USER_REFUSED_CONFIRMATIONMESSAGE",
						},
					}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_CONFIRMATIONMESSAGE",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_CONFIRMATIONMESSAGE",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE",
						},
					}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_USER_REFUSED_CERT_CHOICE",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "USER_REFUSED_CERT_CHOICE",
						},
					}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_CERT_CHOICE",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "USER_REFUSED_CERT_CHOICE",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_WRONG_VC",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "WRONG_VC",
						},
					}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "WRONG_VC",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "WRONG_VC",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_TIMEOUT",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.SmartIdProviderSessionStatusResponse{
						State: "COMPLETE",
						Result: dto.SmartIdProviderResult{
							EndResult: "TIMEOUT",
						},
					}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "TIMEOUT",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "TIMEOUT",
				}, nil)
			},
			error: nil,
		},
		{
			name: "Error",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetSmartIdSessionStatus(gomock.Any(), sessionId).
					Return(nil, assert.AnError)
			},
			error: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionId, err := uuid.NewRandom()
			assert.NoError(t, err)

			tt.before(sessionId)

			ctx, cancel := context.WithCancel(context.Background())
			worker.Start(ctx)

			smartIdQueue <- &SmartIdQueue{
				ID: sessionId,
			}
			time.Sleep(5 * time.Millisecond)

			cancel()
			worker.Stop()

			if tt.error != nil {
				assert.Error(t, tt.error)
			} else {
				assert.NoError(t, tt.error)
			}
		})
	}
}
