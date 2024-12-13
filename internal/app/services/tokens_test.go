package services

import (
	"context"
	"loki/internal/app/errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/serializers"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

func Test_Tokens_Refresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	database := repositories.NewMockDatabase(ctrl)
	jwtService := jwt.NewMockJwt(ctrl)
	log := logger.NewLogger()
	service := NewTokens(database, jwtService, log)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	token := "refresh-token"

	tests := []struct {
		name     string
		before   func()
		expected *serializers.UserSerializer
		error    error
	}{
		{
			name: "Success",
			before: func() {
				jwtService.EXPECT().Verify(token).Return(true, nil)

				database.EXPECT().FindUserById(ctx, id).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.AccessTokenExp).Return("access-token", nil)
				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.RefreshTokenExp).Return("refresh-token", nil)

				database.EXPECT().RefreshUserTokens(ctx, gomock.Any()).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)
			},
			expected: &serializers.UserSerializer{
				ID:             id,
				IdentityNumber: "PNOEE-123456789",
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
				AccessToken:    "access-token",
				RefreshToken:   "refresh-token",
			},
			error: nil,
		},
		{
			name: "Failed to decode token",
			before: func() {
				jwtService.EXPECT().Verify(token).Return(false, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
		{
			name: "Invalid token",
			before: func() {
				jwtService.EXPECT().Verify(token).Return(false, nil)
			},
			expected: nil,
			error:    errors.ErrInvalidToken,
		},
		{
			name: "Failed to find user",
			before: func() {
				jwtService.EXPECT().Verify(token).Return(true, nil)

				database.EXPECT().FindUserById(ctx, id).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
		{
			name: "Failed to generate access token",
			before: func() {
				jwtService.EXPECT().Verify(token).Return(true, nil)

				database.EXPECT().FindUserById(ctx, id).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.AccessTokenExp).Return("", assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
		{
			name: "Failed to generate refresh token",
			before: func() {
				jwtService.EXPECT().Verify(token).Return(true, nil)

				database.EXPECT().FindUserById(ctx, id).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.AccessTokenExp).Return("access-token", nil)
				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.RefreshTokenExp).Return("", assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
		{
			name: "Failed to create user tokens",
			before: func() {
				jwtService.EXPECT().Verify(token).Return(true, nil)

				database.EXPECT().FindUserById(ctx, id).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.AccessTokenExp).Return("access-token", nil)
				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.RefreshTokenExp).Return("refresh-token", nil)

				database.EXPECT().RefreshUserTokens(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			user, err := service.Refresh(ctx, id, token)

			if tt.error != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expected, user)
			}
		})
	}
}