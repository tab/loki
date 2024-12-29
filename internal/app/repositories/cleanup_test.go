package repositories

import (
	"context"
	"log"
	"os"
	"testing"

	"loki/internal/config"
	"loki/pkg/spec"
)

func TestMain(m *testing.M) {
	if err := spec.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	if os.Getenv("GO_ENV") == "ci" {
		os.Exit(0)
	}

	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

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

	err := spec.TruncateTables(ctx, cfg.DatabaseDSN, tables)
	if err != nil {
		log.Fatalf("Error truncating tables: %v", err)
	}

	err = spec.DbSeed(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Error seeding database: %v", err)
	}

	code := m.Run()

	err = spec.TruncateTables(ctx, cfg.DatabaseDSN, tables)
	if err != nil {
		log.Fatalf("Error truncating tables: %v", err)
	}

	os.Exit(code)
}
