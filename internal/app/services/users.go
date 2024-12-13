package services

import (
	"context"
	"loki/pkg/logger"

	"loki/internal/app/repositories"
	"loki/internal/app/serializers"
)

type Users interface {
	FindByIdentityNumber(ctx context.Context, identityNumber string) (*serializers.UserSerializer, error)
}

type users struct {
	database repositories.Database
	log      *logger.Logger
}

func NewUsers(database repositories.Database, log *logger.Logger) Users {
	return &users{
		database: database,
		log:      log,
	}
}

func (u *users) FindByIdentityNumber(ctx context.Context, identityNumber string) (*serializers.UserSerializer, error) {
	user, err := u.database.FindUserByIdentityNumber(ctx, identityNumber)
	if err != nil {
		u.log.Error().Err(err).Msg("Failed to find user by identity number")
		return nil, err
	}

	return &serializers.UserSerializer{
		ID:             user.ID,
		IdentityNumber: user.IdentityNumber,
		PersonalCode:   user.PersonalCode,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
	}, nil
}
