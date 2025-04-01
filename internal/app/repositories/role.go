package repositories

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories/db"
	"loki/internal/app/repositories/postgres"
)

type RoleRepository interface {
	List(ctx context.Context, limit, offset uint64) ([]models.Role, uint64, error)
	Create(ctx context.Context, params db.CreateRoleParams) (*models.Role, error)
	Update(ctx context.Context, params db.UpdateRoleParams) (*models.Role, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Role, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)

	FindByName(ctx context.Context, name string) (*models.Role, error)

	CreateUserRole(ctx context.Context, params db.CreateUserRoleParams) error
	FindByUserId(ctx context.Context, id uuid.UUID) ([]models.Role, error)

	FindRoleDetailsById(ctx context.Context, id uuid.UUID) (*models.Role, error)
}

type role struct {
	client postgres.Postgres
}

func NewRoleRepository(client postgres.Postgres) RoleRepository {
	return &role{client: client}
}

func (r *role) List(ctx context.Context, limit, offset uint64) ([]models.Role, uint64, error) {
	rows, err := r.client.Queries().FindRoles(ctx, db.FindRolesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	roles := make([]models.Role, 0, len(rows))
	var total uint64

	if len(rows) > 0 {
		total = rows[0].Total
	}

	for _, row := range rows {
		roles = append(roles, models.Role{
			ID:          row.ID,
			Name:        row.Name.String,
			Description: row.Description.String,
		})
	}

	return roles, total, err
}

func (r *role) Create(ctx context.Context, params db.CreateRoleParams) (*models.Role, error) {
	result, err := r.client.Queries().CreateRole(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.Role{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, nil
}

func (r *role) Update(ctx context.Context, params db.UpdateRoleParams) (*models.Role, error) {
	tx, err := r.client.Db().Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := r.client.Queries().WithTx(tx)

	result, err := q.UpdateRole(ctx, params)
	if err != nil {
		return nil, err
	}

	_, err = q.CreateRolePermissions(ctx, db.CreateRolePermissionsParams{
		RoleID:        result.ID,
		PermissionIds: params.PermissionIDs,
	})
	if err != nil {
		return nil, err
	}

	return &models.Role{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, tx.Commit(ctx)
}

func (r *role) FindById(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	result, err := r.client.Queries().FindRoleById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Role{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, nil
}

func (r *role) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	err := r.client.Queries().DeleteRole(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
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

func (r *role) FindRoleDetailsById(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	result, err := r.client.Queries().FindRoleDetailsById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Role{
		ID:            result.ID,
		Name:          result.Name,
		Description:   result.Description,
		PermissionIDs: result.PermissionIds,
	}, nil
}
