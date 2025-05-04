package workers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tab/mobileid"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/services"
	"loki/internal/config"
	"loki/internal/config/logger"
)

func Test_MobileIdWorker_Perform(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	sessionsMock := services.NewMockSessions(ctrl)
	usersMock := services.NewMockUsers(ctrl)
	workerMock := mobileid.NewMockWorker(ctrl)

	worker := NewMobileIdWorker(sessionsMock, usersMock, workerMock, log)

	id := uuid.MustParse("8fdb516d-1a82-43ba-b82d-be63df569b86")
	sessionId := id.String()
	traceId := uuid.New().String()
	userId := uuid.New()

	tests := []struct {
		name     string
		before   func()
		expected *models.Session
	}{
		{
			name: "Success",
			before: func() {
				resultChan := make(chan mobileid.Result, 1)
				resultChan <- mobileid.Result{
					Person: &mobileid.Person{
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					},
				}
				close(resultChan)

				workerMock.
					EXPECT().
					Process(ctx, sessionId).
					Return(resultChan)

				usersMock.
					EXPECT().
					Create(ctx, &models.User{
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					}).
					Return(&models.User{
						ID:             userId,
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					}, nil)

				sessionsMock.
					EXPECT().
					Update(ctx, &models.UpdateSessionParams{
						ID:     id,
						UserId: userId,
						Status: Success,
					}).
					Return(&models.Session{
						ID:     id,
						UserId: userId,
						Status: Success,
					}, nil)
			},
			expected: &models.Session{
				ID:     id,
				UserId: userId,
				Status: Success,
			},
		},
		{
			name: "Failed to get session status",
			before: func() {
				resultChan := make(chan mobileid.Result, 1)
				resultChan <- mobileid.Result{
					Err: assert.AnError,
				}
				close(resultChan)

				workerMock.
					EXPECT().
					Process(ctx, sessionId).
					Return(resultChan)

				sessionsMock.
					EXPECT().
					Update(ctx, &models.UpdateSessionParams{
						ID:     id,
						Status: Error,
						Error:  assert.AnError.Error(),
					}).Return(&models.Session{
					ID:     id,
					Status: Error,
					Error:  assert.AnError.Error(),
				}, nil)
			},
			expected: &models.Session{
				ID:     id,
				Status: Error,
				Error:  assert.AnError.Error(),
			},
		},
		{
			name: "Failed to create user",
			before: func() {
				resultChan := make(chan mobileid.Result, 1)
				resultChan <- mobileid.Result{
					Person: &mobileid.Person{
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					},
				}
				close(resultChan)

				workerMock.
					EXPECT().
					Process(ctx, sessionId).
					Return(resultChan)

				usersMock.
					EXPECT().
					Create(ctx, &models.User{
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					}).
					Return(nil, assert.AnError)

				sessionsMock.
					EXPECT().
					Update(ctx, &models.UpdateSessionParams{
						ID:     id,
						Status: Error,
						Error:  assert.AnError.Error(),
					}).
					Return(&models.Session{
						ID:     id,
						Status: Error,
						Error:  assert.AnError.Error(),
					}, nil)
			},
			expected: &models.Session{
				ID:     id,
				Status: Error,
				Error:  assert.AnError.Error(),
			},
		},
		{
			name: "Failed to update session",
			before: func() {
				resultChan := make(chan mobileid.Result, 1)
				resultChan <- mobileid.Result{
					Person: &mobileid.Person{
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					},
				}
				close(resultChan)

				workerMock.
					EXPECT().
					Process(ctx, sessionId).
					Return(resultChan)

				usersMock.
					EXPECT().
					Create(ctx, &models.User{
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					}).
					Return(&models.User{
						ID:             userId,
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					}, nil)

				sessionsMock.
					EXPECT().
					Update(ctx, &models.UpdateSessionParams{
						ID:     id,
						UserId: userId,
						Status: Success,
					}).
					Return(nil, assert.AnError)
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			result := worker.Perform(ctx, id, traceId)

			assert.Equal(t, tt.expected, result)
		})
	}
}
