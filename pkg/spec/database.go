package spec

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"loki/internal/app/models"
)

func DbSeed(ctx context.Context, dsn string) error {
	roles := []models.Role{
		{
			ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
			Name:        models.AdminRoleType,
			Description: "Admin role",
		},
		{
			ID:          uuid.MustParse("10000000-1000-1000-1000-000000000002"),
			Name:        models.ManagerRoleType,
			Description: "Manager role",
		},
		{
			ID:          uuid.MustParse("10000000-1000-1000-1000-000000000003"),
			Name:        models.UserRoleType,
			Description: "User role",
		},
	}
	for _, role := range roles {
		query := fmt.Sprintf("INSERT INTO roles (id, name, description) VALUES ('%s', '%s', '%s');", role.ID, role.Name, role.Description)
		err := run(ctx, dsn, query)
		if err != nil {
			return err
		}
	}

	scopes := []models.Scope{
		{
			ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
			Name:        models.SsoServiceType,
			Description: "SSO-service scope",
		},
		{
			ID:          uuid.MustParse("10000000-1000-1000-2000-000000000002"),
			Name:        models.SelfServiceType,
			Description: "Self-service scope",
		},
	}
	for _, scope := range scopes {
		query := fmt.Sprintf("INSERT INTO scopes (id, name, description) VALUES ('%s', '%s', '%s');", scope.ID, scope.Name, scope.Description)
		err := run(ctx, dsn, query)
		if err != nil {
			return err
		}
	}

	permissions := []models.Permission{
		{
			ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
			Name:        "read:self",
			Description: "Read own data",
		},
		{
			ID:          uuid.MustParse("10000000-1000-1000-3000-000000000002"),
			Name:        "write:self",
			Description: "Write own data",
		},
		{
			ID:          uuid.MustParse("10000000-1000-1000-3000-000000000003"),
			Name:        "read:users",
			Description: "Read user data",
		},
		{
			ID:          uuid.MustParse("10000000-1000-1000-3000-000000000004"),
			Name:        "write:users",
			Description: "Write user data",
		},
	}
	for _, permission := range permissions {
		query := fmt.Sprintf("INSERT INTO permissions (id, name, description) VALUES ('%s', '%s', '%s');", permission.ID, permission.Name, permission.Description)
		err := run(ctx, dsn, query)
		if err != nil {
			return err
		}
	}

	rolePermissionMapping := map[string][]string{
		models.AdminRoleType:   {"read:self", "write:self", "read:users", "write:users"},
		models.ManagerRoleType: {"read:self", "write:self", "read:users"},
		models.UserRoleType:    {"read:self", "write:self"},
	}

	roleMap := make(map[string]uuid.UUID)
	for _, r := range roles {
		roleMap[r.Name] = r.ID
	}

	permMap := make(map[string]uuid.UUID)
	for _, p := range permissions {
		permMap[p.Name] = p.ID
	}

	for rName, permNames := range rolePermissionMapping {
		for _, pName := range permNames {
			query := fmt.Sprintf(
				"INSERT INTO role_permissions (role_id, permission_id) VALUES ('%s', '%s');",
				roleMap[rName],
				permMap[pName],
			)
			err := run(ctx, dsn, query)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TruncateTables(ctx context.Context, dsn string, tables []string) error {
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)
		err := run(ctx, dsn, query)
		if err != nil {
			return err
		}
	}

	return nil
}

func run(ctx context.Context, dsn string, query string) error {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
