package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki/internal/app/serializers"
)

func Test_HealthController_HandleLiveness(t *testing.T) {
	handler := NewHealthController()

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
