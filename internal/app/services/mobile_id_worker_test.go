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

func Test_MobileIdWorker_Perform(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		MobileId: config.MobileId{
			BaseURL:          "https://tsp.demo.sk.ee/mid-api",
			RelyingPartyName: "DEMO",
			RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
			Text:             "Enter PIN1",
		},
	}

	mobileIdQueue := make(chan *MobileIdQueue, 1)
	authenticationService := NewMockAuthentication(ctrl)
	sessionService := NewMockSessions(ctrl)
	log := logger.NewLogger()
	worker := NewMobileIdWorker(cfg, authenticationService, sessionService, mobileIdQueue, log)

	tests := []struct {
		name   string
		before func(sessionId uuid.UUID)
		error  error
	}{
		{
			name: "Success",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "OK",
						Signature: dto.MobileIdProviderSignature{
							Value:     "signature",
							Algorithm: "algorithm",
						},
						Cert: "certificate",
					}, nil)
				sessionService.EXPECT().Update(gomock.Any(), models.Session{
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
			name: "SESSION_RESULT_NOT_MID_CLIENT",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "NOT_MID_CLIENT",
					}, nil)
				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "NOT_MID_CLIENT",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "NOT_MID_CLIENT",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "NOT_MID_CLIENT",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_USER_CANCELLED",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "USER_CANCELLED",
					}, nil)
				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "USER_CANCELLED",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "USER_CANCELLED",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "USER_CANCELLED",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_SIGNATURE_HASH_MISMATCH",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "SIGNATURE_HASH_MISMATCH",
					}, nil)
				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "SIGNATURE_HASH_MISMATCH",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "SIGNATURE_HASH_MISMATCH",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "SIGNATURE_HASH_MISMATCH",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_PHONE_ABSENT",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "PHONE_ABSENT",
					}, nil)
				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "PHONE_ABSENT",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "PHONE_ABSENT",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "PHONE_ABSENT",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_DELIVERY_ERROR",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "DELIVERY_ERROR",
					}, nil)
				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "DELIVERY_ERROR",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "DELIVERY_ERROR",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "DELIVERY_ERROR",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_SIM_ERROR",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "SIM_ERROR",
					}, nil)
				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "SIM_ERROR",
					Payload: models.SessionPayload{
						State:  "COMPLETE",
						Result: "SIM_ERROR",
					},
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: "ERROR",
					Error:  "SIM_ERROR",
				}, nil)
			},
			error: nil,
		},
		{
			name: "SESSION_RESULT_TIMEOUT",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "TIMEOUT",
					}, nil)
				sessionService.EXPECT().Update(gomock.Any(), models.Session{
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
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
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

			mobileIdQueue <- &MobileIdQueue{
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
