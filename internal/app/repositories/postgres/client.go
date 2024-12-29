package postgres

import (
	"context"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"

	"loki/internal/app/repositories/db"
	"loki/internal/config"
)

type Postgres interface {
	Db() *pgxpool.Pool
	Queries() *db.Queries
}

type pgClient struct {
	db      *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresClient(cfg *config.Config) (Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	queries := db.New(pool)

	return &pgClient{
		db:      pool,
		queries: queries,
	}, nil
}

func (p *pgClient) Db() *pgxpool.Pool {
	return p.db
}

func (p *pgClient) Queries() *db.Queries {
	return p.queries
}
