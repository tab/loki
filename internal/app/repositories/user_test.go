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

func Test_UserRepository_List(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	userRepository := NewUserRepository(client)

	_, err = userRepository.Create(ctx, db.CreateUserParams{
		IdentityNumber: "PNOEE-30303039914",
		PersonalCode:   "30303039914",
		FirstName:      "TESTNUMBER",
		LastName:       "OK",
	})
	assert.NoError(t, err)

	_, err = userRepository.Create(ctx, db.CreateUserParams{
		IdentityNumber: "PNOEE-50001029996",
		PersonalCode:   "50001029996",
		FirstName:      "TESTNUMBER",
		LastName:       "ADULT",
	})
	assert.NoError(t, err)

	tests := []struct {
		name     string
		limit    uint64
		offset   uint64
		total    uint64
		expected []models.User
	}{
		{
			name:   "List users",
			limit:  10,
			offset: 0,
			total:  2,
			expected: []models.User{
				{
					IdentityNumber: "PNOEE-30303039914",
					PersonalCode:   "30303039914",
					FirstName:      "TESTNUMBER",
					LastName:       "OK",
				},
				{
					IdentityNumber: "PNOEE-50001029996",
					PersonalCode:   "50001029996",
					FirstName:      "TESTNUMBER",
					LastName:       "ADULT",
				},
			},
		},
		{
			name:   "List with offset",
			limit:  1,
			offset: 1,
			total:  2,
			expected: []models.User{
				{
					IdentityNumber: "PNOEE-50001029996",
					PersonalCode:   "50001029996",
					FirstName:      "TESTNUMBER",
					LastName:       "ADULT",
				},
			},
		},
		{
			name:     "List with zero limit",
			limit:    0,
			offset:   0,
			total:    0,
			expected: []models.User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, total, err := userRepository.List(ctx, tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(results))
			assert.Equal(t, tt.total, total)
		})
	}
}

func Test_UserRepository_Create(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	userRepository := NewUserRepository(client)

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
				_, err = userRepository.Create(ctx, db.CreateUserParams{
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

			result, err := userRepository.Create(ctx, tt.params)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.NotEqual(t, uuid.Nil, result.ID)
				assert.Equal(t, tt.expected.IdentityNumber, result.IdentityNumber)
				assert.Equal(t, tt.expected.PersonalCode, result.PersonalCode)
				assert.Equal(t, tt.expected.FirstName, result.FirstName)
				assert.Equal(t, tt.expected.LastName, result.LastName)
			}
		})
	}
}

func Test_UserRepository_Update(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	repo := NewUserRepository(client)

	account, err := repo.Create(ctx, db.CreateUserParams{
		IdentityNumber: "PNOEE-30303039914",
		PersonalCode:   "30303039914",
		FirstName:      "TESTNUMBER",
		LastName:       "OK",
	})
	assert.NoError(t, err)

	tests := []struct {
		name     string
		before   func()
		params   db.UpdateUserParams
		expected *models.User
		error    bool
	}{
		{
			name:   "Success",
			before: func() {},
			params: db.UpdateUserParams{
				ID:             account.ID,
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "JOHN",
				LastName:       "DOE",
			},
			expected: &models.User{
				ID:             account.ID,
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "JOHN",
				LastName:       "DOE",
			},
			error: false,
		},
		{
			name:   "User not found",
			before: func() {},
			params: db.UpdateUserParams{
				ID:             uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "JOHN",
				LastName:       "DOE",
			},
			expected: nil,
			error:    true,
		},
		{
			name:   "Invalid identity number",
			before: func() {},
			params: db.UpdateUserParams{
				ID:             account.ID,
				IdentityNumber: "",
				PersonalCode:   "30303039914",
				FirstName:      "JOHN",
				LastName:       "DOE",
			},
			expected: nil,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := repo.Update(ctx, tt.params)

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

func Test_UserRepository_FindById(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	repo := NewUserRepository(client)

	account, err := repo.Create(ctx, db.CreateUserParams{
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
			id:   account.ID,
			expected: &models.User{
				ID:             account.ID,
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
			result, err := repo.FindById(ctx, tt.id)

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

func Test_UserRepository_Delete(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	repo := NewUserRepository(client)

	account, err := repo.Create(ctx, db.CreateUserParams{
		IdentityNumber: "PNOEE-30303039914",
		PersonalCode:   "30303039914",
		FirstName:      "TESTNUMBER",
		LastName:       "OK",
	})
	assert.NoError(t, err)

	tests := []struct {
		name  string
		id    uuid.UUID
		error bool
	}{
		{
			name:  "Success",
			id:    account.ID,
			error: false,
		},
		{
			name:  "User not found",
			id:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			error: true,
		},
		{
			name:  "Delete with nil UUID",
			id:    uuid.Nil,
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Delete(ctx, tt.id)

			assert.NoError(t, err)
			assert.Equal(t, true, result)
		})
	}
}

func Test_UserRepository_FindByIdentityNumber(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	repo := NewUserRepository(client)

	account, err := repo.Create(ctx, db.CreateUserParams{
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
			identityNumber: account.IdentityNumber,
			expected: &models.User{
				ID:             account.ID,
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
			result, err := repo.FindByIdentityNumber(ctx, tt.identityNumber)

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

func Test_UserRepository_FindUserDetailsById(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	client, err := postgres.NewPostgresClient(cfg)
	assert.NoError(t, err)

	repo := NewUserRepository(client)

	account, err := repo.Create(ctx, db.CreateUserParams{
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
			id:   account.ID,
			expected: &models.User{
				ID:             account.ID,
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "OK",
				RoleIDs: []uuid.UUID{
					uuid.MustParse("10000000-1000-1000-1000-000000000003"),
				},
				ScopeIDs: []uuid.UUID{
					uuid.MustParse("10000000-1000-1000-2000-000000000002"),
				},
			},
			error: false,
		},
		{
			name:     "User not found",
			id:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			expected: nil,
			error:    true,
		},
		{
			name:     "Delete with nil UUID",
			id:       uuid.Nil,
			expected: nil,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindUserDetailsById(ctx, tt.id)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.IdentityNumber, result.IdentityNumber)
				assert.Equal(t, tt.expected.PersonalCode, result.PersonalCode)
				assert.Equal(t, tt.expected.FirstName, result.FirstName)
				assert.Equal(t, tt.expected.LastName, result.LastName)
				assert.Equal(t, tt.expected.RoleIDs, result.RoleIDs)
				assert.Equal(t, tt.expected.ScopeIDs, result.ScopeIDs)
			}
		})
	}
}
