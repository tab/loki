package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"loki/internal/config/logger"
)

type LoggerMiddleware interface {
	Log(next http.Handler) http.Handler
}

type loggerMiddleware struct {
	log *logger.Logger
}

func NewLoggerMiddleware(log *logger.Logger) LoggerMiddleware {
	return &loggerMiddleware{
		log: log,
	}
}

func (m *loggerMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		traceId, _ := CurrentTraceIdFromContext(r.Context())
		requestId := middleware.GetReqID(r.Context())

		reqLogger := m.log.
			WithComponent("http").
			WithRequestId(requestId).
			WithTraceId(traceId)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		status := ww.Status()
		size := ww.BytesWritten()
		duration := time.Since(startTime)

		reqLogger.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Str("proto", r.Proto).
			Str("scheme", scheme).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msgf("%s %s - %d %dB in %s", r.Method, r.RequestURI, status, size, duration)
	})
}
