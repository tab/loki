package interceptors

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	RequestId = "X-Request-Id"
	TraceId   = "X-Trace-Id"
)

type TraceInterceptor interface {
	Trace() grpc.UnaryServerInterceptor
}

type traceInterceptor struct{}

func NewTraceInterceptor() TraceInterceptor {
	return &traceInterceptor{}
}

func (i *traceInterceptor) Trace() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		traceId := extractTraceId(ctx)
		if traceId == "" {
			traceId = uuid.New().String()
		}

		requestId := extractRequestId(ctx)
		if requestId == "" {
			requestId = uuid.New().String()
		}

		ctx = metadata.AppendToOutgoingContext(ctx, TraceId, traceId)
		ctx = metadata.AppendToOutgoingContext(ctx, RequestId, requestId)

		return handler(ctx, req)
	}
}

func extractTraceId(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	traceId := md.Get(TraceId)
	if len(traceId) == 0 {
		return ""
	}

	return traceId[0]
}

func extractRequestId(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	requestId := md.Get(RequestId)
	if len(requestId) == 0 {
		return ""
	}

	return requestId[0]
}
