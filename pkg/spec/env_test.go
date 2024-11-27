package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if os.Getenv("GO_ENV") == "ci" {
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
}

func Test_LoadEnv(t *testing.T) {
	type env struct {
		AppEnv  string
		AppAddr string
	}

	tests := []struct {
		name     string
		expected env
	}{
		{
			name: "Success",
			expected: env{
				AppEnv:  "test",
				AppAddr: "localhost:8080",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LoadEnv()
			assert.NoError(t, err)

			hash := []struct{ key, value string }{
				{"GO_ENV", tt.expected.AppEnv},
				{"APP_ADDRESS", tt.expected.AppAddr},
			}

			for _, h := range hash {
				envValue := os.Getenv(h.key)
				assert.Equal(t, h.value, envValue)
			}
		})
	}
}
