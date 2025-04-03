package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/repositories/db"
	"loki/pkg/logger"
)

func Test_Users_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockUserRepository(ctrl)
	log := logger.NewLogger()
	service := NewUsers(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected []models.User
		total    uint64
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().List(ctx, uint64(10), uint64(0)).Return([]models.User{
					{
						ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
						IdentityNumber: "PNOEE-123456789",
						PersonalCode:   "123456789",
						FirstName:      "John",
						LastName:       "Doe",
					},
					{
						ID:             uuid.MustParse("10000000-1000-1000-1000-000000000002"),
						IdentityNumber: "PNOEE-987654321",
						PersonalCode:   "987654321",
						FirstName:      "Jane",
						LastName:       "Doe",
					},
				}, uint64(2), nil)
			},
			expected: []models.User{
				{
					ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				},
				{
					ID:             uuid.MustParse("10000000-1000-1000-1000-000000000002"),
					IdentityNumber: "PNOEE-987654321",
					PersonalCode:   "987654321",
					FirstName:      "Jane",
					LastName:       "Doe",
				},
			},
			total: uint64(2),
			error: nil,
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().List(ctx, uint64(10), uint64(0)).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
			},
			expected: nil,
			total:    uint64(0),
			error:    errors.ErrFailedToFetchResults,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, total, err := service.List(ctx, &Pagination{
				Page:    uint64(1),
				PerPage: uint64(10),
			})

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Zero(t, total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
				assert.Equal(t, tt.total, total)
			}
		})
	}
}

func Test_Users_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockUserRepository(ctrl)
	log := logger.NewLogger()
	service := NewUsers(repository, log)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		before   func()
		expected *models.User
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Create(ctx, db.CreateUserParams{
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)
			},
			expected: &models.User{
				ID:             id,
				IdentityNumber: "PNOEE-123456789",
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
				AccessToken:    "access-token",
				RefreshToken:   "refresh-token",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().Create(ctx, db.CreateUserParams{
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}).Return(nil, errors.ErrFailedToCreateRecord)
			},
			expected: nil,
			error:    errors.ErrFailedToCreateRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Create(ctx, &models.User{
				IdentityNumber: "PNOEE-123456789",
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			})

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Users_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockUserRepository(ctrl)
	log := logger.NewLogger()
	service := NewUsers(repository, log)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		before   func()
		expected *models.User
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Update(ctx, db.UpdateUserParams{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)
			},
			expected: &models.User{
				ID:             id,
				IdentityNumber: "PNOEE-123456789",
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
				AccessToken:    "access-token",
				RefreshToken:   "refresh-token",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().Update(ctx, db.UpdateUserParams{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}).Return(nil, errors.ErrFailedToUpdateRecord)
			},
			expected: nil,
			error:    errors.ErrFailedToUpdateRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Update(ctx, &models.User{
				ID:             id,
				IdentityNumber: "PNOEE-123456789",
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			})

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Users_FindById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockUserRepository(ctrl)
	log := logger.NewLogger()
	service := NewUsers(repository, log)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		before   func()
		expected *models.User
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().FindById(ctx, id).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}, nil)
			},
			expected: &models.User{
				ID:             id,
				IdentityNumber: "PNOEE-123456789",
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().FindById(ctx, id).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindById(ctx, id)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Users_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockUserRepository(ctrl)
	log := logger.NewLogger()
	service := NewUsers(repository, log)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		before   func()
		expected bool
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Delete(ctx, id).Return(true, nil)
			},
			expected: true,
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().Delete(ctx, id).Return(false, errors.ErrFailedToDeleteRecord)
			},
			expected: false,
			error:    errors.ErrFailedToDeleteRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Delete(ctx, id)

			if tt.error != nil {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.True(t, result)
			}
		})
	}
}

func Test_Users_FindByIdentityNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockUserRepository(ctrl)
	log := logger.NewLogger()
	service := NewUsers(repository, log)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	identityNumber := "PNOEE-123456789"

	tests := []struct {
		name     string
		before   func()
		expected *models.User
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().FindByIdentityNumber(ctx, identityNumber).Return(&models.User{
					ID:             id,
					IdentityNumber: identityNumber,
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}, nil)
			},
			expected: &models.User{
				ID:             id,
				IdentityNumber: identityNumber,
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().FindByIdentityNumber(ctx, identityNumber).Return(nil, errors.ErrRecordNotFound)
			},
			expected: nil,
			error:    errors.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindByIdentityNumber(ctx, identityNumber)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Users_FindUserDetailsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockUserRepository(ctrl)
	log := logger.NewLogger()
	service := NewUsers(repository, log)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		before   func()
		expected *models.User
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().FindUserDetailsById(ctx, id).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
					RoleIDs: []uuid.UUID{
						uuid.MustParse("10000000-1000-1000-1000-000000000003"),
					},
					ScopeIDs: []uuid.UUID{
						uuid.MustParse("10000000-1000-1000-2000-000000000002"),
					},
				}, nil)
			},
			expected: &models.User{
				ID:             id,
				IdentityNumber: "PNOEE-123456789",
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
				RoleIDs: []uuid.UUID{
					uuid.MustParse("10000000-1000-1000-1000-000000000003"),
				},
				ScopeIDs: []uuid.UUID{
					uuid.MustParse("10000000-1000-1000-2000-000000000002"),
				},
			},
		},
		{
			name: "Without Roles and Scopes",
			before: func() {
				repository.EXPECT().FindUserDetailsById(ctx, id).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}, nil)
			},
			expected: &models.User{
				ID:             id,
				IdentityNumber: "PNOEE-123456789",
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			},
			error: nil,
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().FindUserDetailsById(ctx, id).Return(nil, errors.ErrRecordNotFound)
			},
			expected: nil,
			error:    errors.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindUserDetailsById(ctx, id)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
