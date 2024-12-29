package middlewares

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	AuthenticationTraceKey  = "X-Trace-ID"
	AuthenticationTraceName = "authentication"
)

type TelemetryMiddleware interface {
	Trace(next http.Handler) http.Handler
}

type telemetryMiddleware struct{}

func NewTelemetryMiddleware() TelemetryMiddleware {
	return &telemetryMiddleware{}
}

func (m *telemetryMiddleware) Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer(AuthenticationTraceName)
		traceId := r.Header.Get(AuthenticationTraceKey)

		if traceId == "" {
			traceId = uuid.New().String()
			r.Header.Set(AuthenticationTraceKey, traceId)
		}

		id, _ := trace.TraceIDFromHex(formatToTraceID(traceId))
		spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: id,
			Remote:  true,
		})

		ctx := trace.ContextWithSpanContext(r.Context(), spanCtx)
		ctx, span := tracer.Start(ctx, formatToOperationName(r.URL.Path))
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.path", r.URL.Path),
			attribute.String("http.user_agent", r.UserAgent()),
			attribute.String("http.remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func formatToOperationName(path string) string {
	parts := strings.Split(path, "/")

	for i, part := range parts {
		if isUUID(part) {
			parts[i] = "{id}"
		}
	}

	return strings.Join(parts, "/")
}

func formatToTraceID(uuid string) string {
	return strings.ReplaceAll(uuid, "-", "")
}

func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
