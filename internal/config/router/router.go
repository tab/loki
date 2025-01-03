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

	authentication middlewares.AuthMiddleware,
	telemetry middlewares.TelemetryMiddleware,

	smartId controllers.SmartIdController,
	mobileID controllers.MobileIdController,
	sessions controllers.SessionsController,
	tokens controllers.TokensController,
	users controllers.UsersController,
) http.Handler {
	r := chi.NewRouter()

	r.Use(telemetry.Trace)
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{cfg.ClientURL},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type"},
			MaxAge:         300,
		}),
	)

	r.Post("/api/auth/smart_id", smartId.CreateSession)
	r.Post("/api/auth/mobile_id", mobileID.CreateSession)

	r.Get("/api/sessions/{id}", sessions.GetStatus)
	r.Post("/api/sessions/{id}", sessions.Authenticate)

	r.Post("/api/tokens/refresh", tokens.Refresh)

	r.Group(func(r chi.Router) {
		r.Use(authentication.Authenticate)
		r.Get("/api/me", users.Me)
	})

	return r
}
