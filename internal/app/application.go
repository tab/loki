package app

import (
	"context"
	"net/http"
	"time"

	"loki/internal/config"
	"loki/internal/config/router"
	"loki/internal/config/server"
	"loki/pkg/logger"
)

type Application struct {
	cfg    *config.Config
	log    *logger.Logger
	server server.Server
}

func NewApplication(_ context.Context) (*Application, error) {
	appLogger := logger.NewLogger()
	cfg := config.LoadConfig(appLogger)

	appRouter := router.NewRouter(cfg)
	appServer := server.NewServer(cfg, appRouter)

	return &Application{
		cfg:    cfg,
		log:    appLogger,
		server: appServer,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	serverErrors := make(chan error, 1)
	go func() {
		if err := a.server.Run(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	a.log.Info().Msgf("Application starting in %s", a.cfg.AppEnv)
	a.log.Info().Msgf("Listening on %s", a.cfg.AppAddr)

	select {
	case <-ctx.Done():
		a.log.Info().Msg("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		a.log.Info().Msg("Server gracefully stopped")
		return nil
	case err := <-serverErrors:
		return err
	}
}
