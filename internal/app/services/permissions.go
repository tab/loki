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

type Permissions interface {
	List(ctx context.Context, pagination *Pagination) ([]models.Permission, uint64, error)
	Create(ctx context.Context, params *models.Permission) (*models.Permission, error)
	Update(ctx context.Context, params *models.Permission) (*models.Permission, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type permissions struct {
	repository repositories.PermissionRepository
	log        *logger.Logger
}

func NewPermissions(repository repositories.PermissionRepository, log *logger.Logger) Permissions {
	return &permissions{
		repository: repository,
		log:        log,
	}
}

func (p *permissions) List(ctx context.Context, pagination *Pagination) ([]models.Permission, uint64, error) {
	collection, total, err := p.repository.List(ctx, pagination.Limit(), pagination.Offset())

	if err != nil {
		return nil, 0, errors.ErrFailedToFetchResults
	}

	return collection, total, err
}

func (p *permissions) Create(ctx context.Context, params *models.Permission) (*models.Permission, error) {
	permission, err := p.repository.Create(ctx, db.CreatePermissionParams{
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (p *permissions) Update(ctx context.Context, params *models.Permission) (*models.Permission, error) {
	permission, err := p.repository.Update(ctx, db.UpdatePermissionParams{
		ID:          params.ID,
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (p *permissions) FindById(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	permission, err := p.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (p *permissions) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	ok, err := p.repository.Delete(ctx, id)
	if err != nil {
		return false, err
	}

	return ok, nil
}
