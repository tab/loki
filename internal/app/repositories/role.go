package repositories

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
)

type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*models.Role, error)

	CreateUserRole(ctx context.Context, params db.CreateUserRoleParams) error
	FindByUserId(ctx context.Context, id uuid.UUID) ([]models.Role, error)
}

type role struct {
	client postgres.Postgres
}

func NewRoleRepository(client postgres.Postgres) RoleRepository {
	return &role{client: client}
}

func (r *role) FindByName(ctx context.Context, name string) (*models.Role, error) {
	record, err := r.client.Queries().FindRoleByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &models.Role{
		ID:   record.ID,
		Name: record.Name,
	}, nil
}

func (r *role) CreateUserRole(ctx context.Context, params db.CreateUserRoleParams) error {
	_, err := r.client.Queries().CreateUserRole(ctx, params)
	return err
}

func (r *role) FindByUserId(ctx context.Context, id uuid.UUID) ([]models.Role, error) {
	records, err := r.client.Queries().FindUserRoles(ctx, id)
	if err != nil {
		return nil, err
	}

	roles := make([]models.Role, 0, len(records))
	for _, record := range records {
		roles = append(roles, models.Role{
			ID:   record.ID,
			Name: record.Name,
		})
	}

	return roles, nil
}
