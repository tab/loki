package repositories

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewDatabase),
	fx.Provide(NewRedis),
)
