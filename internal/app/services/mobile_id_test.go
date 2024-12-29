package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"loki/internal/app/models/dto"
	"loki/internal/config"
	"loki/pkg/logger"
)

func Test_MobileIdProvider_CreateSession(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		before   func(w http.ResponseWriter, r *http.Request)
		params   dto.CreateMobileIdSessionRequest
		expected *dto.MobileIdProviderSessionResponse
		error    bool
	}{
		{
			name: "Success",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"sessionID": "5eab0e6a-c3e7-4526-a47e-398f0d31f514", "code": "1234"}`))
			},
			params: dto.CreateMobileIdSessionRequest{
				Locale:       "ENG",
				PhoneNumber:  "+37268000769",
				PersonalCode: "60001017869",
			},
			expected: &dto.MobileIdProviderSessionResponse{
				ID:   "5eab0e6a-c3e7-4526-a47e-398f0d31f514",
				Code: "1234",
			},
			error: false,
		},
		{
			name:   "Error",
			before: func(w http.ResponseWriter, r *http.Request) {},
			params: dto.CreateMobileIdSessionRequest{
				Locale:       "ENG",
				PhoneNumber:  "not-a-phone-number",
				PersonalCode: "not-a-personal-code",
			},
			expected: &dto.MobileIdProviderSessionResponse{},
			error:    true,
		},
		{
			name: "Not found",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"title": "Not Found", "status": 404}`))
			},
			params: dto.CreateMobileIdSessionRequest{
				Locale:       "ENG",
				PhoneNumber:  "+37268000769",
				PersonalCode: "60001017869",
			},
			expected: &dto.MobileIdProviderSessionResponse{},
			error:    true,
		},
		{
			name: "Bad Request",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"title": "Bad Request", "status": 400}`))
			},
			params: dto.CreateMobileIdSessionRequest{
				Locale:       "ENG",
				PhoneNumber:  "+37268000769",
				PersonalCode: "60001017869",
			},
			expected: &dto.MobileIdProviderSessionResponse{},
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(tt.before))
			defer testServer.Close()

			cfg := &config.Config{
				MobileId: config.MobileId{
					BaseURL:          testServer.URL,
					RelyingPartyName: "DEMO",
					RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
					Text:             "Enter PIN1",
				},
			}
			log := logger.NewLogger()
			provider := NewMobileId(cfg, log)

			session, err := provider.CreateSession(ctx, tt.params)

			if tt.error {
				assert.Error(t, err)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, session.ID)
			}
		})
	}
}

func Test_MobileIdProvider_GetSessionStatus(t *testing.T) {
	id := uuid.MustParse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")

	tests := []struct {
		name      string
		before    func(w http.ResponseWriter, r *http.Request)
		sessionId uuid.UUID
		expected  *dto.MobileIdProviderSessionStatusResponse
		error     bool
	}{
		{
			name: "Success",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`
{
	"state": "COMPLETE",
	"result": "OK",
	"signature": {
		"value": "tFib1+JlILD1ZBI3w4zUeyQLCh7aVcveelCNmgmYCzKn6ZdZ6keaZRiE0wlpBRjI+eHC1CHp4NKww5kp+M+Wsg==",
		"algorithm": "SHA512WithECEncryption"
	},
	"cert": "MIIFXzCCA0egAwIBAgIQbZGUgkUm/YdiVSIaf13G+zANBgkqhkiG9w0BAQsFADBoMQswCQYDVQQGEwJFRTEiMCAGA1UECgwZQVMgU2VydGlmaXRzZWVyaW1pc2tlc2t1czEXMBUGA1UEYQwOTlRSRUUtMTA3NDcwMTMxHDAaBgNVBAMME1RFU1Qgb2YgRUlELVNLIDIwMTYwHhcNMjIwNDEyMDY1NDE4WhcNMjcwNDEyMjA1OTU5WjBtMQswCQYDVQQGEwJFRTEbMBkGA1UEAwwSRUlEMjAxNixURVNUTlVNQkVSMRMwEQYDVQQEDApURVNUTlVNQkVSMRAwDgYDVQQqDAdFSUQyMDE2MRowGAYDVQQFExFQTk9FRS02MDAwMTAxNzg2OTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABN5THR+i6stpLG0lq12yggjyLvyvu+0tUW2BF33CTrC019eG5oeeBYMVEm+ZBpdZlfwJUtlWXEuGw3XDdBgf14KjggHJMIIBxTAJBgNVHRMEAjAAMA4GA1UdDwEB/wQEAwIHgDB7BgNVHSAEdDByMGYGCSsGAQQBzh8SATBZMDcGCCsGAQUFBwIBFitodHRwczovL3NraWRzb2x1dGlvbnMuZXUvZW4vcmVwb3NpdG9yeS9DUFMvMB4GCCsGAQUFBwICMBIaEE9ubHkgZm9yIFRFU1RJTkcwCAYGBACPegECMB0GA1UdDgQWBBTTbr/pqvMoWlqZULY3pwCzp998ijAfBgNVHSMEGDAWgBSusOrhNvgmq6XMC2ZV/jodAr8StDB9BggrBgEFBQcBAQRxMG8wKQYIKwYBBQUHMAGGHWh0dHA6Ly9haWEuZGVtby5zay5lZS9laWQyMDE2MEIGCCsGAQUFBzAChjZodHRwczovL3NrLmVlL3VwbG9hZC9maWxlcy9URVNUX29mX0VJRC1TS18yMDE2LmRlci5jcnQwbAYIKwYBBQUHAQMEYDBeMFwGBgQAjkYBBTBSMFAWSmh0dHBzOi8vc2tpZHNvbHV0aW9ucy5ldS9lbi9yZXBvc2l0b3J5L2NvbmRpdGlvbnMtZm9yLXVzZS1vZi1jZXJ0aWZpY2F0ZXMvEwJFTjANBgkqhkiG9w0BAQsFAAOCAgEACfHa53mBsmnnnNlTa5DwXmI3R9tTjcNrjMa8alUdmvi50pipPpjvkYaCsiSJpUNNZ4EvfdRI1kKWbuCLc66MqQbZ2KNaOMZ8TODkx5uOhbOqGqWr0mBTJZJu7JboN3UB9/5lrpZlUuuhjAivpQRO1OqmQXMVfIRi/gy8sFc2l10gICZiBt2JBLmC7KiafEXK8WQJpEs8iKMHYLWubURcCrMvZXz2XdbSLfnjqS40P3uFhEJfeo6+aAsnD2M2AwBwmiF0P2/k8Vk/wc6miL8SZROJMMHqcvO9vQYlUihDSNn3Jz4fz0nzTZ4DtTSI6jGU0zy2SS2j7Srgdt+jDnIaANUEIHjsKGWqBWY77D3iNiSmN7OXLXhmsv10yjLxJSwBunFs+RvYQh1SiyvoM+/yq0SPrS/xzg0unRWAjmpGFMJ4Gw7Vn8+faJeTounOmrnZhG9LYrkLPsWxRUcxc+GfM11xUcwMExQvh7oRZ6iBibhKgLdYiP31Q4QEZpNZbjNIPVJtKEj4Z/zJlbySgb4TrXy+SgcmUilUPGowwo7wkSrw5wiG/lX9QxnQHbyUugY693BRH/twRk1jnYXie+p4wRUS4tNX9fb3hKSkIaWIMpryuHxl/WoG5/FmhDN6YpGqEaaf8/rQhUEAYae2dy5A4RgxkdtzQ2Q9uOz1qsw3T3g=",
	"time": "2024-12-08T22:07:52",
	"traceId": "9866bf3f8f4642eb"
}`))
			},
			sessionId: id,
			expected: &dto.MobileIdProviderSessionStatusResponse{
				State:  "COMPLETE",
				Result: "OK",
				Signature: dto.MobileIdProviderSignature{
					Value:     "tFib1+JlILD1ZBI3w4zUeyQLCh7aVcveelCNmgmYCzKn6ZdZ6keaZRiE0wlpBRjI+eHC1CHp4NKww5kp+M+Wsg==",
					Algorithm: "SHA512WithECEncryption",
				},
				Cert: "MIIFXzCCA0egAwIBAgIQbZGUgkUm/YdiVSIaf13G+zANBgkqhkiG9w0BAQsFADBoMQswCQYDVQQGEwJFRTEiMCAGA1UECgwZQVMgU2VydGlmaXRzZWVyaW1pc2tlc2t1czEXMBUGA1UEYQwOTlRSRUUtMTA3NDcwMTMxHDAaBgNVBAMME1RFU1Qgb2YgRUlELVNLIDIwMTYwHhcNMjIwNDEyMDY1NDE4WhcNMjcwNDEyMjA1OTU5WjBtMQswCQYDVQQGEwJFRTEbMBkGA1UEAwwSRUlEMjAxNixURVNUTlVNQkVSMRMwEQYDVQQEDApURVNUTlVNQkVSMRAwDgYDVQQqDAdFSUQyMDE2MRowGAYDVQQFExFQTk9FRS02MDAwMTAxNzg2OTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABN5THR+i6stpLG0lq12yggjyLvyvu+0tUW2BF33CTrC019eG5oeeBYMVEm+ZBpdZlfwJUtlWXEuGw3XDdBgf14KjggHJMIIBxTAJBgNVHRMEAjAAMA4GA1UdDwEB/wQEAwIHgDB7BgNVHSAEdDByMGYGCSsGAQQBzh8SATBZMDcGCCsGAQUFBwIBFitodHRwczovL3NraWRzb2x1dGlvbnMuZXUvZW4vcmVwb3NpdG9yeS9DUFMvMB4GCCsGAQUFBwICMBIaEE9ubHkgZm9yIFRFU1RJTkcwCAYGBACPegECMB0GA1UdDgQWBBTTbr/pqvMoWlqZULY3pwCzp998ijAfBgNVHSMEGDAWgBSusOrhNvgmq6XMC2ZV/jodAr8StDB9BggrBgEFBQcBAQRxMG8wKQYIKwYBBQUHMAGGHWh0dHA6Ly9haWEuZGVtby5zay5lZS9laWQyMDE2MEIGCCsGAQUFBzAChjZodHRwczovL3NrLmVlL3VwbG9hZC9maWxlcy9URVNUX29mX0VJRC1TS18yMDE2LmRlci5jcnQwbAYIKwYBBQUHAQMEYDBeMFwGBgQAjkYBBTBSMFAWSmh0dHBzOi8vc2tpZHNvbHV0aW9ucy5ldS9lbi9yZXBvc2l0b3J5L2NvbmRpdGlvbnMtZm9yLXVzZS1vZi1jZXJ0aWZpY2F0ZXMvEwJFTjANBgkqhkiG9w0BAQsFAAOCAgEACfHa53mBsmnnnNlTa5DwXmI3R9tTjcNrjMa8alUdmvi50pipPpjvkYaCsiSJpUNNZ4EvfdRI1kKWbuCLc66MqQbZ2KNaOMZ8TODkx5uOhbOqGqWr0mBTJZJu7JboN3UB9/5lrpZlUuuhjAivpQRO1OqmQXMVfIRi/gy8sFc2l10gICZiBt2JBLmC7KiafEXK8WQJpEs8iKMHYLWubURcCrMvZXz2XdbSLfnjqS40P3uFhEJfeo6+aAsnD2M2AwBwmiF0P2/k8Vk/wc6miL8SZROJMMHqcvO9vQYlUihDSNn3Jz4fz0nzTZ4DtTSI6jGU0zy2SS2j7Srgdt+jDnIaANUEIHjsKGWqBWY77D3iNiSmN7OXLXhmsv10yjLxJSwBunFs+RvYQh1SiyvoM+/yq0SPrS/xzg0unRWAjmpGFMJ4Gw7Vn8+faJeTounOmrnZhG9LYrkLPsWxRUcxc+GfM11xUcwMExQvh7oRZ6iBibhKgLdYiP31Q4QEZpNZbjNIPVJtKEj4Z/zJlbySgb4TrXy+SgcmUilUPGowwo7wkSrw5wiG/lX9QxnQHbyUugY693BRH/twRk1jnYXie+p4wRUS4tNX9fb3hKSkIaWIMpryuHxl/WoG5/FmhDN6YpGqEaaf8/rQhUEAYae2dy5A4RgxkdtzQ2Q9uOz1qsw3T3g=",
			},
			error: false,
		},
		{
			name: "NOT_MID_CLIENT",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"state": "COMPLETE", "result": "NOT_MID_CLIENT"}`))
			},
			sessionId: id,
			expected: &dto.MobileIdProviderSessionStatusResponse{
				State:  "COMPLETE",
				Result: "NOT_MID_CLIENT",
			},
			error: false,
		},
		{
			name: "TIMEOUT",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"state": "COMPLETE", "result": "TIMEOUT"}`))
			},
			sessionId: id,
			expected: &dto.MobileIdProviderSessionStatusResponse{
				State:  "COMPLETE",
				Result: "TIMEOUT",
			},
			error: false,
		},
		{
			name: "Not found",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"title": "Not Found", "status": 404}`))
			},
			sessionId: id,
			expected:  &dto.MobileIdProviderSessionStatusResponse{},
			error:     true,
		},
		{
			name: "Bad Request",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"title": "Bad Request", "status": 400}`))
			},
			sessionId: id,
			expected:  &dto.MobileIdProviderSessionStatusResponse{},
			error:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(tt.before))
			defer testServer.Close()

			cfg := &config.Config{
				MobileId: config.MobileId{
					BaseURL:          testServer.URL,
					RelyingPartyName: "DEMO",
					RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
					Text:             "Enter PIN1",
				},
			}
			log := logger.NewLogger()
			provider := NewMobileId(cfg, log)

			session, err := provider.GetSessionStatus(tt.sessionId)

			if tt.error {
				assert.Error(t, err)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.State, session.State)
				assert.Equal(t, tt.expected.Result, session.Result)
				assert.Equal(t, tt.expected.Signature, session.Signature)
				assert.Equal(t, tt.expected.Cert, session.Cert)
			}
		})
	}
}
