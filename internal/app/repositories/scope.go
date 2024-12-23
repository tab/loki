package repositories

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
)

type ScopeRepository interface {
	FindByName(ctx context.Context, name string) (*models.Scope, error)

	CreateUserScope(ctx context.Context, params db.CreateUserScopeParams) error
	FindByUserId(ctx context.Context, id uuid.UUID) ([]models.Scope, error)
}

type scope struct {
	client postgres.Postgres
}

func NewScopeRepository(client postgres.Postgres) ScopeRepository {
	return &scope{client: client}
}

func (s *scope) FindByName(ctx context.Context, name string) (*models.Scope, error) {
	record, err := s.client.Queries().FindScopeByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &models.Scope{
		ID:   record.ID,
		Name: record.Name,
	}, nil
}

func (s *scope) CreateUserScope(ctx context.Context, params db.CreateUserScopeParams) error {
	_, err := s.client.Queries().CreateUserScope(ctx, params)
	return err
}

func (s *scope) FindByUserId(ctx context.Context, id uuid.UUID) ([]models.Scope, error) {
	records, err := s.client.Queries().FindUserScopes(ctx, id)
	if err != nil {
		return nil, err
	}

	scopes := make([]models.Scope, 0, len(records))
	for _, record := range records {
		scopes = append(scopes, models.Scope{
			ID:   record.ID,
			Name: record.Name,
		})
	}

	return scopes, nil
}
