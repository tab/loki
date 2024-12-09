package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/errors"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

func Test_Session_GetStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := gomock.Any()
	authentication := services.NewMockAuthentication(ctrl)
	controller := NewSessionController(authentication)

	sessionId := "8fdb516d-1a82-43ba-b82d-be63df569b86"
	id, _ := uuid.Parse(sessionId)

	type result struct {
		response serializers.SessionSerializer
		error    serializers.ErrorSerializer
		status   string
		code     int
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				authentication.EXPECT().FindSessionById(ctx, sessionId).Return(&serializers.SessionSerializer{
					ID:     id,
					Status: "COMPLETED",
				}, nil)
			},
			expected: result{
				response: serializers.SessionSerializer{
					ID:     id,
					Status: "COMPLETED",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
		},
		{
			name: "Not found",
			before: func() {
				authentication.EXPECT().FindSessionById(ctx, sessionId).Return(nil, errors.ErrSessionNotFound)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: "session not found"},
				status: "404 Not Found",
				code:   http.StatusNotFound,
			},
		},
		{
			name: "Error",
			before: func() {
				authentication.EXPECT().FindSessionById(ctx, sessionId).Return(nil, fmt.Errorf("Redis error"))
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: "Redis error"},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/sessions/%s", sessionId), nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/auth/sessions/{id}", controller.GetStatus)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.expected.error.Error != "" {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.error.Error, response.Error)
			} else {
				var response serializers.SessionSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Session_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := gomock.Any()
	authentication := services.NewMockAuthentication(ctrl)
	controller := NewSessionController(authentication)

	sessionId := "8fdb516d-1a82-43ba-b82d-be63df569b86"
	id, _ := uuid.Parse(sessionId)

	type result struct {
		response serializers.UserSerializer
		error    serializers.ErrorSerializer
		status   string
		code     int
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				authentication.EXPECT().Complete(ctx, sessionId).Return(&serializers.UserSerializer{
					ID:             id,
					IdentityNumber: "PNOEE-30303039914",
					PersonalCode:   "30303039914",
					FirstName:      "TESTNUMBER",
					LastName:       "OK",
					AccessToken:    "aaa.bbb.ccc",
					RefreshToken:   "ddd.eee.fff",
				}, nil)
			},
			expected: result{
				response: serializers.UserSerializer{
					ID:             id,
					IdentityNumber: "PNOEE-30303039914",
					PersonalCode:   "30303039914",
					FirstName:      "TESTNUMBER",
					LastName:       "OK",
					AccessToken:    "aaa.bbb.ccc",
					RefreshToken:   "ddd.eee.fff",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
		},
		{
			name: "Error",
			before: func() {
				authentication.EXPECT().Complete(ctx, sessionId).Return(nil, fmt.Errorf("Failed to complete session"))
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: "Failed to complete session"},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/sessions/%s/authenticate", sessionId), nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/auth/sessions/{id}/authenticate", controller.Authenticate)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.expected.error.Error != "" {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.error.Error, response.Error)
			} else {
				var response serializers.UserSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}
