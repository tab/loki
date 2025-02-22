package services

import (
	"context"

	"loki/internal/app/repositories"
)

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type health struct {
	repository repositories.HealthRepository
}

func NewHealthChecker(repository repositories.HealthRepository) HealthChecker {
	return &health{
		repository: repository,
	}
}

func (h *health) Ping(ctx context.Context) error {
	return h.repository.Ping(ctx)
}
