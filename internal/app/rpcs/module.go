package rpcs

import (
	"go.uber.org/fx"

	"loki/internal/app/rpcs/interceptors"
	"loki/internal/app/rpcs/services"
)

var Module = fx.Options(
	interceptors.Module,
	services.Module,

	fx.Provide(NewRegistry),
)
