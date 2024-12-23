package services

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/repositories/db"
	"loki/pkg/logger"
)

type Users interface {
	Create(ctx context.Context, params *models.User) (*models.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error)
}

type users struct {
	repository repositories.UserRepository
	log        *logger.Logger
}

func NewUsers(repository repositories.UserRepository, log *logger.Logger) Users {
	return &users{
		repository: repository,
		log:        log,
	}
}

func (u *users) Create(ctx context.Context, params *models.User) (*models.User, error) {
	user, err := u.repository.Create(ctx, db.CreateUserParams{
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

func (u *users) FindById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := u.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *users) FindByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error) {
	user, err := u.repository.FindByIdentityNumber(ctx, identityNumber)
	if err != nil {
		return nil, err
	}

	return user, nil
}
