package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	"loki/internal/config"
	"loki/internal/config/logger"
)

func Test_LoggerInterceptor_UnaryServerInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)
	interceptorInstance := NewLoggerInterceptor(log)

	tests := []struct {
		name   string
		method string
		error  bool
	}{
		{
			name:   "Success",
			method: "/sso.v1.PermissionService/List",
			error:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := interceptorInstance.Log()

			mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "test response", nil
			}

			resp, err := interceptor(
				context.Background(),
				"test request",
				&grpc.UnaryServerInfo{FullMethod: tt.method},
				mockHandler,
			)

			assert.NoError(t, err)
			assert.Equal(t, "test response", resp)
		})
	}
}
