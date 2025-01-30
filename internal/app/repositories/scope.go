package repositories

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
)

type ScopeRepository interface {
	List(ctx context.Context, limit, offset int32) ([]models.Scope, int, error)
	Create(ctx context.Context, params db.CreateScopeParams) (*models.Scope, error)
	Update(ctx context.Context, params db.UpdateScopeParams) (*models.Scope, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Scope, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)

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

func (s *scope) List(ctx context.Context, limit, offset int32) ([]models.Scope, int, error) {
	rows, err := s.client.Queries().FindScopes(ctx, db.FindScopesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	roles := make([]models.Scope, 0, len(rows))
	var total int

	if len(rows) > 0 {
		total = int(rows[0].Total)
	}

	for _, row := range rows {
		roles = append(roles, models.Scope{
			ID:          row.ID,
			Name:        row.Name.String,
			Description: row.Description.String,
		})
	}

	return roles, total, err
}

func (s *scope) Create(ctx context.Context, params db.CreateScopeParams) (*models.Scope, error) {
	result, err := s.client.Queries().CreateScope(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.Scope{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, nil
}

func (s *scope) Update(ctx context.Context, params db.UpdateScopeParams) (*models.Scope, error) {
	result, err := s.client.Queries().UpdateScope(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.Scope{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, nil
}

func (s *scope) FindById(ctx context.Context, id uuid.UUID) (*models.Scope, error) {
	result, err := s.client.Queries().FindScopeById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Scope{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, nil
}

func (s *scope) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	err := s.client.Queries().DeleteScope(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
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
