package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/app/repositories"
	"loki/internal/app/serializers"
	"loki/internal/config"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

const (
	AuthenticationSuccess = "SUCCESS"
	AuthenticationError   = "ERROR"

	AuthenticationTraceName = "authentication"
)

type Authentication interface {
	CreateSmartIdSession(ctx context.Context, params dto.CreateSmartIdSessionRequest) (*serializers.SessionSerializer, error)
	GetSmartIdSessionStatus(ctx context.Context, id uuid.UUID) (*dto.SmartIdProviderSessionStatusResponse, error)
	CreateMobileIdSession(ctx context.Context, params dto.CreateMobileIdSessionRequest) (*serializers.SessionSerializer, error)
	GetMobileIdSessionStatus(ctx context.Context, id uuid.UUID) (*dto.MobileIdProviderSessionStatusResponse, error)
	Complete(ctx context.Context, id string) (*serializers.UserSerializer, error)
}

type authentication struct {
	cfg              *config.Config
	smartIdProvider  SmartIdProvider
	smartIdQueue     chan<- *SmartIdQueue
	mobileIdProvider MobileIdProvider
	mobileIdQueue    chan<- *MobileIdQueue
	database         repositories.Database
	redis            repositories.Redis
	jwt              jwt.Jwt
	log              *logger.Logger
}

func NewAuthentication(
	cfg *config.Config,
	smartIdProvider SmartIdProvider,
	smartIdQueue chan *SmartIdQueue,
	mobileIdProvider MobileIdProvider,
	mobileIdQueue chan *MobileIdQueue,
	database repositories.Database,
	redis repositories.Redis,
	jwt jwt.Jwt,
	log *logger.Logger,
) Authentication {
	return &authentication{
		cfg:              cfg,
		smartIdProvider:  smartIdProvider,
		smartIdQueue:     smartIdQueue,
		mobileIdProvider: mobileIdProvider,
		mobileIdQueue:    mobileIdQueue,
		database:         database,
		redis:            redis,
		jwt:              jwt,
		log:              log,
	}
}

func (a *authentication) CreateSmartIdSession(ctx context.Context, params dto.CreateSmartIdSessionRequest) (response *serializers.SessionSerializer, error error) {
	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	result, err := a.smartIdProvider.CreateSession(ctx, params)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to initiate SmartId authentication")
		return nil, err
	}

	id, err := a.createSession(ctx, models.CreateSessionParams{
		SessionId: result.ID,
		Code:      result.Code,
	})
	if err != nil {
		return nil, err
	}

	a.smartIdQueue <- &SmartIdQueue{
		ID:      id,
		TraceId: traceId,
	}

	return &serializers.SessionSerializer{
		ID:   id,
		Code: result.Code,
	}, nil
}

func (a *authentication) GetSmartIdSessionStatus(_ context.Context, id uuid.UUID) (response *dto.SmartIdProviderSessionStatusResponse, error error) {
	result, err := a.smartIdProvider.GetSessionStatus(id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to get session status")
		return nil, err
	}

	return result, nil
}

func (a *authentication) CreateMobileIdSession(ctx context.Context, params dto.CreateMobileIdSessionRequest) (response *serializers.SessionSerializer, error error) {
	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	result, err := a.mobileIdProvider.CreateSession(ctx, params)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to initiate SmartId authentication")
		return nil, err
	}

	id, err := a.createSession(ctx, models.CreateSessionParams{
		SessionId: result.ID,
		Code:      result.Code,
	})
	if err != nil {
		return nil, err
	}

	a.mobileIdQueue <- &MobileIdQueue{
		ID:      id,
		TraceId: traceId,
	}

	return &serializers.SessionSerializer{
		ID:   id,
		Code: result.Code,
	}, nil
}

func (a *authentication) GetMobileIdSessionStatus(_ context.Context, id uuid.UUID) (response *dto.MobileIdProviderSessionStatusResponse, error error) {
	result, err := a.mobileIdProvider.GetSessionStatus(id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to get session status")
		return nil, err
	}

	return result, nil
}

func (a *authentication) Complete(ctx context.Context, sessionId string) (response *serializers.UserSerializer, error error) {
	id, err := uuid.Parse(sessionId)
	if err != nil {
		a.log.Error().Err(err).Msg("Invalid session ID format")
		return nil, err
	}

	result, err := a.redis.FindSessionById(ctx, id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to find session")
		return nil, err
	}

	user, err := a.createTokens(ctx, result.UserId)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	err = a.redis.DeleteSessionByID(ctx, id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to delete session")
		return nil, err
	}

	return &serializers.UserSerializer{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		IdentityNumber: user.IdentityNumber,
		PersonalCode:   user.PersonalCode,
		AccessToken:    user.AccessToken,
		RefreshToken:   user.RefreshToken,
	}, nil
}

func (a *authentication) createSession(ctx context.Context, params models.CreateSessionParams) (uuid.UUID, error) {
	id, err := uuid.Parse(params.SessionId)
	if err != nil {
		a.log.Error().Err(err).Msg("Invalid session ID format")
		return uuid.Nil, err
	}

	err = a.redis.CreateSession(ctx, &models.Session{
		ID:     id,
		Status: models.SessionRunning,
	})
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create session in redis")
		return uuid.Nil, err
	}

	return id, nil
}

func (a *authentication) createTokens(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := a.database.FindUserById(ctx, id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to find user")
		return nil, err
	}

	userRoles, err := a.database.FindUserRoles(ctx, id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to find user roles")
		return nil, err
	}
	roles := make([]string, 0, len(userRoles))
	for _, role := range userRoles {
		roles = append(roles, role.Name)
	}

	userPermissions, err := a.database.FindUserPermissions(ctx, id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to find user permissions")
		return nil, err
	}
	permissions := make([]string, 0, len(userPermissions))
	for _, permission := range userPermissions {
		permissions = append(permissions, permission.Name)
	}

	userScopes, err := a.database.FindUserScopes(ctx, id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to find user scopes")
		return nil, err
	}
	scopes := make([]string, 0, len(userScopes))
	for _, scope := range userScopes {
		scopes = append(scopes, scope.Name)
	}

	accessToken, err := a.jwt.Generate(jwt.Payload{
		ID:          user.IdentityNumber,
		Roles:       roles,
		Permissions: permissions,
		Scope:       scopes,
	}, models.AccessTokenExp)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create access token")
		return nil, err
	}

	refreshToken, err := a.jwt.Generate(jwt.Payload{
		ID: user.IdentityNumber,
	}, models.RefreshTokenExp)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create refresh token")
		return nil, err
	}

	result, err := a.database.CreateUserTokens(ctx, dto.CreateUserParams{
		IdentityNumber: user.IdentityNumber,
		AccessToken: dto.CreateTokenParams{
			Type:      models.AccessTokenType,
			Value:     accessToken,
			ExpiresAt: time.Now().Add(models.AccessTokenExp),
		},
		RefreshToken: dto.CreateTokenParams{
			Type:      models.RefreshTokenType,
			Value:     refreshToken,
			ExpiresAt: time.Now().Add(models.RefreshTokenExp),
		},
	})
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	return result, nil
}
