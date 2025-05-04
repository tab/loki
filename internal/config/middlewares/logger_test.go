package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki/internal/config"
	"loki/internal/config/logger"
)

func Test_NewLoggerMiddleware(t *testing.T) {
	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	middleware := NewLoggerMiddleware(log)
	assert.NotNil(t, middleware)
}

func Test_LoggerMiddleware_Log(t *testing.T) {
	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)
	middleware := NewLoggerMiddleware(log)

	type result struct {
		code   int
		status string
	}

	tests := []struct {
		name     string
		traceId  string
		expected result
	}{
		{
			name:    "Success",
			traceId: "test-trace-id",
			expected: result{
				code:   http.StatusOK,
				status: "200 OK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("Success"))
			})

			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)

			req.Header.Set(TraceKey, tt.traceId)

			ctx := NewContextModifier(req.Context()).
				WithTraceId(tt.traceId).
				Context()
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			middleware.Log(handler).ServeHTTP(rr, req)

			assert.Equal(t, tt.expected.code, rr.Code)
			assert.Equal(t, tt.expected.status, rr.Result().Status)
		})
	}
}
