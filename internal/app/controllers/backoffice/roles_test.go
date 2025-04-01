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

func Test_Backoffice_Roles_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesService := services.NewMockRoles(ctrl)
	controller := NewRolesController(rolesService)

	type result struct {
		response serializers.PaginationResponse[serializers.RoleSerializer]
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
				rolesService.EXPECT().List(gomock.Any(), gomock.Any()).Return([]models.Role{
					{
						ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
						Name:        models.AdminRoleType,
						Description: "Admin role",
					},
					{
						ID:          uuid.MustParse("10000000-1000-1000-1000-000000000002"),
						Name:        models.ManagerRoleType,
						Description: "Manager role",
					},
					{
						ID:          uuid.MustParse("10000000-1000-1000-1000-000000000003"),
						Name:        models.UserRoleType,
						Description: "User role",
					},
				}, uint64(3), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.RoleSerializer]{
					Data: []serializers.RoleSerializer{
						{
							ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
							Name:        models.AdminRoleType,
							Description: "Admin role",
						},
						{
							ID:          uuid.MustParse("10000000-1000-1000-1000-000000000002"),
							Name:        models.ManagerRoleType,
							Description: "Manager role",
						},
						{
							ID:          uuid.MustParse("10000000-1000-1000-1000-000000000003"),
							Name:        models.UserRoleType,
							Description: "User role",
						},
					},
					Meta: serializers.PaginationMeta{
						Page:  1,
						Per:   25,
						Total: uint64(3),
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
				rolesService.EXPECT().List(gomock.Any(), gomock.Any()).Return([]models.Role{}, uint64(0), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.RoleSerializer]{
					Data: []serializers.RoleSerializer{},
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
				rolesService.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/roles", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/roles", controller.List)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PaginationResponse[serializers.RoleSerializer]
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Roles_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesService := services.NewMockRoles(ctrl)
	controller := NewRolesController(rolesService)

	type result struct {
		response serializers.RoleSerializer
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
				rolesService.EXPECT().FindRoleDetailsById(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(&models.Role{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
					PermissionIDs: []uuid.UUID{
						uuid.MustParse("10000000-1000-1000-3000-000000000001"),
						uuid.MustParse("10000000-1000-1000-3000-000000000002"),
					},
				}, nil)
			},
			expected: result{
				response: serializers.RoleSerializer{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
					PermissionIDs: []uuid.UUID{
						uuid.MustParse("10000000-1000-1000-3000-000000000001"),
						uuid.MustParse("10000000-1000-1000-3000-000000000002"),
					},
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Not found",
			before: func() {
				rolesService.EXPECT().FindRoleDetailsById(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(&models.Role{}, errors.ErrRoleNotFound)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrRoleNotFound.Error()},
				status: "404 Not Found",
				code:   http.StatusNotFound,
			},
			error: true,
		},
		{
			name: "Error",
			before: func() {
				rolesService.EXPECT().FindRoleDetailsById(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(&models.Role{}, assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/roles/10000000-1000-1000-1000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/roles/{id}", controller.Get)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.RoleSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Roles_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesService := services.NewMockRoles(ctrl)
	controller := NewRolesController(rolesService)

	type result struct {
		response serializers.RoleSerializer
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
				rolesService.EXPECT().Create(gomock.Any(), &models.Role{
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}).Return(&models.Role{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}, nil)
			},
			body: strings.NewReader(`{"name":"admin","description":"Admin role"}`),
			expected: result{
				response: serializers.RoleSerializer{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Admin role",
				},
				status: "201 Created",
				code:   http.StatusCreated,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				rolesService.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"name":"admin"}`),
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
				rolesService.EXPECT().Create(gomock.Any(), &models.Role{
					Name:        models.AdminRoleType,
					Description: "Admin role",
				}).Return(&models.Role{}, assert.AnError)
			},
			body: strings.NewReader(`{"name":"admin","description":"Admin role"}`),
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

			req := httptest.NewRequest(http.MethodPost, "/api/backoffice/roles", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/api/backoffice/roles", controller.Create)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.RoleSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Roles_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesService := services.NewMockRoles(ctrl)
	controller := NewRolesController(rolesService)

	type result struct {
		response serializers.RoleSerializer
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
				rolesService.EXPECT().Update(gomock.Any(), &models.Role{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Administrator role updated",
					PermissionIDs: []uuid.UUID{
						uuid.MustParse("10000000-1000-1000-3000-000000000001"),
						uuid.MustParse("10000000-1000-1000-3000-000000000002"),
					},
				}).Return(&models.Role{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Administrator role updated",
					PermissionIDs: []uuid.UUID{
						uuid.MustParse("10000000-1000-1000-3000-000000000001"),
						uuid.MustParse("10000000-1000-1000-3000-000000000002"),
					},
				}, nil)
			},
			body: strings.NewReader(`{"name":"admin","description":"Administrator role updated","permission_ids":["10000000-1000-1000-3000-000000000001","10000000-1000-1000-3000-000000000002"]}`),
			expected: result{
				response: serializers.RoleSerializer{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Administrator role updated",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				rolesService.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"name":"admin"}`),
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
				rolesService.EXPECT().Update(gomock.Any(), &models.Role{
					ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					Name:        models.AdminRoleType,
					Description: "Administrator role updated",
					PermissionIDs: []uuid.UUID{
						uuid.MustParse("10000000-1000-1000-3000-000000000001"),
						uuid.MustParse("10000000-1000-1000-3000-000000000002"),
					},
				}).Return(&models.Role{}, assert.AnError)
			},
			body: strings.NewReader(`{"name":"admin","description":"Administrator role updated","permission_ids":["10000000-1000-1000-3000-000000000001","10000000-1000-1000-3000-000000000002"]}`),
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

			req := httptest.NewRequest(http.MethodPut, "/api/backoffice/roles/10000000-1000-1000-1000-000000000001", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Put("/api/backoffice/roles/{id}", controller.Update)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.RoleSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Roles_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesService := services.NewMockRoles(ctrl)
	controller := NewRolesController(rolesService)

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
				rolesService.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(true, nil)
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
				rolesService.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(false, assert.AnError)
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

			req := httptest.NewRequest(http.MethodDelete, "/api/backoffice/roles/10000000-1000-1000-1000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Delete("/api/backoffice/roles/{id}", controller.Delete)
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
