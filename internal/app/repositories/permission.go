package repositories

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
)

type PermissionRepository interface {
	List(ctx context.Context, limit, offset int32) ([]models.Permission, int, error)
	Create(ctx context.Context, params db.CreatePermissionParams) (*models.Permission, error)
	Update(ctx context.Context, params db.UpdatePermissionParams) (*models.Permission, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	FindByUserId(ctx context.Context, id uuid.UUID) ([]models.Permission, error)
}

type permission struct {
	client postgres.Postgres
}

func NewPermissionRepository(client postgres.Postgres) PermissionRepository {
	return &permission{client: client}
}

func (p *permission) List(ctx context.Context, limit, offset int32) ([]models.Permission, int, error) {
	rows, err := p.client.Queries().FindPermissions(ctx, db.FindPermissionsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	permissions := make([]models.Permission, 0, len(rows))
	var total int

	if len(rows) > 0 {
		total = int(rows[0].Total)
	}

	for _, row := range rows {
		permissions = append(permissions, models.Permission{
			ID:          row.ID,
			Name:        row.Name.String,
			Description: row.Description.String,
		})
	}

	return permissions, total, err
}

func (p *permission) Create(ctx context.Context, params db.CreatePermissionParams) (*models.Permission, error) {
	result, err := p.client.Queries().CreatePermission(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.Permission{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, nil
}

func (p *permission) Update(ctx context.Context, params db.UpdatePermissionParams) (*models.Permission, error) {
	result, err := p.client.Queries().UpdatePermission(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.Permission{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, nil
}

func (p *permission) FindById(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	result, err := p.client.Queries().FindPermissionById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Permission{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, nil
}

func (p *permission) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	err := p.client.Queries().DeletePermission(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *permission) FindByUserId(ctx context.Context, id uuid.UUID) ([]models.Permission, error) {
	records, err := p.client.Queries().FindUserPermissions(ctx, id)
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
