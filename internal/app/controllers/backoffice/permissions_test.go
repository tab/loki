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

func Test_Backoffice_Permissions_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissions := services.NewMockPermissions(ctrl)
	controller := NewPermissionsController(permissions)

	type result struct {
		response serializers.PaginationResponse[serializers.PermissionSerializer]
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
				permissions.EXPECT().List(gomock.Any(), gomock.Any()).Return([]models.Permission{
					{
						ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
						Name:        "read:self",
						Description: "Read own data",
					},
					{
						ID:          uuid.MustParse("10000000-1000-1000-3000-000000000002"),
						Name:        "write:self",
						Description: "Write own data",
					},
				}, uint64(2), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.PermissionSerializer]{
					Data: []serializers.PermissionSerializer{
						{
							ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
							Name:        "read:self",
							Description: "Read own data",
						},
						{
							ID:          uuid.MustParse("10000000-1000-1000-3000-000000000002"),
							Name:        "write:self",
							Description: "Write own data",
						},
					},
					Meta: serializers.PaginationMeta{
						Page:  1,
						Per:   25,
						Total: uint64(2),
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
				permissions.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.PermissionSerializer]{
					Data: []serializers.PermissionSerializer{},
					Meta: serializers.PaginationMeta{
						Page:  1,
						Per:   25,
						Total: uint64(0),
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
				permissions.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/permissions", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/permissions", controller.List)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PaginationResponse[serializers.PermissionSerializer]
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Permissions_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissions := services.NewMockPermissions(ctrl)
	controller := NewPermissionsController(permissions)

	type result struct {
		response serializers.PermissionSerializer
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
				permissions.EXPECT().FindById(gomock.Any(), uuid.MustParse("10000000-1000-1000-3000-000000000001")).Return(&models.Permission{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				}, nil)
			},
			expected: result{
				response: serializers.PermissionSerializer{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Not found",
			before: func() {
				permissions.EXPECT().FindById(gomock.Any(), uuid.MustParse("10000000-1000-1000-3000-000000000001")).Return(&models.Permission{}, errors.ErrPermissionNotFound)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrPermissionNotFound.Error()},
				status: "404 Not Found",
				code:   http.StatusNotFound,
			},
		},
		{
			name: "Error",
			before: func() {
				permissions.EXPECT().FindById(gomock.Any(), uuid.MustParse("10000000-1000-1000-3000-000000000001")).Return(&models.Permission{}, assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/permissions/10000000-1000-1000-3000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/permissions/{id}", controller.Get)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PermissionSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Permissions_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissions := services.NewMockPermissions(ctrl)
	controller := NewPermissionsController(permissions)

	type result struct {
		response serializers.PermissionSerializer
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
				permissions.EXPECT().Create(gomock.Any(), &models.Permission{
					Name:        "read:self",
					Description: "Read own data",
				}).Return(&models.Permission{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				}, nil)
			},
			body: strings.NewReader(`{"name": "read:self", "description": "Read own data"}`),
			expected: result{
				response: serializers.PermissionSerializer{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				},
				status: "201 Created",
				code:   http.StatusCreated,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				permissions.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"name": "read:self"}`),
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
				permissions.EXPECT().Create(gomock.Any(), &models.Permission{
					Name:        "read:self",
					Description: "Read own data",
				}).Return(&models.Permission{}, assert.AnError)
			},
			body: strings.NewReader(`{"name": "read:self", "description": "Read own data"}`),
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

			req := httptest.NewRequest(http.MethodPost, "/api/backoffice/permissions", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/api/backoffice/permissions", controller.Create)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PermissionSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Permissions_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissions := services.NewMockPermissions(ctrl)
	controller := NewPermissionsController(permissions)

	type result struct {
		response serializers.PermissionSerializer
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
				permissions.EXPECT().Update(gomock.Any(), &models.Permission{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				}).Return(&models.Permission{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				}, nil)
			},
			body: strings.NewReader(`{"name": "read:self", "description": "Read own data"}`),
			expected: result{
				response: serializers.PermissionSerializer{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				permissions.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"name": "read:self"}`),
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
				permissions.EXPECT().Update(gomock.Any(), &models.Permission{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				}).Return(&models.Permission{}, assert.AnError)
			},
			body: strings.NewReader(`{"name": "read:self", "description": "Read own data"}`),
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

			req := httptest.NewRequest(http.MethodPut, "/api/backoffice/permissions/10000000-1000-1000-3000-000000000001", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Put("/api/backoffice/permissions/{id}", controller.Update)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PermissionSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Permissions_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissions := services.NewMockPermissions(ctrl)
	controller := NewPermissionsController(permissions)

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
				permissions.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-3000-000000000001")).Return(true, nil)
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
				permissions.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-3000-000000000001")).Return(false, assert.AnError)
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

			req := httptest.NewRequest(http.MethodDelete, "/api/backoffice/permissions/10000000-1000-1000-3000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Delete("/api/backoffice/permissions/{id}", controller.Delete)
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
