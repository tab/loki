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
		{ID: uuid.MustParse("10000000-1000-1000-1000-000000000001"), Name: models.AdminRoleType},
		{ID: uuid.MustParse("10000000-1000-1000-1000-000000000002"), Name: models.ManagerRoleType},
		{ID: uuid.MustParse("10000000-1000-1000-1000-000000000003"), Name: models.UserRoleType},
	}
	for _, role := range roles {
		query := fmt.Sprintf("INSERT INTO roles (id, name) VALUES ('%s', '%s');", role.ID, role.Name)
		err := run(ctx, dsn, query)
		if err != nil {
			return err
		}
	}

	scopes := []models.Scope{
		{ID: uuid.MustParse("10000000-1000-1000-2000-000000000001"), Name: models.SelfServiceType},
	}
	for _, scope := range scopes {
		query := fmt.Sprintf("INSERT INTO scopes (id, name) VALUES ('%s', '%s');", scope.ID, scope.Name)
		err := run(ctx, dsn, query)
		if err != nil {
			return err
		}
	}

	permissions := []models.Permission{
		{ID: uuid.MustParse("10000000-1000-1000-3000-000000000001"), Name: "read:self"},
		{ID: uuid.MustParse("10000000-1000-1000-3000-000000000002"), Name: "write:self"},
		{ID: uuid.MustParse("10000000-1000-1000-3000-000000000003"), Name: "read:users"},
		{ID: uuid.MustParse("10000000-1000-1000-3000-000000000004"), Name: "write:users"},
	}
	for _, permission := range permissions {
		query := fmt.Sprintf("INSERT INTO permissions (id, name) VALUES ('%s', '%s');", permission.ID, permission.Name)
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

	_, err = db.Exec(ctx, query)
	if err != nil {
		return err
	}
	db.Close()

	return nil
}
