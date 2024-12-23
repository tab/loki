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
				}).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
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
				repository.EXPECT().FindByIdentityNumber(ctx, identityNumber).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
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
