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
	log := logger.NewLogger()

	worker := NewSmartIdWorker(cfg, authenticationService, smartIdQueue, log)

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
							Value:            "certificate",
							CertificateLevel: "certificateLevel",
						},
						InteractionFlowUsed: "interactionFlowUsed",
					}, nil)
				authenticationService.EXPECT().UpdateSession(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "COMPLETE",
					Payload: models.SessionPayload{
						State:     "COMPLETE",
						Result:    "OK",
						Signature: "signature",
						Cert:      "certificate",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "COMPLETE",
				}, nil)
			},
			error: nil,
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
				authenticationService.EXPECT().UpdateSession(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "USER_REFUSED",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "USER_REFUSED",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
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
				authenticationService.EXPECT().UpdateSession(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "USER_REFUSED_VC_CHOICE",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "USER_REFUSED_VC_CHOICE",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
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
				authenticationService.EXPECT().UpdateSession(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "USER_REFUSED_DISPLAYTEXTANDPIN",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "USER_REFUSED_DISPLAYTEXTANDPIN",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
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
				authenticationService.EXPECT().UpdateSession(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "USER_REFUSED_VC_CHOICE",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "USER_REFUSED_VC_CHOICE",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
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
				authenticationService.EXPECT().UpdateSession(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "USER_REFUSED_CONFIRMATIONMESSAGE",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "USER_REFUSED_CONFIRMATIONMESSAGE",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
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
				authenticationService.EXPECT().UpdateSession(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
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
				authenticationService.EXPECT().UpdateSession(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "USER_REFUSED_CERT_CHOICE",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "USER_REFUSED_CERT_CHOICE",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
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
				authenticationService.EXPECT().UpdateSession(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "WRONG_VC",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "WRONG_VC",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
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
				authenticationService.EXPECT().UpdateSession(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "TIMEOUT",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "TIMEOUT",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
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
