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

	authMiddleware middlewares.AuthMiddleware,

	smartId controllers.SmartIdController,
	mobileID controllers.MobileIdController,
	sessions controllers.SessionsController,
	users controllers.UsersController,
) http.Handler {
	r := chi.NewRouter()

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
	r.Get("/api/auth/sessions/{id}", sessions.GetStatus)
	r.Post("/api/auth/sessions/{id}/authenticate", sessions.Authenticate)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Get("/api/me", users.Me)
	})

	return r
}
