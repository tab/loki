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

func Test_RoleRepository_FindByName(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	repo := NewRoleRepository(client)

	tests := []struct {
		name     string
		param    string
		expected *models.Role
		error    bool
	}{
		{
			name:  "Success (admin)",
			param: models.AdminRoleType,
			expected: &models.Role{
				ID:   uuid.MustParse("10000000-1000-1000-1000-000000000001"),
				Name: models.AdminRoleType,
			},
			error: false,
		},
		{
			name:  "Success (manager)",
			param: models.ManagerRoleType,
			expected: &models.Role{
				ID:   uuid.MustParse("10000000-1000-1000-1000-000000000002"),
				Name: models.ManagerRoleType,
			},
			error: false,
		},
		{
			name:  "Success (user)",
			param: models.UserRoleType,
			expected: &models.Role{
				ID:   uuid.MustParse("10000000-1000-1000-1000-000000000003"),
				Name: models.UserRoleType,
			},
			error: false,
		},
		{
			name:     "Role not found",
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

func Test_RoleRepository_CreateUserRole(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	userRepository := NewUserRepository(client)
	roleRepository := NewRoleRepository(client)

	account, err := userRepository.Create(ctx, db.CreateUserParams{
		IdentityNumber: "PNOEE-30303039914",
		PersonalCode:   "30303039914",
		FirstName:      "TESTNUMBER",
		LastName:       "OK",
	})
	assert.NoError(t, err)

	tests := []struct {
		name     string
		params   db.CreateUserRoleParams
		expected error
	}{
		{
			name: "Success",
			params: db.CreateUserRoleParams{
				UserID: account.ID,
				RoleID: uuid.MustParse("10000000-1000-1000-1000-000000000003"),
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := roleRepository.CreateUserRole(ctx, tt.params)
			assert.NoError(t, err)
		})
	}
}

func Test_RoleRepository_FindByUserId(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	userRepository := NewUserRepository(client)
	roleRepository := NewRoleRepository(client)

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
		expected []models.Role
	}{
		{
			name:  "Success",
			param: account.ID,
			expected: []models.Role{
				{ID: uuid.MustParse("10000000-1000-1000-1000-000000000003"), Name: models.UserRoleType},
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
			results, err := roleRepository.FindByUserId(ctx, tt.param)
			assert.NoError(t, err)

			assert.Equal(t, len(tt.expected), len(results))
			for i, result := range results {
				assert.Equal(t, tt.expected[i].ID, result.ID)
				assert.Equal(t, tt.expected[i].Name, result.Name)
			}
		})
	}
}
