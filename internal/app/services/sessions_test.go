package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/config"
	"loki/internal/config/logger"
)

func Test_Sessions_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	repository := repositories.NewMockSessionRepository(ctrl)
	service := NewSessions(repository, log)

	id := uuid.MustParse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")
	sessionId := id.String()

	tests := []struct {
		name     string
		before   func()
		params   *models.CreateSessionParams
		expected *models.Session
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Create(ctx, &models.Session{
					ID:     id,
					Code:   "1234",
					Status: "RUNNING",
				}).Return(nil)
			},
			params: &models.CreateSessionParams{
				SessionId: sessionId,
				Code:      "1234",
			},
			expected: &models.Session{
				ID:     id,
				Code:   "1234",
				Status: "RUNNING",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().Create(ctx, &models.Session{
					ID:     id,
					Code:   "1234",
					Status: "RUNNING",
				}).Return(assert.AnError)
			},
			params: &models.CreateSessionParams{
				SessionId: sessionId,
				Code:      "1234",
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Create(ctx, tt.params)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Authentication_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	repository := repositories.NewMockSessionRepository(ctrl)
	service := NewSessions(repository, log)

	id := uuid.MustParse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")

	tests := []struct {
		name     string
		before   func()
		params   *models.UpdateSessionParams
		expected *models.Session
		error    error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Update(ctx, &models.Session{
					ID:     id,
					Status: "COMPLETE",
				}).Return(nil)
			},
			params: &models.UpdateSessionParams{
				ID:     id,
				Status: "COMPLETE",
			},
			expected: &models.Session{
				ID:     id,
				Status: "COMPLETE",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().Update(ctx, &models.Session{
					ID:     id,
					Status: "COMPLETE",
				}).Return(assert.AnError)
			},
			params: &models.UpdateSessionParams{
				ID:     id,
				Status: "COMPLETE",
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Update(ctx, tt.params)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Sessions_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	repository := repositories.NewMockSessionRepository(ctrl)
	service := NewSessions(repository, log)

	id := uuid.MustParse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")
	sessionId := id.String()

	tests := []struct {
		name   string
		before func()
		error  error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().Delete(ctx, id).Return(nil)
			},
			error: nil,
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().Delete(ctx, id).Return(assert.AnError)
			},
			error: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			err := service.Delete(ctx, sessionId)

			if tt.error != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Sessions_FindById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	repository := repositories.NewMockSessionRepository(ctrl)
	service := NewSessions(repository, log)

	id := uuid.MustParse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")
	sessionId := id.String()

	tests := []struct {
		name      string
		before    func()
		sessionId string
		expected  *models.Session
		error     error
	}{
		{
			name: "Success",
			before: func() {
				repository.EXPECT().FindById(ctx, id).Return(&models.Session{
					ID:     id,
					Status: "COMPLETE",
				}, nil)
			},
			sessionId: sessionId,
			expected: &models.Session{
				ID:     id,
				Status: "COMPLETE",
			},
		},
		{
			name: "Error",
			before: func() {
				repository.EXPECT().FindById(ctx, id).Return(nil, assert.AnError)
			},
			sessionId: sessionId,
			expected:  nil,
			error:     assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindById(ctx, tt.sessionId)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
