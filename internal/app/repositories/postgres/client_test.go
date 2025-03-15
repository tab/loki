package postgres

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"loki/internal/app/repositories/db"
	"loki/internal/config"
)

func Test_NewPostgresClient(t *testing.T) {
	type args struct {
		cfg *config.Config
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{
				cfg: &config.Config{
					DatabaseDSN: "postgres://localhost:5432",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewPostgresClient(tt.args.cfg)
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func Test_Postgres_Db(t *testing.T) {
	type args struct {
		cfg *config.Config
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{
				cfg: &config.Config{
					DatabaseDSN: "postgres://localhost:5432",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewPostgresClient(tt.args.cfg)
			assert.NoError(t, err)

			pool := client.Db()

			assert.NotNil(t, pool)
			assert.IsType(t, &pgxpool.Pool{}, pool)
		})
	}
}

func Test_Postgres_Queries(t *testing.T) {
	type args struct {
		cfg *config.Config
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{
				cfg: &config.Config{
					DatabaseDSN: "postgres://localhost:5432",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewPostgresClient(tt.args.cfg)
			assert.NoError(t, err)

			queries := client.Queries()

			assert.NotNil(t, queries)
			assert.IsType(t, &db.Queries{}, queries)
		})
	}
}
