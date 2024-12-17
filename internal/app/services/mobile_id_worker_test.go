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
	certificateService := NewMockCertificate(ctrl)
	sessionService := NewMockSessions(ctrl)
	usersService := NewMockUsers(ctrl)
	log := logger.NewLogger()
	worker := NewMobileIdWorker(
		cfg,
		authenticationService,
		certificateService,
		sessionService,
		usersService,
		mobileIdQueue,
		log)

	cert := "MIIFXzCCA0egAwIBAgIQbZGUgkUm/YdiVSIaf13G+zANBgkqhkiG9w0BAQsFADBoMQswCQYDVQQGEwJFRTEiMCAGA1UECgwZQVMgU2VydGlmaXRzZWVyaW1pc2tlc2t1czEXMBUGA1UEYQwOTlRSRUUtMTA3NDcwMTMxHDAaBgNVBAMME1RFU1Qgb2YgRUlELVNLIDIwMTYwHhcNMjIwNDEyMDY1NDE4WhcNMjcwNDEyMjA1OTU5WjBtMQswCQYDVQQGEwJFRTEbMBkGA1UEAwwSRUlEMjAxNixURVNUTlVNQkVSMRMwEQYDVQQEDApURVNUTlVNQkVSMRAwDgYDVQQqDAdFSUQyMDE2MRowGAYDVQQFExFQTk9FRS02MDAwMTAxNzg2OTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABN5THR+i6stpLG0lq12yggjyLvyvu+0tUW2BF33CTrC019eG5oeeBYMVEm+ZBpdZlfwJUtlWXEuGw3XDdBgf14KjggHJMIIBxTAJBgNVHRMEAjAAMA4GA1UdDwEB/wQEAwIHgDB7BgNVHSAEdDByMGYGCSsGAQQBzh8SATBZMDcGCCsGAQUFBwIBFitodHRwczovL3NraWRzb2x1dGlvbnMuZXUvZW4vcmVwb3NpdG9yeS9DUFMvMB4GCCsGAQUFBwICMBIaEE9ubHkgZm9yIFRFU1RJTkcwCAYGBACPegECMB0GA1UdDgQWBBTTbr/pqvMoWlqZULY3pwCzp998ijAfBgNVHSMEGDAWgBSusOrhNvgmq6XMC2ZV/jodAr8StDB9BggrBgEFBQcBAQRxMG8wKQYIKwYBBQUHMAGGHWh0dHA6Ly9haWEuZGVtby5zay5lZS9laWQyMDE2MEIGCCsGAQUFBzAChjZodHRwczovL3NrLmVlL3VwbG9hZC9maWxlcy9URVNUX29mX0VJRC1TS18yMDE2LmRlci5jcnQwbAYIKwYBBQUHAQMEYDBeMFwGBgQAjkYBBTBSMFAWSmh0dHBzOi8vc2tpZHNvbHV0aW9ucy5ldS9lbi9yZXBvc2l0b3J5L2NvbmRpdGlvbnMtZm9yLXVzZS1vZi1jZXJ0aWZpY2F0ZXMvEwJFTjANBgkqhkiG9w0BAQsFAAOCAgEACfHa53mBsmnnnNlTa5DwXmI3R9tTjcNrjMa8alUdmvi50pipPpjvkYaCsiSJpUNNZ4EvfdRI1kKWbuCLc66MqQbZ2KNaOMZ8TODkx5uOhbOqGqWr0mBTJZJu7JboN3UB9/5lrpZlUuuhjAivpQRO1OqmQXMVfIRi/gy8sFc2l10gICZiBt2JBLmC7KiafEXK8WQJpEs8iKMHYLWubURcCrMvZXz2XdbSLfnjqS40P3uFhEJfeo6+aAsnD2M2AwBwmiF0P2/k8Vk/wc6miL8SZROJMMHqcvO9vQYlUihDSNn3Jz4fz0nzTZ4DtTSI6jGU0zy2SS2j7Srgdt+jDnIaANUEIHjsKGWqBWY77D3iNiSmN7OXLXhmsv10yjLxJSwBunFs+RvYQh1SiyvoM+/yq0SPrS/xzg0unRWAjmpGFMJ4Gw7Vn8+faJeTounOmrnZhG9LYrkLPsWxRUcxc+GfM11xUcwMExQvh7oRZ6iBibhKgLdYiP31Q4QEZpNZbjNIPVJtKEj4Z/zJlbySgb4TrXy+SgcmUilUPGowwo7wkSrw5wiG/lX9QxnQHbyUugY693BRH/twRk1jnYXie+p4wRUS4tNX9fb3hKSkIaWIMpryuHxl/WoG5/FmhDN6YpGqEaaf8/rQhUEAYae2dy5A4RgxkdtzQ2Q9uOz1qsw3T3g="

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
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "OK",
						Signature: dto.MobileIdProviderSignature{
							Value:     "signature",
							Algorithm: "algorithm",
						},
						Cert: cert,
					}, nil)

				certificateService.EXPECT().Extract(cert).Return(&CertificatePayload{
					IdentityNumber: "PNOEE-1234567890",
					PersonalCode:   "1234567890",
					FirstName:      "John",
					LastName:       "Doe",
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
					Status: "COMPLETE",
				}, nil)
			},
			error: nil,
		},
		{
			name: "Fail to create user",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "OK",
						Signature: dto.MobileIdProviderSignature{
							Value:     "signature",
							Algorithm: "algorithm",
						},
						Cert: cert,
					}, nil)

				certificateService.EXPECT().Extract(cert).Return(&CertificatePayload{
					IdentityNumber: "PNOEE-1234567890",
					PersonalCode:   "1234567890",
					FirstName:      "John",
					LastName:       "Doe",
				}, nil)

				usersService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
			error: assert.AnError,
		},
		{
			name: "Fail to update session",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "OK",
						Signature: dto.MobileIdProviderSignature{
							Value:     "signature",
							Algorithm: "algorithm",
						},
						Cert: cert,
					}, nil)

				certificateService.EXPECT().Extract(cert).Return(&CertificatePayload{
					IdentityNumber: "PNOEE-1234567890",
					PersonalCode:   "1234567890",
					FirstName:      "John",
					LastName:       "Doe",
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
			name: "SESSION_RESULT_NOT_MID_CLIENT",
			before: func(sessionId uuid.UUID) {
				authenticationService.EXPECT().GetMobileIdSessionStatus(gomock.Any(), sessionId).
					Return(&dto.MobileIdProviderSessionStatusResponse{
						State:  "COMPLETE",
						Result: "NOT_MID_CLIENT",
					}, nil)

				sessionService.EXPECT().Update(gomock.Any(), models.Session{
					ID:     sessionId,
					Status: AuthenticationError,
					Error:  "NOT_MID_CLIENT",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
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
					Status: AuthenticationError,
					Error:  "USER_CANCELLED",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
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
					Status: AuthenticationError,
					Error:  "SIGNATURE_HASH_MISMATCH",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
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
					Status: AuthenticationError,
					Error:  "PHONE_ABSENT",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
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
					Status: AuthenticationError,
					Error:  "DELIVERY_ERROR",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
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
					Status: AuthenticationError,
					Error:  "SIM_ERROR",
				}).Return(&serializers.SessionSerializer{
					ID:     sessionId,
					Status: AuthenticationError,
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
