package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"loki/internal/app/rpcs"
	"loki/internal/app/rpcs/interceptors"
	"loki/internal/config"
	"loki/internal/config/logger"
)

func Test_NewGrpcServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	certDir := generateTestCertificates(t)

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		GrpcAddr: "localhost:50051",
		CertPath: certDir,
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	authInterceptor := interceptors.NewMockAuthenticationInterceptor(ctrl)
	traceInterceptor := interceptors.NewMockTraceInterceptor(ctrl)
	loggerInterceptor := interceptors.NewMockLoggerInterceptor(ctrl)

	authInterceptor.EXPECT().Authenticate(gomock.Any()).AnyTimes()
	traceInterceptor.EXPECT().Trace().AnyTimes()
	loggerInterceptor.EXPECT().Log().AnyTimes()

	registry := &rpcs.Registry{}

	srv := NewGrpcServer(cfg, registry, authInterceptor, traceInterceptor, loggerInterceptor, log)
	assert.NotNil(t, srv)

	s, ok := srv.(*grpcServer)
	assert.True(t, ok)
	assert.Equal(t, cfg, s.cfg)
	assert.NotNil(t, s.server)
}

func Test_GrpcServer_RunAndShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	certDir := generateTestCertificates(t)

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		GrpcAddr: "localhost:50051",
		CertPath: certDir,
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	authInterceptor := interceptors.NewMockAuthenticationInterceptor(ctrl)
	traceInterceptor := interceptors.NewMockTraceInterceptor(ctrl)
	loggerInterceptor := interceptors.NewMockLoggerInterceptor(ctrl)

	authInterceptor.EXPECT().Authenticate(gomock.Any()).AnyTimes()
	traceInterceptor.EXPECT().Trace().AnyTimes()
	loggerInterceptor.EXPECT().Log().AnyTimes()

	registry := &rpcs.Registry{}

	srv := NewGrpcServer(cfg, registry, authInterceptor, traceInterceptor, loggerInterceptor, log)
	assert.NotNil(t, srv)

	runErrCh := make(chan error, 1)
	go func() {
		err := srv.Run()
		runErrCh <- err
	}()

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	assert.NoError(t, err)

	err = <-runErrCh
	assert.NoError(t, err)
}

func generateTestCertificates(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "tls-test-*")
	require.NoError(t, err)

	t.Cleanup(func() { os.RemoveAll(tempDir) })

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	caTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"ACME CA"}, CommonName: "ACME CA"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caCertDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	require.NoError(t, err)
	caCertPath := filepath.Join(tempDir, "ca.pem")
	caCertFile, err := os.Create(caCertPath)
	require.NoError(t, err)
	defer caCertFile.Close()
	err = pem.Encode(caCertFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertDER,
	})
	require.NoError(t, err)

	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	serverTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{Organization: []string{"ACME"}, CommonName: "localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost", "127.0.0.1", "0.0.0.0"},
	}
	serverCertDER, err := x509.CreateCertificate(rand.Reader, &serverTemplate, &caTemplate, &serverKey.PublicKey, caKey)
	require.NoError(t, err)
	serverCertPath := filepath.Join(tempDir, "server.pem")
	serverCertFile, err := os.Create(serverCertPath)
	require.NoError(t, err)
	defer serverCertFile.Close()
	err = pem.Encode(serverCertFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertDER,
	})
	require.NoError(t, err)

	serverKeyPath := filepath.Join(tempDir, "server.key")
	serverKeyFile, err := os.Create(serverKeyPath)
	require.NoError(t, err)
	defer serverKeyFile.Close()
	err = pem.Encode(serverKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverKey),
	})
	require.NoError(t, err)

	return tempDir
}
