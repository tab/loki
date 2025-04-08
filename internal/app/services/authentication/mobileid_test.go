package authentication

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tab/mobileid"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/app/services"
	"loki/internal/app/workers"
	"loki/internal/config"
	"loki/internal/config/logger"
)

func Test_MobileId_CreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	clientMock := mobileid.NewMockClient(ctrl)
	sessionsMock := services.NewMockSessions(ctrl)
	usersMock := services.NewMockUsers(ctrl)
	workerMock := workers.NewMockMobileIdWorker(ctrl)

	service := NewMobileId(clientMock, sessionsMock, usersMock, workerMock, log)

	personalCode := "51307149560"
	phoneNumber := "+37269930366"

	id := uuid.MustParse("8fdb516d-1a82-43ba-b82d-be63df569b86")
	sessionId := id.String()

	tests := []struct {
		name     string
		before   func()
		params   dto.CreateMobileIdSessionRequest
		expected *models.Session
		error    error
	}{
		{
			name: "Success",
			before: func() {
				clientMock.EXPECT().CreateSession(ctx, phoneNumber, personalCode).Return(&mobileid.Session{
					Id:   sessionId,
					Code: "1234",
				}, nil)

				sessionsMock.EXPECT().Create(ctx, &models.CreateSessionParams{
					SessionId: sessionId,
					Code:      "1234",
				}).Return(&models.Session{
					ID:   id,
					Code: "1234",
				}, nil)

				workerMock.EXPECT().Perform(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			},
			params: dto.CreateMobileIdSessionRequest{
				PersonalCode: personalCode,
				PhoneNumber:  phoneNumber,
			},
			expected: &models.Session{
				ID:     id,
				Code:   "1234",
				Status: models.SessionRunning,
			},
			error: nil,
		},
		{
			name: "Error to create smart-id session",
			before: func() {
				clientMock.EXPECT().CreateSession(ctx, phoneNumber, personalCode).Return(nil, assert.AnError)
			},
			params: dto.CreateMobileIdSessionRequest{
				PersonalCode: personalCode,
				PhoneNumber:  phoneNumber,
			},
			expected: nil,
			error:    assert.AnError,
		},
		{
			name: "Error to save smart-id session",
			before: func() {
				clientMock.EXPECT().CreateSession(ctx, phoneNumber, personalCode).Return(&mobileid.Session{
					Id:   sessionId,
					Code: "1234",
				}, nil)

				sessionsMock.EXPECT().Create(ctx, &models.CreateSessionParams{
					SessionId: sessionId,
					Code:      "1234",
				}).Return(nil, assert.AnError)
			},
			params: dto.CreateMobileIdSessionRequest{
				PersonalCode: personalCode,
				PhoneNumber:  phoneNumber,
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.CreateSession(ctx, tt.params)

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
