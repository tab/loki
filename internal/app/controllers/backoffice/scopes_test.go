package backoffice

import (
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
	"loki/internal/app/models"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

func Test_Backoffice_Scopes_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scopes := services.NewMockScopes(ctrl)
	controller := NewScopesController(scopes)

	type result struct {
		response serializers.PaginationResponse[serializers.ScopeSerializer]
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
				scopes.EXPECT().List(gomock.Any(), gomock.Any()).Return([]models.Scope{
					{
						ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
						Name:        "sso-service",
						Description: "SSO-service scope",
					},
					{
						ID:          uuid.MustParse("10000000-1000-1000-2000-000000000002"),
						Name:        "self-service",
						Description: "Self-service scope",
					},
				}, 2, nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.ScopeSerializer]{
					Data: []serializers.ScopeSerializer{
						{
							ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
							Name:        "sso-service",
							Description: "SSO-service scope",
						},
						{
							ID:          uuid.MustParse("10000000-1000-1000-2000-000000000002"),
							Name:        "self-service",
							Description: "Self-service scope",
						},
					},
					Meta: serializers.PaginationMeta{
						Page:  1,
						Per:   25,
						Total: 2,
					},
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Empty",
			before: func() {
				scopes.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, 0, nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.ScopeSerializer]{
					Data: []serializers.ScopeSerializer{},
					Meta: serializers.PaginationMeta{
						Page:  1,
						Per:   25,
						Total: 0,
					},
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Error",
			before: func() {
				scopes.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, 0, assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/scopes", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/scopes", controller.List)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PaginationResponse[serializers.ScopeSerializer]
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Scopes_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scopes := services.NewMockScopes(ctrl)
	controller := NewScopesController(scopes)

	type result struct {
		response serializers.ScopeSerializer
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
				scopes.EXPECT().FindById(gomock.Any(), uuid.MustParse("10000000-1000-1000-2000-000000000001")).Return(&models.Scope{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        "sso-service",
					Description: "SSO-service scope",
				}, nil)
			},
			expected: result{
				response: serializers.ScopeSerializer{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        "sso-service",
					Description: "SSO-service scope",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Not found",
			before: func() {
				scopes.EXPECT().FindById(gomock.Any(), uuid.MustParse("10000000-1000-1000-2000-000000000001")).Return(&models.Scope{}, errors.ErrScopeNotFound)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrScopeNotFound.Error()},
				status: "404 Not Found",
				code:   http.StatusNotFound,
			},
		},
		{
			name: "Error",
			before: func() {
				scopes.EXPECT().FindById(gomock.Any(), uuid.MustParse("10000000-1000-1000-2000-000000000001")).Return(&models.Scope{}, assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/scopes/10000000-1000-1000-2000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/scopes/{id}", controller.Get)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.ScopeSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Scopes_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scopes := services.NewMockScopes(ctrl)
	controller := NewScopesController(scopes)

	type result struct {
		response serializers.ScopeSerializer
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
				scopes.EXPECT().Create(gomock.Any(), &models.Scope{
					Name:        "sso-service",
					Description: "SSO-service scope",
				}).Return(&models.Scope{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        "sso-service",
					Description: "SSO-service scope",
				}, nil)
			},
			body: strings.NewReader(`{"name":"sso-service","description":"SSO-service scope"}`),
			expected: result{
				response: serializers.ScopeSerializer{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        "sso-service",
					Description: "SSO-service scope",
				},
				status: "201 Created",
				code:   http.StatusCreated,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				scopes.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"name":"sso-service"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "empty description"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Error",
			before: func() {
				scopes.EXPECT().Create(gomock.Any(), &models.Scope{
					Name:        "sso-service",
					Description: "SSO-service scope",
				}).Return(&models.Scope{}, assert.AnError)
			},
			body: strings.NewReader(`{"name":"sso-service","description":"SSO-service scope"}`),
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

			req := httptest.NewRequest(http.MethodPost, "/api/backoffice/scopes", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/api/backoffice/scopes", controller.Create)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.ScopeSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Scopes_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scopes := services.NewMockScopes(ctrl)
	controller := NewScopesController(scopes)

	type result struct {
		response serializers.ScopeSerializer
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
				scopes.EXPECT().Update(gomock.Any(), &models.Scope{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        "sso-service",
					Description: "SSO-service scope",
				}).Return(&models.Scope{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        "sso-service",
					Description: "SSO-service scope",
				}, nil)
			},
			body: strings.NewReader(`{"name":"sso-service","description":"SSO-service scope"}`),
			expected: result{
				response: serializers.ScopeSerializer{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        "sso-service",
					Description: "SSO-service scope",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				scopes.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"name":"sso-service"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "empty description"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Error",
			before: func() {
				scopes.EXPECT().Update(gomock.Any(), &models.Scope{
					ID:          uuid.MustParse("10000000-1000-1000-2000-000000000001"),
					Name:        "sso-service",
					Description: "SSO-service scope",
				}).Return(&models.Scope{}, assert.AnError)
			},
			body: strings.NewReader(`{"name":"sso-service","description":"SSO-service scope"}`),
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

			req := httptest.NewRequest(http.MethodPut, "/api/backoffice/scopes/10000000-1000-1000-2000-000000000001", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Put("/api/backoffice/scopes/{id}", controller.Update)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.ScopeSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Scopes_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scopes := services.NewMockScopes(ctrl)
	controller := NewScopesController(scopes)

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
				scopes.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-2000-000000000001")).Return(true, nil)
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
				scopes.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-2000-000000000001")).Return(false, assert.AnError)
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

			req := httptest.NewRequest(http.MethodDelete, "/api/backoffice/scopes/10000000-1000-1000-2000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Delete("/api/backoffice/scopes/{id}", controller.Delete)
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
