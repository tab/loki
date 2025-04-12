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
	"loki/internal/config"
	"loki/internal/config/logger"
	"loki/pkg/jwt"
)

func Test_Tokens_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	permissionRepository := repositories.NewMockPermissionRepository(ctrl)
	roleRepository := repositories.NewMockRoleRepository(ctrl)
	scopeRepository := repositories.NewMockScopeRepository(ctrl)
	tokenRepository := repositories.NewMockTokenRepository(ctrl)
	userRepository := repositories.NewMockUserRepository(ctrl)

	jwtService := jwt.NewMockJwt(ctrl)
	service := NewTokens(
		jwtService,
		permissionRepository,
		roleRepository,
		scopeRepository,
		tokenRepository,
		userRepository,
		log,
	)

	tests := []struct {
		name     string
		before   func()
		expected []models.Token
		total    uint64
		error    error
	}{
		{
			name: "Success",
			before: func() {
				tokenRepository.EXPECT().List(ctx, uint64(10), uint64(0)).Return([]models.Token{}, uint64(0), nil)
			},
			expected: []models.Token{},
			total:    uint64(0),
			error:    nil,
		},
		{
			name: "Failed to fetch results",
			before: func() {
				tokenRepository.EXPECT().List(ctx, uint64(10), uint64(0)).Return(nil, uint64(0), assert.AnError)
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

func Test_Tokens_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	permissionRepository := repositories.NewMockPermissionRepository(ctrl)
	roleRepository := repositories.NewMockRoleRepository(ctrl)
	scopeRepository := repositories.NewMockScopeRepository(ctrl)
	tokenRepository := repositories.NewMockTokenRepository(ctrl)
	userRepository := repositories.NewMockUserRepository(ctrl)

	jwtService := jwt.NewMockJwt(ctrl)
	service := NewTokens(
		jwtService,
		permissionRepository,
		roleRepository,
		scopeRepository,
		tokenRepository,
		userRepository,
		log,
	)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

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
		expected *models.User
		err      error
	}{
		{
			name: "Success",
			before: func() {
				userRepository.EXPECT().FindById(ctx, user.ID).Return(user, nil)

				roleRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Role{}, nil)
				permissionRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Permission{}, nil)
				scopeRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID:          "PNOEE-123456789",
						Roles:       []string{},
						Permissions: []string{},
						Scope:       []string{},
					},
					models.AccessTokenExp,
				).Return("access-token", nil)

				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID: "PNOEE-123456789",
					},
					models.RefreshTokenExp,
				).Return("refresh-token", nil)

				tokenRepository.EXPECT().Create(ctx, gomock.Any()).Return([]models.Token{}, nil)
			},
			expected: &models.User{
				ID:             user.ID,
				IdentityNumber: user.IdentityNumber,
				PersonalCode:   user.PersonalCode,
				FirstName:      user.FirstName,
				LastName:       user.LastName,
				AccessToken:    "access-token",
				RefreshToken:   "refresh-token",
			},
			err: nil,
		},
		{
			name: "Failed to find user",
			before: func() {
				userRepository.EXPECT().FindById(ctx, user.ID).Return(user, errors.ErrRecordNotFound)
			},
			expected: nil,
			err:      errors.ErrRecordNotFound,
		},
		{
			name: "Failed to find user roles",
			before: func() {
				userRepository.EXPECT().FindById(ctx, user.ID).Return(user, nil)
				roleRepository.EXPECT().FindByUserId(ctx, user.ID).Return(nil, assert.AnError)
			},
			expected: nil,
			err:      assert.AnError,
		},
		{
			name: "Failed to find user permissions",
			before: func() {
				userRepository.EXPECT().FindById(ctx, user.ID).Return(user, nil)

				roleRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Role{}, nil)
				permissionRepository.EXPECT().FindByUserId(ctx, user.ID).Return(nil, assert.AnError)
			},
			expected: nil,
			err:      assert.AnError,
		},
		{
			name: "Failed to find user scopes",
			before: func() {
				userRepository.EXPECT().FindById(ctx, user.ID).Return(user, nil)

				roleRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Role{}, nil)
				permissionRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Permission{}, nil)
				scopeRepository.EXPECT().FindByUserId(ctx, user.ID).Return(nil, assert.AnError)
			},
			expected: nil,
			err:      assert.AnError,
		},
		{
			name: "Failed to generate access token",
			before: func() {
				userRepository.EXPECT().FindById(ctx, user.ID).Return(user, nil)

				roleRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Role{}, nil)
				permissionRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Permission{}, nil)
				scopeRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID:          "PNOEE-123456789",
						Roles:       []string{},
						Permissions: []string{},
						Scope:       []string{},
					},
					models.AccessTokenExp,
				).Return("", assert.AnError)
			},
			expected: nil,
			err:      assert.AnError,
		},
		{
			name: "Failed to save user tokens",
			before: func() {
				userRepository.EXPECT().FindById(ctx, user.ID).Return(user, nil)

				roleRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Role{}, nil)
				permissionRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Permission{}, nil)
				scopeRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID:          "PNOEE-123456789",
						Roles:       []string{},
						Permissions: []string{},
						Scope:       []string{},
					},
					models.AccessTokenExp,
				).Return("access-token", nil)

				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID: "PNOEE-123456789",
					},
					models.RefreshTokenExp,
				).Return("refresh-token", nil)

				tokenRepository.EXPECT().Create(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			expected: nil,
			err:      assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Create(ctx, user.ID)

			if tt.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.AccessToken, result.AccessToken)
				assert.Equal(t, tt.expected.RefreshToken, result.RefreshToken)
			}
		})
	}
}

func Test_Tokens_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	permissionRepository := repositories.NewMockPermissionRepository(ctrl)
	roleRepository := repositories.NewMockRoleRepository(ctrl)
	scopeRepository := repositories.NewMockScopeRepository(ctrl)
	tokenRepository := repositories.NewMockTokenRepository(ctrl)
	userRepository := repositories.NewMockUserRepository(ctrl)

	jwtService := jwt.NewMockJwt(ctrl)
	service := NewTokens(
		jwtService,
		permissionRepository,
		roleRepository,
		scopeRepository,
		tokenRepository,
		userRepository,
		log,
	)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	refreshTokenValue := "refresh-token"

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
		expected *models.User
		err      error
	}{
		{
			name: "Success",
			before: func() {
				jwtService.EXPECT().Decode(refreshTokenValue).Return(&jwt.Payload{
					ID: "PNOEE-123456789",
				}, nil)
				userRepository.EXPECT().FindByIdentityNumber(ctx, "PNOEE-123456789").Return(user, nil)

				roleRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Role{}, nil)
				permissionRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Permission{}, nil)
				scopeRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID:          "PNOEE-123456789",
						Roles:       []string{},
						Permissions: []string{},
						Scope:       []string{},
					},
					models.AccessTokenExp,
				).Return("new-access-token", nil)
				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID: "PNOEE-123456789",
					},
					models.RefreshTokenExp,
				).Return("new-refresh-token", nil)

				tokenRepository.EXPECT().Create(ctx, gomock.Any()).Return([]models.Token{}, nil)
			},
			expected: &models.User{
				ID:             user.ID,
				IdentityNumber: user.IdentityNumber,
				PersonalCode:   user.PersonalCode,
				FirstName:      user.FirstName,
				LastName:       user.LastName,
				AccessToken:    "new-access-token",
				RefreshToken:   "new-refresh-token",
			},
			err: nil,
		},
		{
			name: "Failed to decode token",
			before: func() {
				jwtService.EXPECT().Decode(refreshTokenValue).Return(nil, assert.AnError)
				userRepository.EXPECT().FindByIdentityNumber(ctx, "PNOEE-123456789").Return(user, nil)
			},
			expected: nil,
			err:      assert.AnError,
		},
		{
			name: "Invalid token",
			before: func() {
				jwtService.EXPECT().Decode(refreshTokenValue).Return(nil, errors.ErrInvalidToken)
			},
			expected: nil,
			err:      errors.ErrInvalidToken,
		},
		{
			name: "Failed to generate access token",
			before: func() {
				jwtService.EXPECT().Decode(refreshTokenValue).Return(&jwt.Payload{
					ID: "PNOEE-123456789",
				}, nil)
				userRepository.EXPECT().FindByIdentityNumber(ctx, "PNOEE-123456789").Return(user, nil)

				roleRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Role{}, nil)
				permissionRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Permission{}, nil)
				scopeRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID:          "PNOEE-123456789",
						Roles:       []string{},
						Permissions: []string{},
						Scope:       []string{},
					},
					models.AccessTokenExp,
				).Return("", assert.AnError)
			},
			expected: nil,
			err:      assert.AnError,
		},
		{
			name: "Failed to generate refresh token",
			before: func() {
				jwtService.EXPECT().Decode(refreshTokenValue).Return(&jwt.Payload{
					ID: "PNOEE-123456789",
				}, nil)
				userRepository.EXPECT().FindByIdentityNumber(ctx, "PNOEE-123456789").Return(user, nil)

				roleRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Role{}, nil)
				permissionRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Permission{}, nil)
				scopeRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID:          "PNOEE-123456789",
						Roles:       []string{},
						Permissions: []string{},
						Scope:       []string{},
					},
					models.AccessTokenExp,
				).Return("new-access-token", nil)
				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID: "PNOEE-123456789",
					},
					models.RefreshTokenExp,
				).Return("", assert.AnError)
			},
			expected: nil,
			err:      assert.AnError,
		},
		{
			name: "Failed to create user tokens",
			before: func() {
				jwtService.EXPECT().Decode(refreshTokenValue).Return(&jwt.Payload{
					ID: "PNOEE-123456789",
				}, nil)

				roleRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Role{}, nil)
				permissionRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Permission{}, nil)
				scopeRepository.EXPECT().FindByUserId(ctx, user.ID).Return([]models.Scope{}, nil)

				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID:          "PNOEE-123456789",
						Roles:       []string{},
						Permissions: []string{},
						Scope:       []string{},
					},
					models.AccessTokenExp,
				).Return("new-access-token", nil)
				jwtService.EXPECT().Generate(
					jwt.Payload{
						ID: "PNOEE-123456789",
					},
					models.RefreshTokenExp,
				).Return("new-refresh-token", nil)

				tokenRepository.EXPECT().Create(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			expected: nil,
			err:      assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Update(ctx, refreshTokenValue)

			if tt.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.AccessToken, result.AccessToken)
				assert.Equal(t, tt.expected.RefreshToken, result.RefreshToken)
			}
		})
	}
}

func Test_Tokens_FindById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	permissionRepository := repositories.NewMockPermissionRepository(ctrl)
	roleRepository := repositories.NewMockRoleRepository(ctrl)
	scopeRepository := repositories.NewMockScopeRepository(ctrl)
	tokenRepository := repositories.NewMockTokenRepository(ctrl)
	userRepository := repositories.NewMockUserRepository(ctrl)

	jwtService := jwt.NewMockJwt(ctrl)
	service := NewTokens(
		jwtService,
		permissionRepository,
		roleRepository,
		scopeRepository,
		tokenRepository,
		userRepository,
		log,
	)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		before   func()
		expected *models.Token
		error    error
	}{
		{
			name: "Success",
			before: func() {
				tokenRepository.EXPECT().FindById(ctx, id).Return(&models.Token{}, nil)
			},
			expected: &models.Token{},
			error:    nil,
		},
		{
			name: "Error",
			before: func() {
				tokenRepository.EXPECT().FindById(ctx, id).Return(nil, errors.ErrRecordNotFound)
			},
			expected: nil,
			error:    errors.ErrRecordNotFound,
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

func Test_Tokens_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	permissionRepository := repositories.NewMockPermissionRepository(ctrl)
	roleRepository := repositories.NewMockRoleRepository(ctrl)
	scopeRepository := repositories.NewMockScopeRepository(ctrl)
	tokenRepository := repositories.NewMockTokenRepository(ctrl)
	userRepository := repositories.NewMockUserRepository(ctrl)

	jwtService := jwt.NewMockJwt(ctrl)
	service := NewTokens(
		jwtService,
		permissionRepository,
		roleRepository,
		scopeRepository,
		tokenRepository,
		userRepository,
		log,
	)

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
				tokenRepository.EXPECT().Delete(ctx, id).Return(true, nil)
			},
			expected: true,
		},
		{
			name: "Error",
			before: func() {
				tokenRepository.EXPECT().Delete(ctx, id).Return(false, errors.ErrFailedToDeleteRecord)
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
