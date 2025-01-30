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

type Roles interface {
	List(ctx context.Context, pagination *Pagination) ([]models.Role, int, error)
	Create(ctx context.Context, params *models.Role) (*models.Role, error)
	Update(ctx context.Context, params *models.Role) (*models.Role, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Role, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)

	FindRoleDetailsById(ctx context.Context, id uuid.UUID) (*models.Role, error)
}

type roles struct {
	repository repositories.RoleRepository
	log        *logger.Logger
}

func NewRoles(repository repositories.RoleRepository, log *logger.Logger) Roles {
	return &roles{
		repository: repository,
		log:        log,
	}
}

func (r *roles) List(ctx context.Context, pagination *Pagination) ([]models.Role, int, error) {
	collection, total, err := r.repository.List(ctx, pagination.Limit(), pagination.Offset())

	if err != nil {
		return nil, 0, errors.ErrFailedToFetchResults
	}

	return collection, total, err
}

func (r *roles) Create(ctx context.Context, params *models.Role) (*models.Role, error) {
	role, err := r.repository.Create(ctx, db.CreateRoleParams{
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roles) Update(ctx context.Context, params *models.Role) (*models.Role, error) {
	role, err := r.repository.Update(ctx, db.UpdateRoleParams{
		ID:            params.ID,
		Name:          params.Name,
		Description:   params.Description,
		PermissionIDs: params.PermissionIDs,
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roles) FindById(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	role, err := r.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roles) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	ok, err := r.repository.Delete(ctx, id)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (r *roles) FindRoleDetailsById(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	role, err := r.repository.FindRoleDetailsById(ctx, id)
	if err != nil {
		return nil, err
	}

	return role, nil
}
