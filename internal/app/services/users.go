package services

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/repositories/db"
	"loki/pkg/logger"
)

type Users interface {
	List(ctx context.Context, pagination *Pagination) ([]models.User, uint64, error)
	Create(ctx context.Context, params *models.User) (*models.User, error)
	Update(ctx context.Context, params *models.User) (*models.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)

	FindByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error)
	FindUserDetailsById(ctx context.Context, id uuid.UUID) (*models.User, error)
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

func (u *users) List(ctx context.Context, pagination *Pagination) ([]models.User, uint64, error) {
	collection, total, err := u.repository.List(ctx, pagination.Limit(), pagination.Offset())

	if err != nil {
		return nil, 0, errors.ErrFailedToFetchResults
	}

	return collection, total, err
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

func (u *users) Update(ctx context.Context, params *models.User) (*models.User, error) {
	user, err := u.repository.Update(ctx, db.UpdateUserParams{
		ID:             params.ID,
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		RoleIDs:        params.RoleIDs,
		ScopeIDs:       params.ScopeIDs,
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

func (u *users) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	ok, err := u.repository.Delete(ctx, id)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (u *users) FindByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error) {
	user, err := u.repository.FindByIdentityNumber(ctx, identityNumber)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *users) FindUserDetailsById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := u.repository.FindUserDetailsById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
