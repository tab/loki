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

func Test_PermissionRepository_List(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	permissionRepository := NewPermissionRepository(client)

	tests := []struct {
		name     string
		limit    uint64
		offset   uint64
		total    uint64
		expected []models.Permission
	}{
		{
			name:   "List permissions",
			limit:  10,
			offset: 0,
			total:  4,
			expected: []models.Permission{
				{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				},
				{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000002"),
					Name:        "write:self",
					Description: "Write own data",
				},
				{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000003"),
					Name:        "read:users",
					Description: "Read user data",
				},
				{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000004"),
					Name:        "write:users",
					Description: "Write user data",
				},
			},
		},
		{
			name:   "List with offset",
			limit:  2,
			offset: 2,
			total:  4,
			expected: []models.Permission{
				{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000003"),
					Name:        "read:users",
					Description: "Read user data",
				},
				{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000004"),
					Name:        "write:users",
					Description: "Write user data",
				},
			},
		},
		{
			name:     "List with zero limit",
			limit:    0,
			offset:   0,
			total:    0,
			expected: []models.Permission{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, total, err := permissionRepository.List(ctx, tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(results))
			assert.Equal(t, tt.total, total)
		})
	}
}

func Test_PermissionRepository_Create(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	permissionRepository := NewPermissionRepository(client)

	tests := []struct {
		name   string
		params db.CreatePermissionParams
		error  bool
	}{
		{
			name: "Create valid permission",
			params: db.CreatePermissionParams{
				Name:        "delete:self",
				Description: "Delete own data",
			},
			error: false,
		},
		{
			name: "Create existing permission",
			params: db.CreatePermissionParams{
				Name:        "read:self",
				Description: "Read own data",
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := permissionRepository.Create(ctx, tt.params)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.params.Name, result.Name)
				assert.Equal(t, tt.params.Description, result.Description)
			}
		})
	}
}

func Test_PermissionRepository_Update(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	permissionRepository := NewPermissionRepository(client)

	existingPermission, err := permissionRepository.Create(ctx, db.CreatePermissionParams{
		Name:        "read:existing",
		Description: "Existing permission",
	})
	assert.NoError(t, err)
	assert.NotNil(t, existingPermission)

	tests := []struct {
		name   string
		params db.UpdatePermissionParams
		error  bool
	}{
		{
			name: "Update existing permission",
			params: db.UpdatePermissionParams{
				ID:          existingPermission.ID,
				Name:        "read:updated",
				Description: "Updated description",
			},
			error: false,
		},
		{
			name: "Update non-existing permission",
			params: db.UpdatePermissionParams{
				ID:          uuid.New(),
				Name:        "nonexistent",
				Description: "Does not exist",
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := permissionRepository.Update(ctx, tt.params)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.params.Name, result.Name)
				assert.Equal(t, tt.params.Description, result.Description)
			}
		})
	}
}

func Test_PermissionRepository_FindById(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	permissionRepository := NewPermissionRepository(client)

	tests := []struct {
		name     string
		param    uuid.UUID
		expected *models.Permission
		error    bool
	}{
		{
			name:  "Find existing permission",
			param: uuid.MustParse("10000000-1000-1000-3000-000000000001"),
			expected: &models.Permission{
				ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
				Name:        "read:self",
				Description: "Read own data",
			},
			error: false,
		},
		{
			name:     "Find non-existing permission",
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
			result, err := permissionRepository.FindById(ctx, tt.param)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.Description, result.Description)
			}
		})
	}
}

func Test_PermissionRepository_Delete(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	permissionRepository := NewPermissionRepository(client)

	existingPermission, err := permissionRepository.Create(ctx, db.CreatePermissionParams{
		Name:        "temp:delete",
		Description: "Temporary permission",
	})
	assert.NoError(t, err)

	tests := []struct {
		name  string
		param uuid.UUID
	}{
		{
			name:  "Delete existing permission",
			param: existingPermission.ID,
		},
		{
			name:  "Delete non-existing permission",
			param: uuid.New(),
		},
		{
			name:  "Delete with nil UUID",
			param: uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := permissionRepository.Delete(ctx, tt.param)

			assert.NoError(t, err)
			assert.Equal(t, true, result)
		})
	}
}

func Test_PermissionRepository_FindByUserId(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	userRepository := NewUserRepository(client)
	permissionRepository := NewPermissionRepository(client)

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
		expected []models.Permission
	}{
		{
			name:  "Success",
			param: account.ID,
			expected: []models.Permission{
				{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				},
				{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000002"),
					Name:        "write:self",
					Description: "Write own data",
				},
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
			results, err := permissionRepository.FindByUserId(ctx, tt.param)
			assert.NoError(t, err)

			assert.Equal(t, len(tt.expected), len(results))
			for i, result := range results {
				assert.Equal(t, tt.expected[i].ID, result.ID)
				assert.Equal(t, tt.expected[i].Name, result.Name)
			}
		})
	}
}
