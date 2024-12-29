package repositories

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories/postgres"
)

type PermissionRepository interface {
	FindByUserId(ctx context.Context, id uuid.UUID) ([]models.Permission, error)
}

type permission struct {
	client postgres.Postgres
}

func NewPermissionRepository(client postgres.Postgres) PermissionRepository {
	return &permission{client: client}
}

func (r *permission) FindByUserId(ctx context.Context, id uuid.UUID) ([]models.Permission, error) {
	records, err := r.client.Queries().FindUserPermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	permissions := make([]models.Permission, 0, len(records))
	for _, record := range records {
		permissions = append(permissions, models.Permission{
			ID:   record.ID,
			Name: record.Name,
		})
	}

	return permissions, nil
}
