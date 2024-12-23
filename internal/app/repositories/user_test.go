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
