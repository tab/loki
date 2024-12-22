package services

import (
	"context"

	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/repositories/db"
	"loki/internal/app/serializers"
	"loki/pkg/logger"
)

type Users interface {
	Create(ctx context.Context, params *models.User) (*models.User, error)
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

func (u *users) Create(ctx context.Context, params *models.User) (*models.User, error) {
	user, err := u.database.CreateUser(ctx, db.CreateUserParams{
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *users) FindByIdentityNumber(ctx context.Context, identityNumber string) (*serializers.UserSerializer, error) {
	user, err := u.database.FindUserByIdentityNumber(ctx, identityNumber)
	if err != nil {
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
