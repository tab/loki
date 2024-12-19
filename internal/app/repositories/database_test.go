package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/config"
	"loki/pkg/spec"
)

func TestMain(m *testing.M) {
	if err := spec.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	if os.Getenv("GO_ENV") == "ci" {
		os.Exit(0)
	}

	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	err := spec.DbSeed(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Error seeding database: %v", err)
	}

	code := m.Run()

	err = spec.TruncateTables(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Error truncating tables: %v", err)
	}

	os.Exit(code)
}

func Test_Database_CreateUser(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		before   func()
		params   db.CreateUserParams
		expected *models.User
		error    bool
	}{
		{
			name:   "Success",
			before: func() {},
			params: db.CreateUserParams{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "OK",
			},
			expected: &models.User{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "OK",
			},
			error: false,
		},
		{
			name: "User already exists",
			before: func() {
				_, err = repo.CreateUser(ctx, db.CreateUserParams{
					IdentityNumber: "PNOEE-30303039914",
					PersonalCode:   "30303039914",
					FirstName:      "TESTNUMBER",
					LastName:       "OK",
				})
				assert.NoError(t, err)
			},
			params: db.CreateUserParams{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "JOHN",
				LastName:       "DOE",
			},
			expected: &models.User{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "JOHN",
				LastName:       "DOE",
			},
			error: false,
		},
		{
			name:   "Invalid identity number",
			before: func() {},
			params: db.CreateUserParams{
				IdentityNumber: "",
				PersonalCode:   "",
				FirstName:      "TESTNUMBER",
				LastName:       "NOT OK",
			},
			expected: nil,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			user, err := repo.CreateUser(ctx, tt.params)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.NotEqual(t, uuid.Nil, user.ID)
				assert.Equal(t, tt.expected.IdentityNumber, user.IdentityNumber)
				assert.Equal(t, tt.expected.PersonalCode, user.PersonalCode)
				assert.Equal(t, tt.expected.FirstName, user.FirstName)
				assert.Equal(t, tt.expected.LastName, user.LastName)
			}
		})
	}
}

func Test_Database_CreateUserTokens(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
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
				UserID:            user.ID,
				AccessTokenValue:  "aaa.bbb.ccc",
				RefreshTokenValue: "ddd.eee.fff",
			},
			expected: []models.Token{
				{
					UserId: user.ID,
				},
				{
					UserId: user.ID,
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
			results, err := repo.CreateUserTokens(ctx, tt.params)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(results))
			}
		})
	}
}

func Test_Database_CreateUserRole(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
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
				UserID: user.ID,
				RoleID: uuid.MustParse("10000000-1000-1000-1000-000000000003"),
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateUserRole(ctx, tt.params)
			assert.NoError(t, err)
		})
	}
}

func Test_Database_CreateUserScope(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
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
				UserID:  user.ID,
				ScopeID: uuid.MustParse("10000000-1000-1000-2000-000000000001"),
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateUserScope(ctx, tt.params)
			assert.NoError(t, err)
		})
	}
}

func Test_Database_FindUserById(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		IdentityNumber: "PNOEE-30303039914",
		PersonalCode:   "30303039914",
		FirstName:      "TESTNUMBER",
		LastName:       "OK",
	})
	assert.NoError(t, err)

	tests := []struct {
		name     string
		id       uuid.UUID
		expected *models.User
		error    bool
	}{
		{
			name: "Success",
			id:   user.ID,
			expected: &models.User{
				ID:             user.ID,
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "OK",
			},
			error: false,
		},
		{
			name:     "User not found",
			id:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			expected: nil,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindUserById(ctx, tt.id)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.IdentityNumber, result.IdentityNumber)
				assert.Equal(t, tt.expected.PersonalCode, result.PersonalCode)
				assert.Equal(t, tt.expected.FirstName, result.FirstName)
				assert.Equal(t, tt.expected.LastName, result.LastName)
			}
		})
	}
}

func Test_Database_FindUserByIdentityNumber(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		IdentityNumber: "PNOEE-30303039914",
		PersonalCode:   "30303039914",
		FirstName:      "TESTNUMBER",
		LastName:       "OK",
	})
	assert.NoError(t, err)

	tests := []struct {
		name           string
		identityNumber string
		expected       *models.User
		error          bool
	}{
		{
			name:           "Success",
			identityNumber: user.IdentityNumber,
			expected: &models.User{
				ID:             user.ID,
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "OK",
			},
			error: false,
		},
		{
			name:           "User not found",
			identityNumber: "PNOEE-123456789",
			expected:       nil,
			error:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindUserByIdentityNumber(ctx, tt.identityNumber)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.IdentityNumber, result.IdentityNumber)
				assert.Equal(t, tt.expected.PersonalCode, result.PersonalCode)
				assert.Equal(t, tt.expected.FirstName, result.FirstName)
				assert.Equal(t, tt.expected.LastName, result.LastName)
			}
		})
	}
}

func Test_Database_FindRoleByName(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

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
			result, err := repo.FindRoleByName(ctx, tt.param)

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

func Test_Database_FindScopeByName(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

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
			result, err := repo.FindScopeByName(ctx, tt.param)

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

func Test_Database_FindUserRoles(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
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
			param: user.ID,
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
			results, err := repo.FindUserRoles(ctx, tt.param)
			assert.NoError(t, err)

			assert.Equal(t, len(tt.expected), len(results))
			for i, role := range results {
				assert.Equal(t, tt.expected[i].ID, role.ID)
				assert.Equal(t, tt.expected[i].Name, role.Name)
			}
		})
	}
}

func Test_Database_FindUserPermissions(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
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
			param: user.ID,
			expected: []models.Permission{
				{ID: uuid.MustParse("10000000-1000-1000-3000-000000000001"), Name: "read:self"},
				{ID: uuid.MustParse("10000000-1000-1000-3000-000000000002"), Name: "write:self"},
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
			results, err := repo.FindUserPermissions(ctx, tt.param)
			assert.NoError(t, err)

			assert.Equal(t, len(tt.expected), len(results))
			for i, permission := range results {
				assert.Equal(t, tt.expected[i].ID, permission.ID)
				assert.Equal(t, tt.expected[i].Name, permission.Name)
			}
		})
	}
}

func Test_Database_FindUserScopes(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
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
			param: user.ID,
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
			results, err := repo.FindUserScopes(ctx, tt.param)
			assert.NoError(t, err)

			assert.Equal(t, len(tt.expected), len(results))
			for i, scope := range results {
				assert.Equal(t, tt.expected[i].ID, scope.ID)
				assert.Equal(t, tt.expected[i].Name, scope.Name)
			}
		})
	}
}
