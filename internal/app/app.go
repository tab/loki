package app

import (
	"context"
	"net/http"
	"time"

	"github.com/tab/mobileid"
	"github.com/tab/smartid"
	"go.uber.org/fx"

	"loki/internal/app/controllers"
	"loki/internal/app/repositories"
	"loki/internal/app/rpcs"
	"loki/internal/app/services"
	"loki/internal/app/services/authentication"
	"loki/internal/app/workers"
	"loki/internal/config"
	"loki/internal/config/logger"
	"loki/internal/config/middlewares"
	"loki/internal/config/router"
	"loki/internal/config/server"
	"loki/internal/config/telemetry"
	"loki/pkg/jwt"
)

var Module = fx.Options(
	logger.Module,

	authentication.Module,
	controllers.Module,
	repositories.Module,
	jwt.Module,
	services.Module,
	workers.Module,

	rpcs.Module,

	middlewares.Module,

	server.Module,
	router.Module,
	telemetry.Module,

	fx.Invoke(registerWebServer),
	fx.Invoke(registerGrpcServer),
	fx.Invoke(registerWorkers),
	fx.Invoke(registerTelemetry),
)

func registerWebServer(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	server server.WebServer,
	log *logger.Logger,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msgf("Starting web server in %s environment at %s", cfg.AppEnv, cfg.AppAddr)
			go func() {
				if err := server.Run(); err != nil && err != http.ErrServerClosed {
					log.Error().Err(err).Msg("Web server start failed")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			log.Info().Msg("Shutting down web server...")
			return server.Shutdown(shutdownCtx)
		},
	})
}

func registerGrpcServer(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	server server.GrpcServer,
	log *logger.Logger,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msgf("Starting gRPC server in %s environment at %s", cfg.AppEnv, cfg.GrpcAddr)
			go func() {
				if err := server.Run(); err != nil {
					log.Error().Err(err).Msg("gRPC server start failed")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			log.Info().Msg("Shutting down gRPC server...")
			return server.Shutdown(shutdownCtx)
		},
	})
}

func registerWorkers(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	smartId smartid.Worker,
	mobileId mobileid.Worker,
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
