package spec

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GenerateCertificates(t *testing.T) {
	certDir := GenerateCertificates(t)

	_, err := os.Stat(certDir)
	assert.NoError(t, err)

	expectedFiles := []string{
		CaFile,
		ServerCertFile,
		ServerKeyFile,
		ClientCertFile,
		ClientKeyFile,
	}

	for _, filename := range expectedFiles {
		path := filepath.Join(certDir, filename)
		_, err = os.Stat(path)
		assert.NoError(t, err)
	}

	caCertPath := filepath.Join(certDir, CaFile)
	caCertData, err := os.ReadFile(caCertPath)
	require.NoError(t, err)

	caCertBlock, _ := pem.Decode(caCertData)
	require.NotNil(t, caCertBlock)
	assert.Equal(t, "CERTIFICATE", caCertBlock.Type)

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	require.NoError(t, err)

	assert.True(t, caCert.IsCA)
	assert.Equal(t, "Test CA", caCert.Subject.CommonName)

	serverCertPath := filepath.Join(certDir, ServerCertFile)
	serverCertData, err := os.ReadFile(serverCertPath)
	require.NoError(t, err)

	serverCertBlock, _ := pem.Decode(serverCertData)
	require.NotNil(t, serverCertBlock)
	assert.Equal(t, "CERTIFICATE", serverCertBlock.Type)

	serverCert, err := x509.ParseCertificate(serverCertBlock.Bytes)
	require.NoError(t, err)

	assert.Equal(t, "localhost", serverCert.Subject.CommonName)
	assert.Contains(t, serverCert.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	assert.Contains(t, serverCert.DNSNames, "localhost")

	clientCertPath := filepath.Join(certDir, ClientCertFile)
	clientCertData, err := os.ReadFile(clientCertPath)
	require.NoError(t, err)

	clientCertBlock, _ := pem.Decode(clientCertData)
	require.NotNil(t, clientCertBlock)
	assert.Equal(t, "CERTIFICATE", clientCertBlock.Type)

	clientCert, err := x509.ParseCertificate(clientCertBlock.Bytes)
	require.NoError(t, err)

	assert.Equal(t, "client", clientCert.Subject.CommonName)
	assert.Contains(t, clientCert.ExtKeyUsage, x509.ExtKeyUsageClientAuth)

	roots := x509.NewCertPool()
	roots.AddCert(caCert)

	opts := x509.VerifyOptions{
		Roots: roots,
	}

	serverOpts := opts
	serverOpts.KeyUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	_, err = serverCert.Verify(serverOpts)
	assert.NoError(t, err)

	clientOpts := opts
	clientOpts.KeyUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	_, err = clientCert.Verify(clientOpts)
	assert.NoError(t, err)

	serverKeyPath := filepath.Join(certDir, ServerKeyFile)
	serverKeyData, err := os.ReadFile(serverKeyPath)
	require.NoError(t, err)

	serverKeyBlock, _ := pem.Decode(serverKeyData)
	require.NotNil(t, serverKeyBlock)
	assert.Equal(t, "RSA PRIVATE KEY", serverKeyBlock.Type)

	clientKeyPath := filepath.Join(certDir, ClientKeyFile)
	clientKeyData, err := os.ReadFile(clientKeyPath)
	require.NoError(t, err)

	clientKeyBlock, _ := pem.Decode(clientKeyData)
	require.NotNil(t, clientKeyBlock)
	assert.Equal(t, "RSA PRIVATE KEY", clientKeyBlock.Type)
}
