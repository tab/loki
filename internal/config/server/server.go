package server

import (
	"context"
	"net/http"
	"time"

	"loki/internal/config"
)

type WebServer interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type webServer struct {
	httpServer *http.Server
}

func NewWebServer(cfg *config.Config, appRouter http.Handler) WebServer {
	return &webServer{
		httpServer: &http.Server{
			Addr:         cfg.AppAddr,
			Handler:      appRouter,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (s *webServer) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *webServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
