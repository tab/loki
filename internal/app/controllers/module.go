package controllers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthController),
	fx.Provide(NewMobileIdController),
	fx.Provide(NewSmartIdController),
	fx.Provide(NewSessionsController),
	fx.Provide(NewTokensController),
	fx.Provide(NewUsersController),
)
