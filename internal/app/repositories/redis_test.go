package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"loki/internal/app/models"
	"loki/internal/config"
)

func Test_Redis_CreateSession(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		RedisURI: os.Getenv("REDIS_URI"),
	}

	repo, err := NewRedis(cfg)
	assert.NoError(t, err)

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
			err := repo.CreateSession(ctx, tt.params)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func Test_Redis_UpdateSession(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		RedisURI: os.Getenv("REDIS_URI"),
	}

	repo, err := NewRedis(cfg)
	assert.NoError(t, err)

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
				session := &models.Session{
					ID:     id,
					Status: "RUNNING",
				}
				err := repo.CreateSession(ctx, session)
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
			err := repo.UpdateSession(ctx, tt.params)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func Test_Redis_FindSessionById(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		RedisURI: os.Getenv("REDIS_URI"),
	}

	repo, err := NewRedis(cfg)
	assert.NoError(t, err)

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
				session := &models.Session{
					ID:     id,
					Status: "RUNNING",
				}
				err := repo.CreateSession(ctx, session)
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

			session, err := repo.FindSessionById(ctx, tt.sessionId)

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.ID, session.ID)
				assert.Equal(t, tt.expected.Status, session.Status)
			}
		})
	}
}

func Test_Redis_DeleteSessionByID(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		RedisURI: os.Getenv("REDIS_URI"),
	}

	repo, err := NewRedis(cfg)
	assert.NoError(t, err)

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
				session := &models.Session{
					ID:     id,
					Status: "RUNNING",
				}
				err := repo.CreateSession(ctx, session)
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

			err := repo.DeleteSessionByID(ctx, tt.sessionId)
			assert.NoError(t, err)
		})
	}
}
