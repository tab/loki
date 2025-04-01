package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"loki/internal/app/rpcs"
	"loki/internal/app/rpcs/interceptors"
	"loki/internal/config"
	"loki/pkg/logger"
)

const (
	CaFile   = "ca.pem"
	CertFile = "server.pem"
	KeyFile  = "server.key"

	MaxConnectionIdle     = 5 * time.Minute
	MaxConnectionAge      = 5 * time.Minute
	MaxConnectionAgeGrace = 1 * time.Minute
	KeepaliveTime         = 5 * time.Second
	KeepaliveTimeout      = 1 * time.Second
)

type GrpcServer interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type grpcServer struct {
	cfg      *config.Config
	server   *grpc.Server
	registry *rpcs.Registry
	log      *logger.Logger
}

func NewGrpcServer(
	cfg *config.Config,
	authInterceptor interceptors.AuthenticationInterceptor,
	registry *rpcs.Registry,
	log *logger.Logger,
) GrpcServer {
	tlsConfig, err := setupTLS(cfg, log)
	if err != nil {
		return nil
	}

	options := keepalive.ServerParameters{
		MaxConnectionIdle:     MaxConnectionIdle,
		MaxConnectionAge:      MaxConnectionAge,
		MaxConnectionAgeGrace: MaxConnectionAgeGrace,
		Time:                  KeepaliveTime,
		Timeout:               KeepaliveTimeout,
	}

	server := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.KeepaliveParams(options),
		grpc.UnaryInterceptor(
			auth.UnaryServerInterceptor(authInterceptor.Authenticate),
		),
		grpc.StreamInterceptor(
			auth.StreamServerInterceptor(authInterceptor.Authenticate),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	registry.RegisterAll(server)

	return &grpcServer{
		cfg:      cfg,
		server:   server,
		registry: registry,
		log:      log,
	}
}

func (s *grpcServer) Run() error {
	listener, err := net.Listen("tcp", s.cfg.GrpcAddr)
	if err != nil {
		return err
	}

	return s.server.Serve(listener)
}

func (s *grpcServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	case <-done:
		return nil
	}
}

func setupTLS(cfg *config.Config, log *logger.Logger) (*tls.Config, error) {
	caCert, err := os.ReadFile(filepath.Join(cfg.CertPath, CaFile))
	if err != nil {
		log.Error().Err(err).Msg("Failed to load CA certificate")
		return nil, err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(
		filepath.Join(cfg.CertPath, CertFile),
		filepath.Join(cfg.CertPath, KeyFile),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load server certificate and private key")
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caPool,
		MinVersion:   tls.VersionTLS13,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}
