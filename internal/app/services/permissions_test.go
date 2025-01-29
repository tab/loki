package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/repositories/db"
	"loki/pkg/logger"
)

func Test_Permissions_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockPermissionRepository(ctrl)
	log := logger.NewLogger()
	service := NewPermissions(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected []models.Permission
		total    int
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().List(ctx, int32(10), int32(0)).Return([]models.Permission{
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
				}, 2, nil)
			},
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
			total: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, total, err := service.List(ctx, &Pagination{
				Page:    int32(1),
				PerPage: int32(10),
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

func Test_Permissions_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockPermissionRepository(ctrl)
	log := logger.NewLogger()
	service := NewPermissions(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected *models.Permission
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Create(ctx, db.CreatePermissionParams{
					Name:        "read:self",
					Description: "Read own data",
				}).Return(&models.Permission{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				}, nil)
			},
			expected: &models.Permission{
				ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
				Name:        "read:self",
				Description: "Read own data",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().Create(ctx, db.CreatePermissionParams{
					Name:        "read:self",
					Description: "Read own data",
				}).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Create(ctx, &models.Permission{
				Name:        "read:self",
				Description: "Read own data",
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

func Test_Permissions_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockPermissionRepository(ctrl)
	log := logger.NewLogger()
	service := NewPermissions(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected *models.Permission
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Update(ctx, db.UpdatePermissionParams{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				}).Return(&models.Permission{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				}, nil)
			},
			expected: &models.Permission{
				ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
				Name:        "read:self",
				Description: "Read own data",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().Update(ctx, db.UpdatePermissionParams{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				}).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Update(ctx, &models.Permission{
				ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
				Name:        "read:self",
				Description: "Read own data",
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

func Test_Permissions_FindById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockPermissionRepository(ctrl)
	log := logger.NewLogger()
	service := NewPermissions(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected *models.Permission
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().FindById(ctx, uuid.MustParse("10000000-1000-1000-3000-000000000001")).Return(&models.Permission{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				}, nil)
			},
			expected: &models.Permission{
				ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
				Name:        "read:self",
				Description: "Read own data",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().FindById(ctx, uuid.MustParse("10000000-1000-1000-3000-000000000001")).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindById(ctx, uuid.MustParse("10000000-1000-1000-3000-000000000001"))

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

func Test_Permissions_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockPermissionRepository(ctrl)
	log := logger.NewLogger()
	service := NewPermissions(repository, log)

	tests := []struct {
		name  string
		id    uuid.UUID
		error error
	}{
		{
			name: "Success",
			id:   uuid.MustParse("10000000-1000-1000-3000-000000000001"),
		},
		{
			name:  "Error",
			id:    uuid.MustParse("10000000-1000-1000-3000-000000000001"),
			error: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.error != nil {
				repository.EXPECT().Delete(ctx, tt.id).Return(false, assert.AnError)
			} else {
				repository.EXPECT().Delete(ctx, tt.id).Return(true, nil)
			}

			result, err := service.Delete(ctx, tt.id)

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
