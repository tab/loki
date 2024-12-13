package repositories

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
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

	code := m.Run()
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
			name: "Success",
			before: func() {
			},
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
				_, err := repo.CreateUser(ctx, db.CreateUserParams{
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

			t.Cleanup(func() {
				err = spec.TruncateTables(ctx, cfg.DatabaseDSN)
				require.NoError(t, err)
			})
		})
	}
}

func Test_Database_CreateOrUpdateUserWithTokens(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	repo, err := NewDatabase(cfg)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		before   func()
		params   dto.CreateUserParams
		expected *models.User
		error    bool
	}{
		{
			name: "Success",
			before: func() {
			},
			params: dto.CreateUserParams{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "OK",
				AccessToken: dto.CreateTokenParams{
					Type:      "access_token",
					Value:     "aaa.bbb.ccc",
					ExpiresAt: time.Now().Add(time.Hour),
				},
				RefreshToken: dto.CreateTokenParams{
					Type:      "refresh_token",
					Value:     "ddd.eee.fff",
					ExpiresAt: time.Now().Add(time.Hour * 24),
				},
			},
			expected: &models.User{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "OK",
				AccessToken:    "aaa.bbb.ccc",
				RefreshToken:   "ddd.eee.fff",
			},
			error: false,
		},
		{
			name: "User already exists",
			before: func() {
				_, err := repo.CreateUser(ctx, db.CreateUserParams{
					IdentityNumber: "PNOEE-30303039914",
					PersonalCode:   "30303039914",
					FirstName:      "TESTNUMBER",
					LastName:       "OK",
				})
				assert.NoError(t, err)
			},
			params: dto.CreateUserParams{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "JOHN",
				LastName:       "DOE",
				AccessToken: dto.CreateTokenParams{
					Type:      "access_token",
					Value:     "aaa.bbb.ccc",
					ExpiresAt: time.Now().Add(time.Hour),
				},
				RefreshToken: dto.CreateTokenParams{
					Type:      "refresh_token",
					Value:     "ddd.eee.fff",
					ExpiresAt: time.Now().Add(time.Hour * 24),
				},
			},
			expected: &models.User{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "JOHN",
				LastName:       "DOE",
				AccessToken:    "aaa.bbb.ccc",
				RefreshToken:   "ddd.eee.fff",
			},
			error: false,
		},
		{
			name:   "Invalid identity number",
			before: func() {},
			params: dto.CreateUserParams{
				IdentityNumber: "",
				PersonalCode:   "",
				FirstName:      "TESTNUMBER",
				LastName:       "NOT OK",
			},
			expected: nil,
			error:    true,
		},
		{
			name:   "Invalid access token",
			before: func() {},
			params: dto.CreateUserParams{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "NOT OK",
				AccessToken: dto.CreateTokenParams{
					Type:      "",
					Value:     "",
					ExpiresAt: time.Time{},
				},
				RefreshToken: dto.CreateTokenParams{
					Type:      "refresh_token",
					Value:     "ddd.eee.fff",
					ExpiresAt: time.Now().Add(time.Hour * 24),
				},
			},
			expected: nil,
			error:    true,
		},
		{
			name:   "Invalid refresh token",
			before: func() {},
			params: dto.CreateUserParams{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "NOT OK",
				AccessToken: dto.CreateTokenParams{
					Type:      "access_token",
					Value:     "aaa.bbb.ccc",
					ExpiresAt: time.Now().Add(time.Hour),
				},
				RefreshToken: dto.CreateTokenParams{
					Type:      "",
					Value:     "",
					ExpiresAt: time.Time{},
				},
			},
			expected: nil,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			user, err := repo.CreateOrUpdateUserWithTokens(ctx, tt.params)

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

			t.Cleanup(func() {
				err = spec.TruncateTables(ctx, cfg.DatabaseDSN)
				require.NoError(t, err)
			})
		})
	}
}

func Test_Database_RefreshUserTokens(t *testing.T) {
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
		params   dto.RefreshTokenParams
		expected *models.User
		error    bool
	}{
		{
			name: "Success",
			params: dto.RefreshTokenParams{
				UserId: user.ID,
				AccessToken: dto.CreateTokenParams{
					Type:      "access_token",
					Value:     "aaa.bbb.ccc",
					ExpiresAt: time.Now().Add(time.Hour),
				},
				RefreshToken: dto.CreateTokenParams{
					Type:      "refresh_token",
					Value:     "ddd.eee.fff",
					ExpiresAt: time.Now().Add(time.Hour * 24),
				},
			},
			expected: &models.User{
				IdentityNumber: "PNOEE-30303039914",
				PersonalCode:   "30303039914",
				FirstName:      "TESTNUMBER",
				LastName:       "OK",
				AccessToken:    "aaa.bbb.ccc",
				RefreshToken:   "ddd.eee.fff",
			},
			error: false,
		},
		{
			name: "User not found",
			params: dto.RefreshTokenParams{
				UserId: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				AccessToken: dto.CreateTokenParams{
					Type:      "access_token",
					Value:     "aaa.bbb.ccc",
					ExpiresAt: time.Now().Add(time.Hour),
				},
				RefreshToken: dto.CreateTokenParams{
					Type:      "refresh_token",
					Value:     "ddd.eee.fff",
					ExpiresAt: time.Now().Add(time.Hour * 24),
				},
			},
			expected: nil,
			error:    true,
		},
		{
			name: "Invalid access token",
			params: dto.RefreshTokenParams{
				UserId: user.ID,
				AccessToken: dto.CreateTokenParams{
					Type:      "",
					Value:     "",
					ExpiresAt: time.Time{},
				},
				RefreshToken: dto.CreateTokenParams{
					Type:      "refresh_token",
					Value:     "ddd.eee.fff",
					ExpiresAt: time.Now().Add(time.Hour * 24),
				},
			},
			expected: nil,
			error:    true,
		},
		{
			name: "Invalid refresh token",
			params: dto.RefreshTokenParams{
				UserId: user.ID,
				AccessToken: dto.CreateTokenParams{
					Type:      "access_token",
					Value:     "aaa.bbb.ccc",
					ExpiresAt: time.Now().Add(time.Hour),
				},
				RefreshToken: dto.CreateTokenParams{
					Type:      "",
					Value:     "",
					ExpiresAt: time.Time{},
				},
			},
			expected: nil,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.RefreshUserTokens(ctx, tt.params)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.NotEqual(t, uuid.Nil, result.ID)
				assert.Equal(t, tt.expected.IdentityNumber, result.IdentityNumber)
				assert.Equal(t, tt.expected.PersonalCode, result.PersonalCode)
				assert.Equal(t, tt.expected.FirstName, result.FirstName)
				assert.Equal(t, tt.expected.LastName, result.LastName)
				assert.Equal(t, tt.expected.AccessToken, result.AccessToken)
				assert.Equal(t, tt.expected.RefreshToken, result.RefreshToken)
			}

			t.Cleanup(func() {
				err = spec.TruncateTables(ctx, cfg.DatabaseDSN)
				require.NoError(t, err)
			})
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

			t.Cleanup(func() {
				err = spec.TruncateTables(ctx, cfg.DatabaseDSN)
				require.NoError(t, err)
			})
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

			t.Cleanup(func() {
				err = spec.TruncateTables(ctx, cfg.DatabaseDSN)
				require.NoError(t, err)
			})
		})
	}
}
