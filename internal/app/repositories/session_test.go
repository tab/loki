package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"loki/internal/app/models"
	"loki/internal/app/repositories/redis"
	"loki/internal/config"
)

func Test_SessionRepository_Create(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		RedisURI: os.Getenv("REDIS_URI"),
	}

	client, err := redis.NewRedisClient(cfg)
	assert.NoError(t, err)

	repo := NewSessionRepository(client)

	tests := []struct {
		name     string
		params   *models.Session
		expected error
	}{
		{
			name: "Success",
			params: &models.Session{
				ID:     uuid.New(),
				Status: "COMPLETED",
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.params)
			assert.Nil(t, err)
		})
	}
}

func Test_SessionRepository_Update(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		RedisURI: os.Getenv("REDIS_URI"),
	}

	client, err := redis.NewRedisClient(cfg)
	assert.NoError(t, err)

	repo := NewSessionRepository(client)

	id := uuid.MustParse("bf57e208-e6e7-4692-9de7-e75c1f8e5d52")

	tests := []struct {
		name     string
		before   func()
		params   *models.Session
		expected error
	}{
		{
			name: "Success",
			before: func() {
				err := repo.Create(ctx, &models.Session{
					ID:     id,
					Status: "RUNNING",
				})
				assert.NoError(t, err)
			},
			params: &models.Session{
				ID:     id,
				Status: "COMPLETED",
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.params)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func Test_SessionRepository_Delete(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		RedisURI: os.Getenv("REDIS_URI"),
	}

	client, err := redis.NewRedisClient(cfg)
	assert.NoError(t, err)

	repo := NewSessionRepository(client)

	id := uuid.MustParse("a29bfdbd-02d2-4f65-9601-d7309a0da16e")

	tests := []struct {
		name      string
		before    func()
		sessionId uuid.UUID
		expected  *models.Session
	}{
		{
			name: "Success",
			before: func() {
				err := repo.Create(ctx, &models.Session{
					ID:     id,
					Status: "RUNNING",
				})
				assert.NoError(t, err)
			},
			sessionId: id,
			expected: &models.Session{
				ID:     id,
				Status: "RUNNING",
			},
		},
		{
			name:      "Not found",
			before:    func() {},
			sessionId: uuid.New(),
			expected:  &models.Session{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			err := repo.Delete(ctx, tt.sessionId)
			assert.NoError(t, err)
		})
	}
}

func Test_SessionRepository_FindById(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		RedisURI: os.Getenv("REDIS_URI"),
	}

	client, err := redis.NewRedisClient(cfg)
	assert.NoError(t, err)

	repo := NewSessionRepository(client)

	id := uuid.MustParse("8fdb516d-1a82-43ba-b82d-be63df569b86")

	tests := []struct {
		name      string
		before    func()
		sessionId uuid.UUID
		expected  *models.Session
		error     bool
	}{
		{
			name: "Success",
			before: func() {
				err := repo.Create(ctx, &models.Session{
					ID:     id,
					Status: "RUNNING",
				})
				assert.NoError(t, err)
			},
			sessionId: id,
			expected: &models.Session{
				ID:     id,
				Status: "RUNNING",
			},
			error: false,
		},
		{
			name:      "Not found",
			before:    func() {},
			sessionId: uuid.New(),
			expected:  &models.Session{},
			error:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := repo.FindById(ctx, tt.sessionId)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Status, result.Status)
			}
		})
	}
}
