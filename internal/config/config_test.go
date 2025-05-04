package config

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki/pkg/spec"
)

func TestMain(m *testing.M) {
	if err := spec.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	if os.Getenv("GO_ENV") == "ci" {
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
}

func Test_LoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		env      map[string]string
		expected *Config
	}{
		{
			name: "Success",
			args: []string{},
			env:  map[string]string{},
			expected: &Config{
				AppEnv:      "test",
				AppAddr:     "0.0.0.0:8080",
				GrpcAddr:    "0.0.0.0:50051",
				ClientURL:   "http://localhost:3000",
				CertPath:    "./certs",
				DatabaseDSN: "postgres://postgres:postgres@localhost:5432/loki-test?sslmode=disable",
				RedisURI:    "redis://localhost:6379/1",
				SmartId: SmartId{
					BaseURL:          "https://sid.demo.sk.ee/smart-id-rp/v2",
					RelyingPartyName: "DEMO",
					RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
					Text:             "Enter PIN1",
				},
				MobileId: MobileId{
					BaseURL:          "https://tsp.demo.sk.ee/mid-api",
					RelyingPartyName: "DEMO",
					RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
					Text:             "Enter PIN1",
					TextFormat:       "GSM-7",
					Language:         "ENG",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			result := LoadConfig()

			assert.Equal(t, tt.expected.AppEnv, result.AppEnv)
			assert.Equal(t, tt.expected.AppAddr, result.AppAddr)
			assert.Equal(t, tt.expected.GrpcAddr, result.GrpcAddr)
			assert.Equal(t, tt.expected.ClientURL, result.ClientURL)
			assert.Equal(t, tt.expected.CertPath, result.CertPath)
			assert.Equal(t, tt.expected.DatabaseDSN, result.DatabaseDSN)
			assert.Equal(t, tt.expected.RedisURI, result.RedisURI)
			assert.Equal(t, tt.expected.SmartId, result.SmartId)
			assert.Equal(t, tt.expected.MobileId, result.MobileId)

			t.Cleanup(func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			})
		})
	}
}
