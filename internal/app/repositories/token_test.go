package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
	"loki/internal/config"
)

func Test_TokenRepository_Create(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	userRepository := NewUserRepository(client)
	tokenRepository := NewTokenRepository(client)

	account, err := userRepository.Create(ctx, db.CreateUserParams{
		IdentityNumber: "PNOEE-30303039914",
		PersonalCode:   "30303039914",
		FirstName:      "TESTNUMBER",
		LastName:       "OK",
	})
	assert.NoError(t, err)

	tests := []struct {
		name     string
		params   db.CreateTokensParams
		expected []models.Token
		error    bool
	}{
		{
			name: "Success",
			params: db.CreateTokensParams{
				UserID:            account.ID,
				AccessTokenValue:  "aaa.bbb.ccc",
				RefreshTokenValue: "ddd.eee.fff",
			},
			expected: []models.Token{
				{
					UserId: account.ID,
				},
				{
					UserId: account.ID,
				},
			},
			error: false,
		},
		{
			name: "Error",
			params: db.CreateTokensParams{
				UserID:            uuid.Nil,
				AccessTokenValue:  "",
				RefreshTokenValue: "",
			},
			expected: nil,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := tokenRepository.Create(ctx, tt.params)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(results))
			}
		})
	}
}
