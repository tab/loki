package services

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/repositories/db"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

type Tokens interface {
	List(ctx context.Context, pagination *Pagination) ([]models.Token, int, error)
	Create(ctx context.Context, userId uuid.UUID) (*models.User, error)
	Update(ctx context.Context, refreshToken string) (*models.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Token, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type tokens struct {
	jwt        jwt.Jwt
	permission repositories.PermissionRepository
	role       repositories.RoleRepository
	scope      repositories.ScopeRepository
	token      repositories.TokenRepository
	user       repositories.UserRepository
	log        *logger.Logger
}

func NewTokens(
	jwt jwt.Jwt,
	permission repositories.PermissionRepository,
	role repositories.RoleRepository,
	scope repositories.ScopeRepository,
	token repositories.TokenRepository,
	user repositories.UserRepository,
	log *logger.Logger,
) Tokens {
	return &tokens{
		jwt:        jwt,
		permission: permission,
		role:       role,
		scope:      scope,
		token:      token,
		user:       user,
		log:        log,
	}
}

func (t *tokens) List(ctx context.Context, pagination *Pagination) ([]models.Token, int, error) {
	collection, total, err := t.token.List(ctx, pagination.Limit(), pagination.Offset())

	if err != nil {
		return nil, 0, errors.ErrFailedToFetchResults
	}

	return collection, total, err
}

func (t *tokens) Create(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	user, err := t.user.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := t.generate(ctx, user)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             user.ID,
		IdentityNumber: user.IdentityNumber,
		PersonalCode:   user.PersonalCode,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
	}, nil
}

func (t *tokens) Update(ctx context.Context, refreshToken string) (*models.User, error) {
	payload, err := t.jwt.Decode(refreshToken)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to decode token")
		return nil, err
	}

	user, err := t.user.FindByIdentityNumber(ctx, payload.ID)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to find user")
		return nil, err
	}

	accessToken, refreshToken, err := t.generate(ctx, user)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             user.ID,
		IdentityNumber: user.IdentityNumber,
		PersonalCode:   user.PersonalCode,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
	}, nil
}

func (t *tokens) generate(ctx context.Context, user *models.User) (string, string, error) {
	userRoles, err := t.role.FindByUserId(ctx, user.ID)
	if err != nil {
		return "", "", err
	}
	roles := make([]string, 0, len(userRoles))
	for _, role := range userRoles {
		roles = append(roles, role.Name)
	}

	userPermissions, err := t.permission.FindByUserId(ctx, user.ID)
	if err != nil {
		return "", "", err
	}
	permissions := make([]string, 0, len(userPermissions))
	for _, permission := range userPermissions {
		permissions = append(permissions, permission.Name)
	}

	userScopes, err := t.scope.FindByUserId(ctx, user.ID)
	if err != nil {
		return "", "", err
	}
	scopes := make([]string, 0, len(userScopes))
	for _, scope := range userScopes {
		scopes = append(scopes, scope.Name)
	}

	accessToken, err := t.jwt.Generate(jwt.Payload{
		ID:          user.IdentityNumber,
		Roles:       roles,
		Permissions: permissions,
		Scope:       scopes,
	}, models.AccessTokenExp)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := t.jwt.Generate(jwt.Payload{
		ID: user.IdentityNumber,
	}, models.RefreshTokenExp)
	if err != nil {
		return "", "", err
	}

	_, err = t.token.Create(ctx, db.CreateTokensParams{
		UserID:            user.ID,
		AccessTokenValue:  accessToken,
		RefreshTokenValue: refreshToken,
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (t *tokens) FindById(ctx context.Context, id uuid.UUID) (*models.Token, error) {
	token, err := t.token.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *tokens) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	ok, err := t.token.Delete(ctx, id)
	if err != nil {
		return false, err
	}

	return ok, nil
}
