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
	"loki/internal/config"
	"loki/internal/config/logger"
)

func Test_Roles_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	repository := repositories.NewMockRoleRepository(ctrl)
	service := NewRoles(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected []models.Role
		total    uint64
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().List(ctx, uint64(10), uint64(0)).Return([]models.Role{
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
				}, uint64(3), nil)
			},
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
			total: uint64(3),
			error: nil,
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().List(ctx, uint64(10), uint64(0)).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
			},
			expected: nil,
			total:    0,
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

func Test_Roles_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	repository := repositories.NewMockRoleRepository(ctrl)
	service := NewRoles(repository, log)

	tests := []struct {
		name     string
		before   func()
		params   *models.Role
		expected *models.Role
		error    error
	}{
		{
			name: "Success",
			params: &models.Role{
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				repository.EXPECT().Create(ctx, db.CreateRoleParams{
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}).Return(&models.Role{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}, nil)
			},
			expected: &models.Role{
				ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
		},
		{
			name: "Error",
			params: &models.Role{
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				repository.EXPECT().Create(ctx, db.CreateRoleParams{
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}).Return(nil, errors.ErrFailedToCreateRecord)
			},
			expected: nil,
			error:    errors.ErrFailedToCreateRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Create(ctx, tt.params)

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

func Test_Roles_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	repository := repositories.NewMockRoleRepository(ctrl)
	service := NewRoles(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected *models.Role
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Update(ctx, db.UpdateRoleParams{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}).Return(&models.Role{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}, nil)
			},
			expected: &models.Role{
				ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().Update(ctx, db.UpdateRoleParams{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}).Return(nil, errors.ErrFailedToUpdateRecord)
			},
			expected: nil,
			error:    errors.ErrFailedToUpdateRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Update(ctx, &models.Role{
				ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
				Name:        models.AdminRoleType,
				Description: "Admin role",
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

func Test_Roles_FindById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	repository := repositories.NewMockRoleRepository(ctrl)
	service := NewRoles(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected *models.Role
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().FindById(ctx, uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(&models.Role{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}, nil)
			},
			expected: &models.Role{
				ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().FindById(ctx, uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(nil, errors.ErrRecordNotFound)
			},
			expected: nil,
			error:    errors.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindById(ctx, uuid.MustParse("10000000-1000-1000-1000-000000000001"))

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

func Test_Roles_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	repository := repositories.NewMockRoleRepository(ctrl)
	service := NewRoles(repository, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")

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

func Test_Roles_FindRoleDetailsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	repository := repositories.NewMockRoleRepository(ctrl)
	service := NewRoles(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected *models.Role
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().FindRoleDetailsById(ctx, uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(&models.Role{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
					PermissionIDs: []uuid.UUID{
						uuid.MustParse("10000000-1000-1000-3000-000000000001"),
						uuid.MustParse("10000000-1000-1000-3000-000000000002"),
						uuid.MustParse("10000000-1000-1000-3000-000000000003"),
						uuid.MustParse("10000000-1000-1000-3000-000000000004"),
					},
				}, nil)
			},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindRoleDetailsById(ctx, uuid.MustParse("10000000-1000-1000-1000-000000000001"))

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
