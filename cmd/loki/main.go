package main

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"loki/internal/app"
	"loki/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	fx.New(
		fx.WithLogger(
			func() fxevent.Logger {
				if cfg.LogLevel == config.DebugLevel {
					return &fxevent.ConsoleLogger{W: os.Stdout}
				}
				return fxevent.NopLogger
			},
		),
		fx.Supply(cfg),
		app.Module,
	).Run()
}
