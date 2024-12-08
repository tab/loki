package controllers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewMobileIdController),
	fx.Provide(NewSmartIdController),
	fx.Provide(NewSessionController),
)
