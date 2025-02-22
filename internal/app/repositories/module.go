package repositories

import (
	"go.uber.org/fx"

	"loki/internal/app/repositories/postgres"
	"loki/internal/app/repositories/redis"
)

var Module = fx.Options(
	fx.Provide(postgres.NewPostgresClient),
	fx.Provide(redis.NewRedisClient),

	fx.Provide(NewHealthRepository),
	fx.Provide(NewSessionRepository),
	fx.Provide(NewPermissionRepository),
	fx.Provide(NewRoleRepository),
	fx.Provide(NewScopeRepository),
	fx.Provide(NewTokenRepository),
	fx.Provide(NewUserRepository),
)
