package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"loki/internal/config/logger"
)

type LoggerInterceptor interface {
	Log() grpc.UnaryServerInterceptor
}

type loggerInterceptor struct {
	log *logger.Logger
}

func NewLoggerInterceptor(log *logger.Logger) LoggerInterceptor {
	return &loggerInterceptor{
		log: log,
	}
}

func (i *loggerInterceptor) Log() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		traceId := extractTraceId(ctx)
		requestId := extractRequestId(ctx)

		reqLogger := i.log.
			WithComponent("gRPC").
			WithRequestId(requestId).
			WithTraceId(traceId)

		resp, err := handler(ctx, req)

		code := status.Code(err).String()
		duration := time.Since(startTime)

		reqLogger.Info().
			Str("method", info.FullMethod).
			Str("status", code).
			Dur("duration", duration).
			Msgf("%s - %s in %s", info.FullMethod, code, duration)

		return resp, err
	}
}
