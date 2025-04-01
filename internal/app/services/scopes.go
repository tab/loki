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

type Scopes interface {
	List(ctx context.Context, pagination *Pagination) ([]models.Scope, uint64, error)
	Create(ctx context.Context, params *models.Scope) (*models.Scope, error)
	Update(ctx context.Context, params *models.Scope) (*models.Scope, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Scope, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type scopes struct {
	repository repositories.ScopeRepository
	log        *logger.Logger
}

func NewScopes(repository repositories.ScopeRepository, log *logger.Logger) Scopes {
	return &scopes{
		repository: repository,
		log:        log,
	}
}

func (s *scopes) List(ctx context.Context, pagination *Pagination) ([]models.Scope, uint64, error) {
	collection, total, err := s.repository.List(ctx, pagination.Limit(), pagination.Offset())

	if err != nil {
		return nil, 0, errors.ErrFailedToFetchResults
	}

	return collection, total, err
}

func (s *scopes) Create(ctx context.Context, params *models.Scope) (*models.Scope, error) {
	scope, err := s.repository.Create(ctx, db.CreateScopeParams{
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		return nil, err
	}

	return scope, nil
}

func (s *scopes) Update(ctx context.Context, params *models.Scope) (*models.Scope, error) {
	scope, err := s.repository.Update(ctx, db.UpdateScopeParams{
		ID:          params.ID,
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		return nil, err
	}

	return scope, nil
}

func (s *scopes) FindById(ctx context.Context, id uuid.UUID) (*models.Scope, error) {
	scope, err := s.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return scope, nil
}

func (s *scopes) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	ok, err := s.repository.Delete(ctx, id)
	if err != nil {
		return false, err
	}

	return ok, nil
}
