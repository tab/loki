package services

import (
	"context"

	"github.com/google/uuid"

	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/serializers"
	"loki/pkg/logger"
)

type Sessions interface {
	FindById(ctx context.Context, sessionId string) (*serializers.SessionSerializer, error)
	Update(ctx context.Context, params models.Session) (*serializers.SessionSerializer, error)
}

type sessions struct {
	redis repositories.Redis
	log   *logger.Logger
}

func NewSessions(redis repositories.Redis, log *logger.Logger) Sessions {
	return &sessions{
		redis: redis,
		log:   log,
	}
}

func (s *sessions) FindById(ctx context.Context, sessionId string) (*serializers.SessionSerializer, error) {
	id, err := uuid.Parse(sessionId)
	if err != nil {
		s.log.Error().Err(err).Msg("Invalid session ID format")
		return nil, err
	}

	result, err := s.redis.FindSessionById(ctx, id)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to find session")
		return nil, err
	}

	return &serializers.SessionSerializer{
		ID:     result.ID,
		Code:   result.Code,
		Status: result.Status,
		Error:  result.Error,
	}, nil
}

func (s *sessions) Update(ctx context.Context, params models.Session) (*serializers.SessionSerializer, error) {
	err := s.redis.UpdateSession(ctx, &models.Session{
		ID:     params.ID,
		Status: params.Status,
		Error:  params.Error,
		Payload: models.SessionPayload{
			State:     params.Payload.State,
			Result:    params.Payload.Result,
			Signature: params.Payload.Signature,
			Cert:      params.Payload.Cert,
		},
	})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to update session")
		return nil, err
	}

	return &serializers.SessionSerializer{
		ID:     params.ID,
		Status: params.Status,
	}, nil
}
