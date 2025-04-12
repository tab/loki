package middlewares

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewAuthenticationMiddleware),
	fx.Provide(NewAuthorizationMiddleware),
	fx.Provide(NewTelemetryMiddleware),
	fx.Provide(NewLoggerMiddleware),
)
