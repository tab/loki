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
	"loki/internal/app/serializers"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

func Test_Tokens_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	database := repositories.NewMockDatabase(ctrl)
	jwtService := jwt.NewMockJwt(ctrl)
	log := logger.NewLogger()
	service := NewTokens(database, jwtService, log)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	user := &models.User{
		ID:             id,
		IdentityNumber: "PNOEE-123456789",
		PersonalCode:   "123456789",
		FirstName:      "John",
		LastName:       "Doe",
	}

	type result struct {
		accessToken  string
		refreshToken string
		error        error
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				database.EXPECT().FindUserRoles(ctx, user.ID).Return([]models.Role{}, nil)
				database.EXPECT().FindUserPermissions(ctx, user.ID).Return([]models.Permission{}, nil)
				database.EXPECT().FindUserScopes(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID:          "PNOEE-123456789",
					Roles:       []string{},
					Permissions: []string{},
					Scope:       []string{},
				}, models.AccessTokenExp).Return("access-token", nil)
				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.RefreshTokenExp).Return("refresh-token", nil)

				database.EXPECT().CreateUserTokens(ctx, gomock.Any()).Return([]models.Token{}, nil)
			},
			expected: result{
				accessToken:  "access-token",
				refreshToken: "refresh-token",
				error:        nil,
			},
		},
		{
			name: "Failed to find user roles",
			before: func() {
				database.EXPECT().FindUserRoles(ctx, user.ID).Return(nil, assert.AnError)
			},
			expected: result{
				accessToken:  "",
				refreshToken: "",
				error:        assert.AnError,
			},
		},
		{
			name: "Failed to find user permissions",
			before: func() {
				database.EXPECT().FindUserRoles(ctx, user.ID).Return([]models.Role{}, nil)
				database.EXPECT().FindUserPermissions(ctx, user.ID).Return(nil, assert.AnError)
			},
			expected: result{
				accessToken:  "",
				refreshToken: "",
				error:        assert.AnError,
			},
		},
		{
			name: "Failed to find user scopes",
			before: func() {
				database.EXPECT().FindUserRoles(ctx, user.ID).Return([]models.Role{}, nil)
				database.EXPECT().FindUserPermissions(ctx, user.ID).Return([]models.Permission{}, nil)
				database.EXPECT().FindUserScopes(ctx, user.ID).Return(nil, assert.AnError)
			},
			expected: result{
				accessToken:  "",
				refreshToken: "",
				error:        assert.AnError,
			},
		},
		{
			name: "Failed to generate access token",
			before: func() {
				database.EXPECT().FindUserRoles(ctx, user.ID).Return([]models.Role{}, nil)
				database.EXPECT().FindUserPermissions(ctx, user.ID).Return([]models.Permission{}, nil)
				database.EXPECT().FindUserScopes(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID:          "PNOEE-123456789",
					Roles:       []string{},
					Permissions: []string{},
					Scope:       []string{},
				}, models.AccessTokenExp).Return("", assert.AnError)
			},
			expected: result{
				accessToken:  "",
				refreshToken: "",
				error:        assert.AnError,
			},
		},
		{
			name: "Failed to save user tokens",
			before: func() {
				database.EXPECT().FindUserRoles(ctx, user.ID).Return([]models.Role{}, nil)
				database.EXPECT().FindUserPermissions(ctx, user.ID).Return([]models.Permission{}, nil)
				database.EXPECT().FindUserScopes(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID:          "PNOEE-123456789",
					Roles:       []string{},
					Permissions: []string{},
					Scope:       []string{},
				}, models.AccessTokenExp).Return("access-token", nil)
				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.RefreshTokenExp).Return("refresh-token", nil)

				database.EXPECT().CreateUserTokens(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			expected: result{
				accessToken:  "",
				refreshToken: "",
				error:        assert.AnError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			accessToken, refreshToken, err := service.Generate(ctx, user)

			if tt.expected.error != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.accessToken, accessToken)
				assert.Equal(t, tt.expected.refreshToken, refreshToken)
			}
		})
	}
}

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

	user := &models.User{
		ID:             id,
		IdentityNumber: "PNOEE-123456789",
		PersonalCode:   "123456789",
		FirstName:      "John",
		LastName:       "Doe",
	}

	tests := []struct {
		name     string
		before   func()
		expected *serializers.TokensSerializer
		error    error
	}{
		{
			name: "Success",
			before: func() {
				database.EXPECT().FindUserByIdentityNumber(ctx, "PNOEE-123456789").Return(user, nil)

				jwtService.EXPECT().Decode(token).Return(&jwt.Payload{
					ID: "PNOEE-123456789",
				}, nil)

				database.EXPECT().FindUserRoles(ctx, user.ID).Return([]models.Role{}, nil)
				database.EXPECT().FindUserPermissions(ctx, user.ID).Return([]models.Permission{}, nil)
				database.EXPECT().FindUserScopes(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID:          "PNOEE-123456789",
					Roles:       []string{},
					Permissions: []string{},
					Scope:       []string{},
				}, models.AccessTokenExp).Return("access-token", nil)
				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.RefreshTokenExp).Return("refresh-token", nil)

				database.EXPECT().CreateUserTokens(ctx, gomock.Any()).Return([]models.Token{}, nil)
			},
			expected: &serializers.TokensSerializer{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
			error: nil,
		},
		{
			name: "Failed to decode token",
			before: func() {
				database.EXPECT().FindUserByIdentityNumber(ctx, "PNOEE-123456789").Return(user, nil)

				jwtService.EXPECT().Decode(token).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
		{
			name: "Invalid token",
			before: func() {
				jwtService.EXPECT().Decode(token).Return(nil, errors.ErrInvalidToken)
			},
			expected: nil,
			error:    errors.ErrInvalidToken,
		},
		{
			name: "Failed to generate access token",
			before: func() {
				database.EXPECT().FindUserByIdentityNumber(ctx, "PNOEE-123456789").Return(user, nil)

				jwtService.EXPECT().Decode(token).Return(&jwt.Payload{
					ID: "PNOEE-123456789",
				}, nil)

				database.EXPECT().FindUserRoles(ctx, user.ID).Return([]models.Role{}, nil)
				database.EXPECT().FindUserPermissions(ctx, user.ID).Return([]models.Permission{}, nil)
				database.EXPECT().FindUserScopes(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID:          "PNOEE-123456789",
					Roles:       []string{},
					Permissions: []string{},
					Scope:       []string{},
				}, models.AccessTokenExp).Return("", assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
		{
			name: "Failed to generate refresh token",
			before: func() {
				database.EXPECT().FindUserByIdentityNumber(ctx, "PNOEE-123456789").Return(user, nil)

				jwtService.EXPECT().Decode(token).Return(&jwt.Payload{
					ID: "PNOEE-123456789",
				}, nil)

				database.EXPECT().FindUserRoles(ctx, user.ID).Return([]models.Role{}, nil)
				database.EXPECT().FindUserPermissions(ctx, user.ID).Return([]models.Permission{}, nil)
				database.EXPECT().FindUserScopes(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID:          "PNOEE-123456789",
					Roles:       []string{},
					Permissions: []string{},
					Scope:       []string{},
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
				jwtService.EXPECT().Decode(token).Return(&jwt.Payload{
					ID: "PNOEE-123456789",
				}, nil)

				database.EXPECT().FindUserRoles(ctx, user.ID).Return([]models.Role{}, nil)
				database.EXPECT().FindUserPermissions(ctx, user.ID).Return([]models.Permission{}, nil)
				database.EXPECT().FindUserScopes(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(jwt.Payload{
					ID:          "PNOEE-123456789",
					Roles:       []string{},
					Permissions: []string{},
					Scope:       []string{},
				}, models.AccessTokenExp).Return("access-token", nil)
				jwtService.EXPECT().Generate(jwt.Payload{
					ID: "PNOEE-123456789",
				}, models.RefreshTokenExp).Return("refresh-token", nil)

				database.EXPECT().CreateUserTokens(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Refresh(ctx, token)

			if tt.error != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
