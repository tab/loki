package interceptors

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Test_TraceInterceptor_Trace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	interceptorInstance := NewTraceInterceptor()

	traceId := uuid.New().String()
	requestId := uuid.New().String()

	type result struct {
		traceId   string
		requestId string
	}

	tests := []struct {
		name   string
		method string
		md     metadata.MD
		expect result
	}{
		{
			name:   "Success",
			method: "/sso.v1.PermissionService/List",
			md: metadata.Pairs(
				TraceId, traceId,
				RequestId, requestId,
			),
			expect: result{
				traceId:   traceId,
				requestId: requestId,
			},
		},
		{
			name: "No traceId in context",
			md: metadata.Pairs(
				RequestId, requestId,
			),
			expect: result{
				traceId:   "",
				requestId: requestId,
			},
		},
		{
			name: "No requestId in context",
			md: metadata.Pairs(
				TraceId, traceId,
			),
			expect: result{
				traceId:   traceId,
				requestId: "",
			},
		},
		{
			name: "No traceId and requestId in context",
			md:   metadata.Pairs(),
			expect: result{
				traceId:   "",
				requestId: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(), tt.md)

			var actualTraceId, actualRequestId string

			mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				md, ok := metadata.FromOutgoingContext(ctx)
				assert.True(t, ok)

				traceIds := md.Get(TraceId)
				assert.NotEmpty(t, traceIds)
				actualTraceId = traceIds[0]

				requestIds := md.Get(RequestId)
				assert.NotEmpty(t, requestIds)
				actualRequestId = requestIds[0]

				return "test response", nil
			}

			interceptor := interceptorInstance.Trace()

			resp, err := interceptor(
				ctx,
				"test request",
				&grpc.UnaryServerInfo{FullMethod: tt.method},
				mockHandler,
			)

			assert.NoError(t, err)
			assert.Equal(t, "test response", resp)

			if tt.expect.traceId != "" {
				assert.Equal(t, tt.expect.traceId, actualTraceId)
			} else {
				_, err = uuid.Parse(actualTraceId)
				assert.NoError(t, err)
			}

			if tt.expect.requestId != "" {
				assert.Equal(t, tt.expect.requestId, actualRequestId)
			} else {
				_, err = uuid.Parse(actualRequestId)
				assert.NoError(t, err)
			}
		})
	}
}
