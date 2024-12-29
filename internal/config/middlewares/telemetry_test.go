package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_TelemetryMiddleware_Trace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	middleware := NewTelemetryMiddleware()

	type result struct {
		status string
		code   int
	}

	tests := []struct {
		name     string
		header   string
		expected result
		error    error
	}{
		{
			name:   "Success",
			header: "8f963243-726d-4603-af9c-271eeb15c4a2",
			expected: result{
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode("ok")
			})

			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Set("X-Trace-ID", tt.header)
			rw := httptest.NewRecorder()

			middleware.Trace(handler).ServeHTTP(rw, req)

			res := rw.Result()
			defer res.Body.Close()

			if tt.error != nil {
				assert.Error(t, tt.error)
			} else {
				assert.Equal(t, tt.expected.code, res.StatusCode)
				assert.Equal(t, tt.expected.status, res.Status)
			}
		})
	}
}
