package workers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tab/smartid"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/services"
	"loki/pkg/logger"
)

func Test_SmartIdWorker_Perform(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	sessionsMock := services.NewMockSessions(ctrl)
	usersMock := services.NewMockUsers(ctrl)
	workerMock := smartid.NewMockWorker(ctrl)
	log := logger.NewLogger()

	worker := NewSmartIdWorker(sessionsMock, usersMock, workerMock, log)

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
				resultChan := make(chan smartid.Result, 1)
				resultChan <- smartid.Result{
					Person: &smartid.Person{
						IdentityNumber: "PNOEE-30303039914",
						PersonalCode:   "30303039914",
						FirstName:      "TESTNUMBER",
						LastName:       "OK",
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
						IdentityNumber: "PNOEE-30303039914",
						PersonalCode:   "30303039914",
						FirstName:      "TESTNUMBER",
						LastName:       "OK",
					}).
					Return(&models.User{
						ID:             userId,
						IdentityNumber: "PNOEE-30303039914",
						PersonalCode:   "30303039914",
						FirstName:      "TESTNUMBER",
						LastName:       "OK",
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
				resultChan := make(chan smartid.Result, 1)
				resultChan <- smartid.Result{
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
				resultChan := make(chan smartid.Result, 1)
				resultChan <- smartid.Result{
					Person: &smartid.Person{
						IdentityNumber: "PNOEE-30303039914",
						PersonalCode:   "30303039914",
						FirstName:      "TESTNUMBER",
						LastName:       "OK",
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
						IdentityNumber: "PNOEE-30303039914",
						PersonalCode:   "30303039914",
						FirstName:      "TESTNUMBER",
						LastName:       "OK",
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
				resultChan := make(chan smartid.Result, 1)
				resultChan <- smartid.Result{
					Person: &smartid.Person{
						IdentityNumber: "PNOEE-30303039914",
						PersonalCode:   "30303039914",
						FirstName:      "TESTNUMBER",
						LastName:       "OK",
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
						IdentityNumber: "PNOEE-30303039914",
						PersonalCode:   "30303039914",
						FirstName:      "TESTNUMBER",
						LastName:       "OK",
					}).
					Return(&models.User{
						ID:             userId,
						IdentityNumber: "PNOEE-30303039914",
						PersonalCode:   "30303039914",
						FirstName:      "TESTNUMBER",
						LastName:       "OK",
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
