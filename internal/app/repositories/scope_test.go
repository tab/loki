package repositories

import (
	"context"
	"loki/internal/app/models"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
	"loki/internal/config"
)

func Test_ScopeRepository_FindByName(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	repo := NewScopeRepository(client)

	tests := []struct {
		name     string
		param    string
		expected *models.Scope
		error    bool
	}{
		{
			name:  "Success",
			param: models.SelfServiceType,
			expected: &models.Scope{
				ID:   uuid.MustParse("10000000-1000-1000-2000-000000000001"),
				Name: models.SelfServiceType,
			},
			error: false,
		},
		{
			name:     "Scope not found",
			param:    "unknown",
			expected: nil,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindByName(ctx, tt.param)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Name, result.Name)
			}
		})
	}
}

func Test_ScopeRepository_CreateUserScope(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	userRepository := NewUserRepository(client)
	scopeRepository := NewScopeRepository(client)

	account, err := userRepository.Create(ctx, db.CreateUserParams{
		IdentityNumber: "PNOEE-30303039914",
		PersonalCode:   "30303039914",
		FirstName:      "TESTNUMBER",
		LastName:       "OK",
	})
	assert.NoError(t, err)

	tests := []struct {
		name     string
		params   db.CreateUserScopeParams
		expected error
	}{
		{
			name: "Success",
			params: db.CreateUserScopeParams{
				UserID:  account.ID,
				ScopeID: uuid.MustParse("10000000-1000-1000-2000-000000000001"),
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := scopeRepository.CreateUserScope(ctx, tt.params)
			assert.NoError(t, err)
		})
	}
}

func Test_ScopeRepository_FindByUserId(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	userRepository := NewUserRepository(client)
	scopeRepository := NewScopeRepository(client)

	account, err := userRepository.Create(ctx, db.CreateUserParams{
		IdentityNumber: "PNOEE-30303039914",
		PersonalCode:   "30303039914",
		FirstName:      "TESTNUMBER",
		LastName:       "OK",
	})
	assert.NoError(t, err)

	tests := []struct {
		name     string
		param    uuid.UUID
		expected []models.Scope
	}{
		{
			name:  "Success",
			param: account.ID,
			expected: []models.Scope{
				{ID: uuid.MustParse("10000000-1000-1000-2000-000000000001"), Name: models.SelfServiceType},
			},
		},
		{
			name:     "User not found",
			param:    uuid.Nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := scopeRepository.FindByUserId(ctx, tt.param)
			assert.NoError(t, err)

			assert.Equal(t, len(tt.expected), len(results))
			for i, result := range results {
				assert.Equal(t, tt.expected[i].ID, result.ID)
				assert.Equal(t, tt.expected[i].Name, result.Name)
			}
		})
	}
}
