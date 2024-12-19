package services

import (
	"context"

	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/repositories/db"
	"loki/internal/app/serializers"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

type Tokens interface {
	Generate(ctx context.Context, user *models.User) (accessToken, refreshToken string, error error)
	Refresh(ctx context.Context, refreshToken string) (*serializers.UserSerializer, error)
}

type tokens struct {
	database repositories.Database
	jwt      jwt.Jwt
	log      *logger.Logger
}

func NewTokens(database repositories.Database, jwt jwt.Jwt, log *logger.Logger) Tokens {
	return &tokens{
		database: database,
		jwt:      jwt,
		log:      log,
	}
}

func (t *tokens) Generate(ctx context.Context, user *models.User) (accessToken, refreshToken string, error error) {
	return t.handleGenerateTokens(ctx, user)
}

func (t *tokens) Refresh(ctx context.Context, token string) (*serializers.UserSerializer, error) {
	payload, err := t.jwt.Decode(token)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to decode token")
		return nil, err
	}

	user, err := t.database.FindUserByIdentityNumber(ctx, payload.ID)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to find user")
		return nil, err
	}

	accessToken, refreshToken, err := t.handleGenerateTokens(ctx, user)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to generate tokens")
		return nil, err
	}

	return &serializers.UserSerializer{
		ID:             user.ID,
		IdentityNumber: user.IdentityNumber,
		PersonalCode:   user.PersonalCode,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
	}, nil
}

func (t *tokens) handleGenerateTokens(ctx context.Context, user *models.User) (accessToken, refreshToken string, error error) {
	userRoles, err := t.database.FindUserRoles(ctx, user.ID)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to find user roles")
		return "", "", err
	}
	roles := make([]string, 0, len(userRoles))
	for _, role := range userRoles {
		roles = append(roles, role.Name)
	}

	userPermissions, err := t.database.FindUserPermissions(ctx, user.ID)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to find user permissions")
		return "", "", err
	}
	permissions := make([]string, 0, len(userPermissions))
	for _, permission := range userPermissions {
		permissions = append(permissions, permission.Name)
	}

	userScopes, err := t.database.FindUserScopes(ctx, user.ID)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to find user scopes")
		return "", "", err
	}
	scopes := make([]string, 0, len(userScopes))
	for _, scope := range userScopes {
		scopes = append(scopes, scope.Name)
	}

	accessToken, err = t.jwt.Generate(jwt.Payload{
		ID:          user.IdentityNumber,
		Roles:       roles,
		Permissions: permissions,
		Scope:       scopes,
	}, models.AccessTokenExp)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to create access token")
		return "", "", err
	}

	refreshToken, err = t.jwt.Generate(jwt.Payload{
		ID: user.IdentityNumber,
	}, models.RefreshTokenExp)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to create refresh token")
		return "", "", err
	}

	_, err = t.database.CreateUserTokens(ctx, db.CreateTokensParams{
		UserID:            user.ID,
		AccessTokenValue:  accessToken,
		RefreshTokenValue: refreshToken,
	})
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to create user tokens in database")
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
