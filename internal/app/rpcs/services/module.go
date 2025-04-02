package services

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewPermissions),
	fx.Provide(NewRoles),
	fx.Provide(NewScopes),
	fx.Provide(NewTokens),
)
