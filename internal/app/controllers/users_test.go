package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
	"loki/internal/config/middlewares"
)

func Test_UsersController_Me(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users)

	identityNumber := "PNOEE-123456789"

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	type result struct {
		response serializers.UserSerializer
		error    serializers.ErrorSerializer
		status   string
		code     int
	}

	tests := []struct {
		name        string
		currentUser *models.User
		expected    result
	}{
		{
			name: "Success",
			currentUser: &models.User{
				ID:             id,
				IdentityNumber: identityNumber,
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			},
			expected: result{
				response: serializers.UserSerializer{
					ID:             id,
					IdentityNumber: identityNumber,
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
		},
		{
			name:        "Unauthorized",
			currentUser: nil,
			expected: result{
				error:  serializers.ErrorSerializer{Error: "unauthorized"},
				status: "401 Unauthorized",
				code:   http.StatusUnauthorized,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
			if tt.currentUser != nil {
				ctx := context.WithValue(req.Context(), middlewares.CurrentUser{}, tt.currentUser)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/me", controller.Me)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.currentUser != nil {
				var response serializers.UserSerializer
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			} else {
				var response serializers.ErrorSerializer
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}
