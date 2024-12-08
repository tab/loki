package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"loki/internal/app/controllers"
	"loki/internal/config"
)

func NewRouter(
	cfg *config.Config,
	smartId controllers.SmartIdController,
	mobileID controllers.MobileIdController,
	session controllers.SessionController,
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

	r.Post("/auth/smart_id", smartId.CreateSession)
	r.Post("/auth/mobile_id", mobileID.CreateSession)
	r.Get("/auth/sessions/{id}", session.GetStatus)
	r.Post("/auth/sessions/{id}", session.Authenticate)

	return r
}
