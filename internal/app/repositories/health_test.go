package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki/internal/app/repositories/postgres"
	"loki/internal/config"
)

func Test_HealthRepository_Ping(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	healthRepository := NewHealthRepository(client)

	tests := []struct {
		name     string
		expected error
	}{
		{
			name:     "Success",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := healthRepository.Ping(ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}
