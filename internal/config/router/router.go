package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"loki/internal/app/controllers"
	"loki/internal/config"
	"loki/internal/config/middlewares"
)

func NewRouter(
	cfg *config.Config,

	authentication middlewares.AuthenticationMiddleware,
	telemetry middlewares.TelemetryMiddleware,
	logger middlewares.LoggerMiddleware,

	health controllers.HealthController,
	smartId controllers.SmartIdController,
	mobileID controllers.MobileIdController,
	sessions controllers.SessionsController,
	tokens controllers.TokensController,
	users controllers.UsersController,
) http.Handler {
	r := chi.NewRouter()

	r.Use(telemetry.Trace)
	r.Use(middleware.RequestID)
	r.Use(logger.Log)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://*", cfg.ClientURL},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "X-Trace-ID"},
			MaxAge:         300,
		}),
	)

	r.Get("/live", health.HandleLiveness)
	r.Get("/ready", health.HandleReadiness)

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/smart_id", smartId.CreateSession)
		r.Post("/auth/mobile_id", mobileID.CreateSession)

		r.Get("/sessions/{id}", sessions.GetStatus)
		r.Post("/sessions/{id}", sessions.Authenticate)

		r.Post("/tokens/refresh", tokens.Refresh)
	})

	r.Group(func(r chi.Router) {
		r.Use(authentication.Authenticate)
		r.Get("/api/me", users.Me)
	})

	return r
}
