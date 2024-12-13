package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/errors"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
	"loki/internal/config/middlewares"
)

func Test_TokensController_Refresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokens := services.NewMockTokens(ctrl)
	controller := NewTokensController(tokens)

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
		before      func()
		currentUser *serializers.UserSerializer
		body        io.Reader
		expected    result
	}{
		{
			name: "Success",
			before: func() {
				tokens.EXPECT().Refresh(gomock.Any(), id, "refresh-token").Return(&serializers.UserSerializer{
					ID:             id,
					IdentityNumber: identityNumber,
					PersonalCode:   "123456789",
					AccessToken:    "new-access-token",
					RefreshToken:   "new-refresh-token",
				}, nil)
			},
			currentUser: &serializers.UserSerializer{
				ID:             id,
				IdentityNumber: identityNumber,
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			},
			body: strings.NewReader(`{"refresh_token": "refresh-token"}`),
			expected: result{
				response: serializers.UserSerializer{
					ID:             id,
					IdentityNumber: identityNumber,
					PersonalCode:   "123456789",
					AccessToken:    "new-access-token",
					RefreshToken:   "new-refresh-token",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
		},
		{
			name:        "Unauthorized",
			before:      func() {},
			currentUser: nil,
			body:        strings.NewReader(`{"refresh_token": "refresh-token"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrUnauthorized.Error()},
				status: "401 Unauthorized",
				code:   http.StatusUnauthorized,
			},
		},
		{
			name: "Empty refresh token",
			before: func() {
			},
			currentUser: &serializers.UserSerializer{
				ID:             id,
				IdentityNumber: identityNumber,
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			},
			body: strings.NewReader(`{"refresh_token": ""}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrInvalidToken.Error()},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name: "Invalid refresh token",
			before: func() {
				tokens.EXPECT().Refresh(gomock.Any(), id, "invalid-token").Return(nil, errors.ErrInvalidToken)
			},
			currentUser: &serializers.UserSerializer{
				ID:             id,
				IdentityNumber: identityNumber,
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			},
			body: strings.NewReader(`{"refresh_token": "invalid-token"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrInvalidToken.Error()},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Error",
			before: func() {
				tokens.EXPECT().Refresh(gomock.Any(), id, "refresh-token").Return(nil, assert.AnError)
			},
			currentUser: &serializers.UserSerializer{
				ID:             id,
				IdentityNumber: identityNumber,
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			},
			body: strings.NewReader(`{"refresh_token": "refresh-token"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: assert.AnError.Error()},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}

			req := httptest.NewRequest(http.MethodPost, "/api/tokens/refresh", tt.body)
			if tt.currentUser != nil {
				ctx := context.WithValue(req.Context(), middlewares.CurrentUser{}, tt.currentUser)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/api/tokens/refresh", controller.Refresh)
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
