package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

func Test_TokensController_Refresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokens := services.NewMockTokens(ctrl)
	controller := NewTokensController(tokens)

	type result struct {
		response serializers.TokensSerializer
		error    serializers.ErrorSerializer
		status   string
		code     int
	}

	tests := []struct {
		name     string
		before   func()
		body     io.Reader
		expected result
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				tokens.EXPECT().Update(gomock.Any(), "refresh-token").Return(&models.User{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
				}, nil)
			},
			body: strings.NewReader(`{"refresh_token": "refresh-token"}`),
			expected: result{
				response: serializers.TokensSerializer{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Empty refresh token",
			before: func() {
			},
			body: strings.NewReader(`{"refresh_token": ""}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrInvalidToken.Error()},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Invalid refresh token",
			before: func() {
				tokens.EXPECT().Update(gomock.Any(), "invalid-token").Return(nil, errors.ErrInvalidToken)
			},
			body: strings.NewReader(`{"refresh_token": "invalid-token"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrInvalidToken.Error()},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
			error: true,
		},
		{
			name: "Error",
			before: func() {
				tokens.EXPECT().Update(gomock.Any(), "refresh-token").Return(nil, assert.AnError)
			},
			body: strings.NewReader(`{"refresh_token": "refresh-token"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: assert.AnError.Error()},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodPost, "/api/tokens/refresh", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/api/tokens/refresh", controller.Refresh)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.TokensSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}
