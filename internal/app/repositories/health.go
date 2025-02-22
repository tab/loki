package repositories

import (
	"context"

	"loki/internal/app/repositories/postgres"
)

// HealthRepository is an interface for database health checks
type HealthRepository interface {
	Ping(ctx context.Context) error
}

type health struct {
	client postgres.Postgres
}

// NewHealthRepository creates a new health repository instance
func NewHealthRepository(client postgres.Postgres) HealthRepository {
	return &health{client: client}
}

// Ping checks the database connection
func (h *health) Ping(ctx context.Context) error {
	_, err := h.client.Queries().HealthCheck(ctx)
	return err
}
