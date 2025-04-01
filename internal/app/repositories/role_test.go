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

func Test_RoleRepository_List(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	roleRepository := NewRoleRepository(client)

	tests := []struct {
		name     string
		limit    uint64
		offset   uint64
		total    uint64
		expected []models.Role
	}{
		{
			name:   "List roles",
			limit:  10,
			offset: 0,
			total:  3,
			expected: []models.Role{
				{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
				},
				{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000002"),
					Name:        models.ManagerRoleType,
					Description: "Manager role",
				},
				{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000003"),
					Name:        models.UserRoleType,
					Description: "User role",
				},
			},
		},
		{
			name:   "List with offset",
			limit:  2,
			offset: 1,
			total:  3,
			expected: []models.Role{
				{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000002"),
					Name:        models.ManagerRoleType,
					Description: "Manager role",
				},
				{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000003"),
					Name:        models.UserRoleType,
					Description: "User role",
				},
			},
		},
		{
			name:     "List with zero limit",
			limit:    0,
			offset:   0,
			total:    0,
			expected: []models.Role{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, total, err := roleRepository.List(ctx, tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(results))
			assert.Equal(t, tt.total, total)
		})
	}
}

func Test_RoleRepository_Create(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	roleRepository := NewRoleRepository(client)

	tests := []struct {
		name   string
		params db.CreateRoleParams
		error  bool
	}{
		{
			name: "Create valid role",
			params: db.CreateRoleParams{
				Name:        "Developer",
				Description: "Developer role",
			},
			error: false,
		},
		{
			name: "Create existing role",
			params: db.CreateRoleParams{
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := roleRepository.Create(ctx, tt.params)

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

func Test_RoleRepository_Update(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	roleRepository := NewRoleRepository(client)

	existingRole, err := roleRepository.Create(ctx, db.CreateRoleParams{
		Name:        "Tester",
		Description: "Tester role",
	})
	assert.NoError(t, err)
	assert.NotNil(t, existingRole)

	tests := []struct {
		name   string
		params db.UpdateRoleParams
		error  bool
	}{
		{
			name: "Update existing role",
			params: db.UpdateRoleParams{
				ID:          existingRole.ID,
				Name:        "Tester Updated",
				Description: "Updated tester role",
			},
			error: false,
		},
		{
			name: "Update non-existing role",
			params: db.UpdateRoleParams{
				ID:          uuid.New(),
				Name:        "Nonexistent",
				Description: "Does not exist",
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := roleRepository.Update(ctx, tt.params)

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

func Test_RoleRepository_FindById(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	repo := NewRoleRepository(client)

	tests := []struct {
		name     string
		param    uuid.UUID
		expected *models.Role
		error    bool
	}{
		{
			name:  "Find existing role",
			param: uuid.MustParse("10000000-1000-1000-1000-000000000001"),
			expected: &models.Role{
				ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			error: false,
		},
		{
			name:     "Find non-existing role",
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
			result, err := repo.FindById(ctx, tt.param)

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

func Test_RoleRepository_Delete(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	roleRepository := NewRoleRepository(client)

	existingRole, err := roleRepository.Create(ctx, db.CreateRoleParams{
		Name:        "TempRole",
		Description: "Temporary role",
	})
	assert.NoError(t, err)

	tests := []struct {
		name  string
		param uuid.UUID
	}{
		{
			name:  "Delete existing role",
			param: existingRole.ID,
		},
		{
			name:  "Delete non-existing role",
			param: uuid.New(),
		},
		{
			name:  "Delete with nil UUID",
			param: uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := roleRepository.Delete(ctx, tt.param)

			assert.NoError(t, err)
			assert.Equal(t, true, result)
		})
	}
}

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
				ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			error: false,
		},
		{
			name:  "Success (manager)",
			param: models.ManagerRoleType,
			expected: &models.Role{
				ID:          uuid.MustParse("10000000-1000-1000-1000-000000000002"),
				Name:        models.ManagerRoleType,
				Description: "Manager role",
			},
			error: false,
		},
		{
			name:  "Success (user)",
			param: models.UserRoleType,
			expected: &models.Role{
				ID:          uuid.MustParse("10000000-1000-1000-1000-000000000003"),
				Name:        models.UserRoleType,
				Description: "User role",
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

func Test_RoleRepository_FindRoleDetailsById(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	repo := NewRoleRepository(client)

	tests := []struct {
		name     string
		param    uuid.UUID
		expected *models.Role
		error    bool
	}{
		{
			name:  "Find existing role details",
			param: uuid.MustParse("10000000-1000-1000-1000-000000000001"),
			expected: &models.Role{
				ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
				Name:        models.AdminRoleType,
				Description: "Admin role",
				PermissionIDs: []uuid.UUID{
					uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					uuid.MustParse("10000000-1000-1000-3000-000000000002"),
					uuid.MustParse("10000000-1000-1000-3000-000000000003"),
					uuid.MustParse("10000000-1000-1000-3000-000000000004"),
				},
			},
			error: false,
		},
		{
			name:     "Find role details for non-existing role",
			param:    uuid.New(),
			expected: nil,
			error:    true,
		},
		{
			name:     "Find role details with nil UUID",
			param:    uuid.Nil,
			expected: nil,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindRoleDetailsById(ctx, tt.param)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.Description, result.Description)
				assert.ElementsMatch(t, tt.expected.PermissionIDs, result.PermissionIDs)
			}
		})
	}
}
