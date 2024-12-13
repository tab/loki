package app

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/fx"

	"loki/internal/app/controllers"
	"loki/internal/app/repositories"
	"loki/internal/app/services"
	"loki/internal/config"
	"loki/internal/config/middlewares"
	"loki/internal/config/router"
	"loki/internal/config/server"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

var Module = fx.Options(
	logger.Module,

	controllers.Module,
	repositories.Module,
	jwt.Module,
	services.Module,

	middlewares.Module,

	server.Module,
	router.Module,

	fx.Invoke(registerHooks),
)

func registerHooks(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	server server.Server,
	smartId services.SmartIdWorker,
	mobileId services.MobileIdWorker,
	log *logger.Logger,
) {
	var backgroundCtx context.Context
	var backgroundCtxCancel context.CancelFunc

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msgf("Starting server at %s", cfg.AppAddr)

			go func() {
				if err := server.Run(); err != nil && err != http.ErrServerClosed {
					log.Error().Err(err).Msg("Server failed")
				}
			}()

			backgroundCtx, backgroundCtxCancel = context.WithCancel(context.Background())
			smartId.Start(backgroundCtx)
			mobileId.Start(backgroundCtx)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			backgroundCtxCancel()
			smartId.Stop()
			mobileId.Stop()

			log.Info().Msg("Shutting down server...")

			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			return server.Shutdown(shutdownCtx)
		},
	})
}
