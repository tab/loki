package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

func Test_HealthController_HandleLiveness(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := services.NewMockHealthChecker(ctrl)
	handler := NewHealthController(service)

	type result struct {
		response serializers.HealthSerializer
		code     int
		status   string
	}

	tests := []struct {
		name     string
		expected result
	}{
		{
			name: "Success",
			expected: result{
				response: serializers.HealthSerializer{Result: "alive"},
				code:     http.StatusOK,
				status:   "200 OK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/live", nil)
			w := httptest.NewRecorder()

			handler.HandleLiveness(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			var actual serializers.HealthSerializer
			err := json.NewDecoder(resp.Body).Decode(&actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.response.Result, actual.Result)
			assert.Equal(t, tt.expected.status, resp.Status)
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}

func Test_HealthController_HandleReadiness(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := services.NewMockHealthChecker(ctrl)
	handler := NewHealthController(service)

	type result struct {
		response serializers.HealthSerializer
		error    serializers.ErrorSerializer
		code     int
		status   string
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				service.EXPECT().Ping(gomock.Any()).Return(nil)
			},
			expected: result{
				response: serializers.HealthSerializer{Result: "ready"},
				code:     http.StatusOK,
				status:   "200 OK",
			},
		},
		{
			name: "Error",
			before: func() {
				service.EXPECT().Ping(gomock.Any()).Return(assert.AnError)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: "unavailable"},
				code:   http.StatusServiceUnavailable,
				status: "503 Service Unavailable",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest("GET", "/ready", nil)
			w := httptest.NewRecorder()

			handler.HandleReadiness(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.expected.error.Error != "" {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error.Error, response.Error)
			} else {
				var actual serializers.HealthSerializer
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response.Result, actual.Result)
			}

			assert.Equal(t, tt.expected.status, resp.Status)
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}
