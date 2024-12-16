package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/app/repositories"
	"loki/internal/app/serializers"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

type Tokens interface {
	Refresh(ctx context.Context, userId uuid.UUID, refreshToken string) (*serializers.UserSerializer, error)
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

func (t *tokens) Refresh(ctx context.Context, userId uuid.UUID, token string) (*serializers.UserSerializer, error) {
	ok, err := t.jwt.Verify(token)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to decode token")
		return nil, err
	}

	if !ok {
		return nil, errors.ErrInvalidToken
	}

	user, err := t.database.FindUserById(ctx, userId)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to find user")
		return nil, err
	}

	accessToken, err := t.jwt.Generate(jwt.Payload{
		ID: user.IdentityNumber,
	}, models.AccessTokenExp)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to create access token")
		return nil, err
	}

	refreshToken, err := t.jwt.Generate(jwt.Payload{
		ID: user.IdentityNumber,
	}, models.RefreshTokenExp)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to create refresh token")
		return nil, err
	}

	_, err = t.database.CreateUserTokens(ctx, dto.CreateUserTokensParams{
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
