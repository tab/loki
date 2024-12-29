package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/repositories/redis"
)

const SessionTTL = 5 * time.Minute

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindById(ctx context.Context, id uuid.UUID) (*models.Session, error)
}

type session struct {
	client redis.Redis
}

func NewSessionRepository(client redis.Redis) SessionRepository {
	return &session{client: client}
}

func (s *session) Create(ctx context.Context, session *models.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return s.client.Connection().Set(ctx, session.ID.String(), data, SessionTTL).Err()
}

func (s *session) Update(ctx context.Context, session *models.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return s.client.Connection().Set(ctx, session.ID.String(), data, SessionTTL).Err()
}

func (s *session) Delete(ctx context.Context, id uuid.UUID) error {
	return s.client.Connection().Del(ctx, id.String()).Err()
}

func (s *session) FindById(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	data, err := s.client.Connection().Get(ctx, id.String()).Result()
	if err != nil {
		return nil, errors.ErrSessionNotFound
	}

	var result models.Session
	if err = json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
