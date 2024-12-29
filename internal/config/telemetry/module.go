package telemetry

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewTelemetry),
)
