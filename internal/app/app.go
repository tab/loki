package app

import (
	"context"
	"net/http"
	"time"

	"github.com/tab/smartid"
	"go.uber.org/fx"

	"loki/internal/app/controllers"
	"loki/internal/app/controllers/backoffice"
	"loki/internal/app/repositories"
	"loki/internal/app/services"
	"loki/internal/app/services/authentication"
	"loki/internal/app/workers"
	"loki/internal/config"
	"loki/internal/config/middlewares"
	"loki/internal/config/router"
	"loki/internal/config/server"
	"loki/internal/config/telemetry"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

var Module = fx.Options(
	logger.Module,

	authentication.Module,
	controllers.Module,
	backoffice.Module,
	repositories.Module,
	jwt.Module,
	services.Module,
	workers.Module,

	middlewares.Module,

	server.Module,
	router.Module,
	telemetry.Module,

	fx.Invoke(registerHooks),
	fx.Invoke(registerWorkers),
	fx.Invoke(registerTelemetry),
)

func registerHooks(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	server server.Server,
	log *logger.Logger,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msgf("Starting server in %s environment at %s", cfg.AppEnv, cfg.AppAddr)

			go func() {
				if err := server.Run(); err != nil && err != http.ErrServerClosed {
					log.Error().Err(err).Msg("Server failed")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Shutting down server...")

			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			return server.Shutdown(shutdownCtx)
		},
	})
}

func registerWorkers(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	smartId smartid.Worker,
	mobileId services.MobileIdWorker,
	log *logger.Logger,
) {
	var ctx, cancel = context.WithCancel(context.Background())
	workers.Ctx = ctx

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info().Msgf("Starting workers in %s environment", cfg.AppEnv)
			smartId.Start(ctx)
			mobileId.Start(ctx)

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info().Msg("Shutting down workers ...")

			cancel()
			smartId.Stop()
			mobileId.Stop()

			return nil
		},
	})
}

func registerTelemetry(lifecycle fx.Lifecycle, cfg *config.Config) {
	var ctx, cancel = context.WithCancel(context.Background())
	service, _ := telemetry.NewTelemetry(ctx, cfg)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancel()
			return service.Shutdown(ctx)
		},
	})
}
