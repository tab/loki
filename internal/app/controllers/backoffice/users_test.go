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

func Test_Backoffice_Users_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users)

	type result struct {
		response serializers.PaginationResponse[serializers.UserSerializer]
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
				users.EXPECT().List(gomock.Any(), gomock.Any()).Return([]models.User{
					{
						ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
						IdentityNumber: "PNOEE-123456789",
						PersonalCode:   "123456789",
						FirstName:      "John",
						LastName:       "Doe",
					},
					{
						ID:             uuid.MustParse("10000000-1000-1000-1000-000000000002"),
						IdentityNumber: "PNOEE-987654321",
						PersonalCode:   "987654321",
						FirstName:      "Jane",
						LastName:       "Doe",
					},
				}, uint64(2), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.UserSerializer]{
					Data: []serializers.UserSerializer{
						{
							ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
							IdentityNumber: "PNOEE-123456789",
							PersonalCode:   "123456789",
							FirstName:      "John",
							LastName:       "Doe",
						},
						{
							ID:             uuid.MustParse("10000000-1000-1000-1000-000000000002"),
							IdentityNumber: "PNOEE-987654321",
							PersonalCode:   "987654321",
							FirstName:      "Jane",
							LastName:       "Doe",
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
				users.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.UserSerializer]{
					Data: []serializers.UserSerializer{},
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
				users.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/users", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/users", controller.List)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PaginationResponse[serializers.UserSerializer]
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Backoffice_Users_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users)

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
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				users.EXPECT().FindUserDetailsById(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(&models.User{
					ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}, nil)
			},
			expected: result{
				response: serializers.UserSerializer{
					ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Not found",
			before: func() {
				users.EXPECT().FindUserDetailsById(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(&models.User{}, errors.ErrUserNotFound)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrUserNotFound.Error()},
				status: "404 Not Found",
				code:   http.StatusNotFound,
			},
		},
		{
			name: "Error",
			before: func() {
				users.EXPECT().FindUserDetailsById(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(&models.User{}, assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/users/10000000-1000-1000-1000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/users/{id}", controller.Get)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
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

func Test_Backoffice_Users_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users)

	type result struct {
		response serializers.UserSerializer
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
				users.EXPECT().Create(gomock.Any(), &models.User{
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}).Return(&models.User{
					ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}, nil)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-123456789", "personal_code": "123456789", "first_name": "John", "last_name": "Doe"}`),
			expected: result{
				response: serializers.UserSerializer{
					ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				},
				status: "201 Created",
				code:   http.StatusCreated,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				users.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-123456789"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "empty personal code"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Error",
			before: func() {
				users.EXPECT().Create(gomock.Any(), &models.User{
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}).Return(&models.User{}, assert.AnError)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-123456789", "personal_code": "123456789", "first_name": "John", "last_name": "Doe"}`),
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

			req := httptest.NewRequest(http.MethodPost, "/api/backoffice/users", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/api/backoffice/users", controller.Create)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
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

func Test_Backoffice_Users_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users)

	type result struct {
		response serializers.UserSerializer
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
				users.EXPECT().Update(gomock.Any(), &models.User{
					ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "JOHN",
					LastName:       "DOE",
				}).Return(&models.User{
					ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "JOHN",
					LastName:       "DOE",
				}, nil)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-123456789", "personal_code": "123456789", "first_name": "JOHN", "last_name":"DOE"}`),
			expected: result{
				response: serializers.UserSerializer{
					ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "JOHN",
					LastName:       "DOE",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				users.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-123456789"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "empty personal code"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Error",
			before: func() {
				users.EXPECT().Update(gomock.Any(), &models.User{
					ID:             uuid.MustParse("10000000-1000-1000-1000-000000000001"),
					IdentityNumber: "PNOEE-123456789",
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
				}).Return(&models.User{}, assert.AnError)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-123456789", "personal_code": "123456789", "first_name": "John", "last_name": "Doe"}`),
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

			req := httptest.NewRequest(http.MethodPut, "/api/backoffice/users/10000000-1000-1000-1000-000000000001", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Put("/api/backoffice/users/{id}", controller.Update)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
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

func Test_Backoffice_Users_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users)

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
				users.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(true, nil)
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
				users.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(false, assert.AnError)
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

			req := httptest.NewRequest(http.MethodDelete, "/api/backoffice/users/10000000-1000-1000-1000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Delete("/api/backoffice/users/{id}", controller.Delete)
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
