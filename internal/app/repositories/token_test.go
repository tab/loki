package repositories

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
	"loki/internal/config"
)

func Test_TokenRepository_List(t *testing.T) {
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

	_, err = tokenRepository.Create(ctx, db.CreateTokensParams{
		UserID:            account.ID,
		AccessTokenValue:  "aaa.bbb.ccc",
		RefreshTokenValue: "ccc.ccc.ccc",
	})
	assert.NoError(t, err)

	_, err = tokenRepository.Create(ctx, db.CreateTokensParams{
		UserID:            account.ID,
		AccessTokenValue:  "aaa.bbb.ddd",
		RefreshTokenValue: "ddd.ddd.ddd",
	})
	assert.NoError(t, err)

	tests := []struct {
		name     string
		limit    uint64
		offset   uint64
		total    uint64
		expected []models.Token
	}{
		{
			name:   "List tokens",
			limit:  10,
			offset: 0,
			total:  4,
			expected: []models.Token{
				{
					UserId: account.ID,
					Type:   models.AccessTokenType,
					Value:  "aaa.bbb.ccc",
				},
				{
					UserId: account.ID,
					Type:   models.RefreshTokenType,
					Value:  "ccc.ccc.ccc",
				},
				{
					UserId: account.ID,
					Type:   models.AccessTokenType,
					Value:  "aaa.bbb.ddd",
				},
				{
					UserId: account.ID,
					Type:   models.RefreshTokenType,
					Value:  "ddd.ddd.ddd",
				},
			},
		},
		{
			name:   "List with offset",
			limit:  2,
			offset: 2,
			total:  4,
			expected: []models.Token{
				{
					UserId: account.ID,
					Type:   models.AccessTokenType,
					Value:  "aaa.bbb.ddd",
				},
				{
					UserId: account.ID,
					Type:   models.RefreshTokenType,
					Value:  "ddd.ddd.ddd",
				},
			},
		},
		{
			name:     "List with zero limit",
			limit:    0,
			offset:   0,
			total:    0,
			expected: []models.Token{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, total, err := tokenRepository.List(ctx, tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(results))
			assert.Equal(t, tt.total, total)
		})
	}
}

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

func Test_TokenRepository_FindById(t *testing.T) {
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

	existingTokens, err := tokenRepository.Create(ctx, db.CreateTokensParams{
		UserID:            account.ID,
		AccessTokenValue:  "access-token-123",
		RefreshTokenValue: "refresh-token-123",
	})
	assert.NoError(t, err)
	assert.NotNil(t, existingTokens)
	accessToken := existingTokens[0]

	tests := []struct {
		name     string
		param    uuid.UUID
		expected *models.Token
		error    bool
	}{
		{
			name:  "Find existing token",
			param: existingTokens[0].ID,
			expected: &models.Token{
				ID:        accessToken.ID,
				UserId:    accessToken.UserId,
				Type:      accessToken.Type,
				Value:     accessToken.Value,
				ExpiresAt: accessToken.ExpiresAt,
			},
			error: false,
		},
		{
			name:     "Find non-existing token",
			param:    uuid.New(),
			expected: nil,
			error:    true,
		},
		{
			name:     "Find with nil UUID",
			param:    uuid.Nil,
			expected: nil,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tokenRepository.FindById(ctx, tt.param)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Type, result.Type)
				assert.Equal(t, tt.expected.Value, result.Value)
				assert.WithinDuration(t, tt.expected.ExpiresAt, result.ExpiresAt, time.Second)
			}
		})
	}
}

func Test_TokenRepository_Delete(t *testing.T) {
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

	existingTokens, err := tokenRepository.Create(ctx, db.CreateTokensParams{
		UserID:            account.ID,
		AccessTokenValue:  "access-token-123",
		RefreshTokenValue: "refresh-token-123",
	})
	assert.NoError(t, err)
	assert.NotNil(t, existingTokens)

	accessToken := existingTokens[0]
	refreshToken := existingTokens[1]

	tests := []struct {
		name  string
		param uuid.UUID
	}{
		{
			name:  "Delete existing access token",
			param: accessToken.ID,
		},
		{
			name:  "Delete existing refresh token",
			param: refreshToken.ID,
		},
		{
			name:  "Delete non-existing token",
			param: uuid.New(),
		},
		{
			name:  "Delete with nil UUID",
			param: uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tokenRepository.Delete(ctx, tt.param)

			assert.NoError(t, err)
			assert.Equal(t, true, result)
		})
	}
}
