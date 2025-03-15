package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"loki/internal/app/controllers"
	"loki/internal/app/controllers/backoffice"
	"loki/internal/config"
	"loki/internal/config/middlewares"
	"loki/pkg/rbac"
)

func NewRouter(
	cfg *config.Config,

	authentication middlewares.AuthenticationMiddleware,
	authorization middlewares.AuthorizationMiddleware,
	telemetry middlewares.TelemetryMiddleware,

	health controllers.HealthController,
	smartId controllers.SmartIdController,
	mobileID controllers.MobileIdController,
	sessions controllers.SessionsController,
	tokens controllers.TokensController,
	users controllers.UsersController,

	backofficePermissions backoffice.PermissionsController,
	backofficeRoles backoffice.RolesController,
	backofficeScopes backoffice.ScopesController,
	backofficeTokens backoffice.TokensController,
	backofficeUsers backoffice.UsersController,
) http.Handler {
	r := chi.NewRouter()

	r.Use(telemetry.Trace)
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://*", cfg.ClientURL},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Trace-ID"},
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

	r.Group(func(r chi.Router) {
		r.Use(authorization.Authorize)

		r.Route("/api/backoffice", func(r chi.Router) {
			r.With(authorization.Check(rbac.ReadPermissions)).Get("/permissions", backofficePermissions.List)
			r.With(authorization.Check(rbac.ReadPermissions)).Get("/permissions/{id}", backofficePermissions.Get)
			r.With(authorization.Check(rbac.WritePermissions)).Post("/permissions", backofficePermissions.Create)
			r.With(authorization.Check(rbac.WritePermissions)).Put("/permissions/{id}", backofficePermissions.Update)
			r.With(authorization.Check(rbac.WritePermissions)).Delete("/permissions/{id}", backofficePermissions.Delete)

			r.With(authorization.Check(rbac.ReadRoles)).Get("/roles", backofficeRoles.List)
			r.With(authorization.Check(rbac.ReadRoles)).Get("/roles/{id}", backofficeRoles.Get)
			r.With(authorization.Check(rbac.WriteRoles)).Post("/roles", backofficeRoles.Create)
			r.With(authorization.Check(rbac.WriteRoles)).Put("/roles/{id}", backofficeRoles.Update)
			r.With(authorization.Check(rbac.WriteRoles)).Delete("/roles/{id}", backofficeRoles.Delete)

			r.With(authorization.Check(rbac.ReadScopes)).Get("/scopes", backofficeScopes.List)
			r.With(authorization.Check(rbac.ReadScopes)).Get("/scopes/{id}", backofficeScopes.Get)
			r.With(authorization.Check(rbac.WriteScopes)).Post("/scopes", backofficeScopes.Create)
			r.With(authorization.Check(rbac.WriteScopes)).Put("/scopes/{id}", backofficeScopes.Update)
			r.With(authorization.Check(rbac.WriteScopes)).Delete("/scopes/{id}", backofficeScopes.Delete)

			r.With(authorization.Check(rbac.ReadTokens)).Get("/tokens", backofficeTokens.List)
			r.With(authorization.Check(rbac.WriteTokens)).Delete("/tokens/{id}", backofficeTokens.Delete)

			r.With(authorization.Check(rbac.ReadUsers)).Get("/users", backofficeUsers.List)
			r.With(authorization.Check(rbac.ReadUsers)).Get("/users/{id}", backofficeUsers.Get)
			r.With(authorization.Check(rbac.WriteUsers)).Post("/users", backofficeUsers.Create)
			r.With(authorization.Check(rbac.WriteUsers)).Put("/users/{id}", backofficeUsers.Update)
			r.With(authorization.Check(rbac.WriteUsers)).Delete("/users/{id}", backofficeUsers.Delete)
		})
	})

	return r
}
