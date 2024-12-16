package spec

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"loki/internal/app/models"
)

func DbSeed(ctx context.Context, dsn string) error {
	roles := []string{models.AdminRoleType, models.ManagerRoleType, models.UserRoleType}

	for _, role := range roles {
		query := fmt.Sprintf("INSERT INTO roles (name) VALUES ('%s');", role)
		err := run(ctx, dsn, query)
		if err != nil {
			return err
		}
	}

	scopes := []string{models.SelfServiceType}

	for _, scope := range scopes {
		query := fmt.Sprintf("INSERT INTO scopes (name) VALUES ('%s');", scope)
		err := run(ctx, dsn, query)
		if err != nil {
			return err
		}
	}

	return nil
}

func TruncateTables(ctx context.Context, dsn string) error {
	tables := []string{
		"role_permissions",
		"user_roles",
		"user_scopes",
		"permissions",
		"roles",
		"scopes",
		"tokens",
		"users",
	}

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
