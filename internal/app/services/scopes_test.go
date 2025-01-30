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

func Test_Scopes_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockScopeRepository(ctrl)
	log := logger.NewLogger()
	service := NewScopes(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected []models.Scope
		total    int
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().List(ctx, int32(10), int32(0)).Return([]models.Scope{
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
				}, 2, nil)
			},
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

func Test_Scopes_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockScopeRepository(ctrl)
	log := logger.NewLogger()
	service := NewScopes(repository, log)

	tests := []struct {
		name     string
		before   func()
		params   *models.Scope
		expected *models.Scope
		error    error
	}{
		{
			name: "Success",
			params: &models.Scope{
				Name:        models.SsoServiceType,
				Description: "SSO-service scope",
			},
			before: func() {
				repository.EXPECT().Create(ctx, db.CreateScopeParams{
					Name:        models.SsoServiceType,
					Description: "SSO-service scope",
				}).Return(&models.Scope{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        models.SsoServiceType,
					Description: "SSO-service scope",
				}, nil)
			},
			expected: &models.Scope{
				ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
				Name:        models.SsoServiceType,
				Description: "SSO-service scope",
			},
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

func Test_Scopes_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockScopeRepository(ctrl)
	log := logger.NewLogger()
	service := NewScopes(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected *models.Scope
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Update(ctx, db.UpdateScopeParams{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        models.SelfServiceType,
					Description: "Self-service scope",
				}).Return(&models.Scope{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        models.SelfServiceType,
					Description: "Self-service scope",
				}, nil)
			},
			expected: &models.Scope{
				ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
				Name:        models.SelfServiceType,
				Description: "Self-service scope",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Update(ctx, &models.Scope{
				ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
				Name:        models.SelfServiceType,
				Description: "Self-service scope",
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

func Test_Scopes_FindById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockScopeRepository(ctrl)
	log := logger.NewLogger()
	service := NewScopes(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected *models.Scope
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().FindById(ctx, uuid.MustParse("10000000-1000-1000-2000-000000000001")).Return(&models.Scope{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        models.SsoServiceType,
					Description: "SSO-service scope",
				}, nil)
			},
			expected: &models.Scope{
				ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
				Name:        models.SsoServiceType,
				Description: "SSO-service scope",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindById(ctx, uuid.MustParse("10000000-1000-1000-2000-000000000001"))

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

func Test_Scopes_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := repositories.NewMockScopeRepository(ctrl)
	log := logger.NewLogger()
	service := NewScopes(repository, log)

	tests := []struct {
		name     string
		before   func()
		expected bool
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Delete(ctx, uuid.MustParse("10000000-1000-1000-2000-000000000001")).Return(true, nil)
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Delete(ctx, uuid.MustParse("10000000-1000-1000-2000-000000000001"))

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
