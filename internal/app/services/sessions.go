package services

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/config/logger"
)

type Sessions interface {
	Create(ctx context.Context, params *models.CreateSessionParams) (*models.Session, error)
	Update(ctx context.Context, params *models.UpdateSessionParams) (*models.Session, error)
	Delete(ctx context.Context, sessionId string) error
	FindById(ctx context.Context, sessionId string) (*models.Session, error)
}

type sessions struct {
	repository repositories.SessionRepository
	log        *logger.Logger
}

func NewSessions(repository repositories.SessionRepository, log *logger.Logger) Sessions {
	return &sessions{
		repository: repository,
		log:        log,
	}
}

func (s *sessions) Create(ctx context.Context, params *models.CreateSessionParams) (*models.Session, error) {
	id, err := uuid.Parse(params.SessionId)
	if err != nil {
		s.log.Error().Err(err).Msg("Invalid session ID format")
		return nil, err
	}

	err = s.repository.Create(ctx, &models.Session{
		ID:     id,
		Code:   params.Code,
		Status: models.SessionRunning,
	})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to create session")
		return nil, err
	}

	return &models.Session{
		ID:     id,
		Code:   params.Code,
		Status: models.SessionRunning,
	}, nil
}

func (s *sessions) Update(ctx context.Context, params *models.UpdateSessionParams) (*models.Session, error) {
	err := s.repository.Update(ctx, &models.Session{
		ID:     params.ID,
		UserId: params.UserId,
		Status: params.Status,
		Error:  params.Error,
	})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to update session")
		return nil, err
	}

	return &models.Session{
		ID:     params.ID,
		UserId: params.UserId,
		Status: params.Status,
		Error:  params.Error,
	}, nil
}

func (s *sessions) Delete(ctx context.Context, sessionId string) error {
	id, err := uuid.Parse(sessionId)
	if err != nil {
		s.log.Error().Err(err).Msg("Invalid session ID format")
		return err
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to delete session")
		return err
	}

	return nil
}

func (s *sessions) FindById(ctx context.Context, sessionId string) (*models.Session, error) {
	id, err := uuid.Parse(sessionId)
	if err != nil {
		s.log.Error().Err(err).Msg("Invalid session ID format")
		return nil, err
	}

	result, err := s.repository.FindById(ctx, id)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to find session")
		return nil, err
	}

	return result, nil
}
