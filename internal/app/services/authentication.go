package services

import (
	"context"

	"loki/internal/app/models"
	"loki/internal/config"
	"loki/internal/config/logger"
)

const AuthenticationSuccess = "SUCCESS"

type Authentication interface {
	Complete(ctx context.Context, id string) (*models.User, error)
}

type authentication struct {
	cfg      *config.Config
	sessions Sessions
	tokens   Tokens
	log      *logger.Logger
}

func NewAuthentication(
	cfg *config.Config,
	sessions Sessions,
	tokens Tokens,
	log *logger.Logger,
) Authentication {
	return &authentication{
		cfg:      cfg,
		sessions: sessions,
		tokens:   tokens,
		log:      log,
	}
}

func (a *authentication) Complete(ctx context.Context, sessionId string) (*models.User, error) {
	session, err := a.sessions.FindById(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	user, err := a.tokens.Create(ctx, session.UserId)
	if err != nil {
		return nil, err
	}

	err = a.sessions.Delete(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
