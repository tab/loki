package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/config"
)

const SessionTTL = 5 * time.Minute

type Redis interface {
	CreateSession(ctx context.Context, session *models.Session) error
	UpdateSession(ctx context.Context, session *models.Session) error
	FindSessionById(ctx context.Context, id uuid.UUID) (*models.Session, error)
	DeleteSessionByID(ctx context.Context, id uuid.UUID) error
}

type store struct {
	client *redis.Client
}

func NewRedis(cfg *config.Config) (Redis, error) {
	options, err := redis.ParseURL(cfg.RedisURI)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	return &store{client: client}, nil
}

func (r *store) CreateSession(ctx context.Context, session *models.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, session.ID.String(), data, SessionTTL).Err()
}

func (r *store) UpdateSession(ctx context.Context, session *models.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, session.ID.String(), data, SessionTTL).Err()
}

func (r *store) FindSessionById(ctx context.Context, sessionId uuid.UUID) (*models.Session, error) {
	data, err := r.client.Get(ctx, sessionId.String()).Result()
	if err != nil {
		return nil, errors.ErrSessionNotFound
	}

	var result models.Session
	if err = json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *store) DeleteSessionByID(ctx context.Context, sessionId uuid.UUID) error {
	return r.client.Del(ctx, sessionId.String()).Err()
}
