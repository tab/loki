package backoffice

import (
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
)

func Test_Backoffice_Tokens_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokens := services.NewMockTokens(ctrl)
	controller := NewTokensController(tokens)

	type result struct {
		response serializers.PaginationResponse[serializers.TokenSerializer]
		error    serializers.ErrorSerializer
		status   string
		code     int
	}

	tests := []struct {
		name     string
		before   func()
		expected result
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				tokens.EXPECT().List(gomock.Any(), gomock.Any()).Return([]models.Token{
					{
						ID:     uuid.MustParse("10000000-1000-1000-1111-000000000001"),
						UserId: uuid.MustParse("10000000-1000-1000-2222-000000000001"),
						Type:   models.AccessTokenType,
						Value:  "access-token-value",
					},
					{
						ID:     uuid.MustParse("10000000-1000-1000-1111-000000000002"),
						UserId: uuid.MustParse("10000000-1000-1000-2222-000000000002"),
						Type:   models.RefreshTokenType,
						Value:  "refresh-token-value",
					},
				}, 2, nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.TokenSerializer]{
					Data: []serializers.TokenSerializer{
						{
							ID:     uuid.MustParse("10000000-1000-1000-1111-000000000001"),
							UserId: uuid.MustParse("10000000-1000-1000-2222-000000000001"),
							Type:   models.AccessTokenType,
							Value:  "access-token-value",
						},
						{
							ID:     uuid.MustParse("10000000-1000-1000-1111-000000000002"),
							UserId: uuid.MustParse("10000000-1000-1000-2222-000000000002"),
							Type:   models.RefreshTokenType,
							Value:  "refresh-token-value",
						},
					},
					Meta: serializers.PaginationMeta{
						Page:  1,
						Per:   25,
						Total: 2,
					},
				},
				status: "200 OK",
				code:   200,
			},
			error: false,
		},
		{
			name: "Error",
			before: func() {
				tokens.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, 0, assert.AnError)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: assert.AnError.Error()},
				status: "422 Unprocessable Entity",
				code:   422,
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/tokens", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/tokens", controller.List)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PaginationResponse[serializers.TokenSerializer]
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Tokens_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokens := services.NewMockTokens(ctrl)
	controller := NewTokensController(tokens)

	type result struct {
		error  serializers.ErrorSerializer
		status string
		code   int
	}

	tests := []struct {
		name     string
		before   func()
		expected result
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				tokens.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-1111-000000000001")).Return(true, nil)
			},
			expected: result{
				status: "204 No Content",
				code:   http.StatusNoContent,
			},
			error: false,
		},
		{
			name: "Error",
			before: func() {
				tokens.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-1111-000000000001")).Return(false, assert.AnError)
			},
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

			req := httptest.NewRequest(http.MethodDelete, "/api/backoffice/tokens/10000000-1000-1000-1111-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Delete("/api/backoffice/tokens/{id}", controller.Delete)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}
