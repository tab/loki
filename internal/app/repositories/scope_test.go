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

func Test_ScopeRepository_List(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	scopeRepository := NewScopeRepository(client)

	tests := []struct {
		name     string
		limit    uint64
		offset   uint64
		total    uint64
		expected []models.Scope
	}{
		{
			name:   "List scopes",
			limit:  10,
			offset: 0,
			total:  2,
			expected: []models.Scope{
				{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        models.SsoServiceType,
					Description: "SSO-service scope",
				},
				{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000002"),
					Name:        models.SelfServiceType,
					Description: "Self-service scope",
				},
			},
		},
		{
			name:   "List with offset",
			limit:  1,
			offset: 1,
			total:  2,
			expected: []models.Scope{
				{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000002"),
					Name:        models.SelfServiceType,
					Description: "Self-service scope",
				},
			},
		},
		{
			name:     "List with zero limit",
			limit:    0,
			offset:   0,
			total:    0,
			expected: []models.Scope{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, total, err := scopeRepository.List(ctx, tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(results))
			assert.Equal(t, tt.total, total)
		})
	}
}

func Test_ScopeRepository_Create(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	scopeRepository := NewScopeRepository(client)

	tests := []struct {
		name   string
		params db.CreateScopeParams
		error  bool
	}{
		{
			name: "Create valid scope",
			params: db.CreateScopeParams{
				Name:        "Admin-service",
				Description: "Admin service scope",
			},
			error: false,
		},
		{
			name: "Create existing scope",
			params: db.CreateScopeParams{
				Name:        models.SelfServiceType,
				Description: "Self-service scope",
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := scopeRepository.Create(ctx, tt.params)

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

func Test_ScopeRepository_Update(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	scopeRepository := NewScopeRepository(client)

	existingScope, err := scopeRepository.Create(ctx, db.CreateScopeParams{
		Name:        "existing-service",
		Description: "Existing service scope",
	})
	assert.NoError(t, err)
	assert.NotNil(t, existingScope)

	tests := []struct {
		name   string
		params db.UpdateScopeParams
		error  bool
	}{
		{
			name: "Update existing scope",
			params: db.UpdateScopeParams{
				ID:          existingScope.ID,
				Name:        "existing-service-updated",
				Description: "Updated description",
			},
			error: false,
		},
		{
			name: "Update non-existing scope",
			params: db.UpdateScopeParams{
				ID:          uuid.New(),
				Name:        "nonexistent",
				Description: "Does not exist",
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := scopeRepository.Update(ctx, tt.params)

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

func Test_ScopeRepository_FindById(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	scopeRepository := NewScopeRepository(client)

	tests := []struct {
		name     string
		param    uuid.UUID
		expected *models.Scope
		error    bool
	}{
		{
			name:  "Find existing scope",
			param: uuid.MustParse("10000000-1000-1000-2000-000000000001"),
			expected: &models.Scope{
				ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
				Name:        models.SsoServiceType,
				Description: "SSO-service scope",
			},
			error: false,
		},
		{
			name:     "Find non-existing scope",
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
			result, err := scopeRepository.FindById(ctx, tt.param)

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

func Test_ScopeRepository_Delete(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	scopeRepository := NewScopeRepository(client)

	existingScope, err := scopeRepository.Create(ctx, db.CreateScopeParams{
		Name:        "Temp-service",
		Description: "Temporary service scope",
	})
	assert.NoError(t, err)

	tests := []struct {
		name  string
		param uuid.UUID
	}{
		{
			name:  "Delete existing scope",
			param: existingScope.ID,
		},
		{
			name:  "Delete non-existing scope",
			param: uuid.New(),
		},
		{
			name:  "Delete with nil UUID",
			param: uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := scopeRepository.Delete(ctx, tt.param)

			assert.NoError(t, err)
			assert.Equal(t, true, result)
		})
	}
}

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
				ID:          uuid.MustParse("10000000-1000-1000-2000-000000000002"),
				Name:        models.SelfServiceType,
				Description: "Self-service scope",
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
				ScopeID: uuid.MustParse("10000000-1000-1000-2000-000000000002"),
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
				{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000002"),
					Name:        models.SelfServiceType,
					Description: "Self-service scope",
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
