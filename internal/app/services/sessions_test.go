package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/serializers"
	"loki/pkg/logger"
)

func Test_Sessions_FindById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	redis := repositories.NewMockRedis(ctrl)
	log := logger.NewLogger()
	service := NewSessions(redis, log)

	id, _ := uuid.Parse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")
	sessionId := id.String()

	tests := []struct {
		name      string
		before    func()
		sessionId string
		expected  *serializers.SessionSerializer
		error     error
	}{
		{
			name: "Success",
			before: func() {
				redis.EXPECT().FindSessionById(ctx, id).Return(&models.Session{
					ID:     id,
					Status: "COMPLETE",
					Payload: models.SessionPayload{
						State:     "COMPLETE",
						Result:    "OK",
						Signature: "signature",
						Cert:      "certificate",
					},
				}, nil)
			},
			sessionId: sessionId,
			expected: &serializers.SessionSerializer{
				ID:     id,
				Status: "COMPLETE",
			},
		},
		{
			name: "Error",
			before: func() {
				redis.EXPECT().FindSessionById(ctx, id).Return(nil, assert.AnError)
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

func Test_Authentication_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	redis := repositories.NewMockRedis(ctrl)
	log := logger.NewLogger()
	service := NewSessions(redis, log)

	id, _ := uuid.Parse("5eab0e6a-c3e7-4526-a47e-398f0d31f514")

	tests := []struct {
		name     string
		before   func()
		params   models.Session
		expected *serializers.SessionSerializer
		error    error
	}{
		{
			name: "Success",
			before: func() {
				redis.EXPECT().UpdateSession(ctx, &models.Session{
					ID:     id,
					Status: "COMPLETE",
					Payload: models.SessionPayload{
						State:     "COMPLETE",
						Result:    "OK",
						Signature: "signature",
						Cert:      "certificate",
					},
				}).Return(nil)
			},
			params: models.Session{
				ID:     id,
				Status: "COMPLETE",
				Payload: models.SessionPayload{
					State:     "COMPLETE",
					Result:    "OK",
					Signature: "signature",
					Cert:      "certificate",
				},
			},
			expected: &serializers.SessionSerializer{
				ID:     id,
				Status: "COMPLETE",
			},
		},
		{
			name: "Error",
			before: func() {
				redis.EXPECT().UpdateSession(ctx, &models.Session{
					ID:     id,
					Status: "COMPLETE",
					Payload: models.SessionPayload{
						State:     "COMPLETE",
						Result:    "OK",
						Signature: "signature",
						Cert:      "certificate",
					},
				}).Return(assert.AnError)
			},
			params: models.Session{
				ID:     id,
				Status: "COMPLETE",
				Payload: models.SessionPayload{
					State:     "COMPLETE",
					Result:    "OK",
					Signature: "signature",
					Cert:      "certificate",
				},
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
