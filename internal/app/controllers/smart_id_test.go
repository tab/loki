package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models/dto"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

func Test_SmartIdController_CreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := gomock.Any()
	authentication := services.NewMockAuthentication(ctrl)
	provider := services.NewMockSmartIdProvider(ctrl)
	controller := NewSmartIdController(authentication, provider)

	sessionId := uuid.MustParse("8fdb516d-1a82-43ba-b82d-be63df569b86")

	type result struct {
		response serializers.SessionSerializer
		error    serializers.ErrorSerializer
		status   string
		code     int
	}

	tests := []struct {
		name     string
		body     io.Reader
		before   func()
		expected result
	}{
		{
			name: "Success",
			body: strings.NewReader(`{"country": "EE", "personal_code": "30303039914"}`),
			before: func() {
				authentication.EXPECT().CreateSmartIdSession(ctx, dto.CreateSmartIdSessionRequest{
					Country:      "EE",
					PersonalCode: "30303039914",
				}).Return(&serializers.SessionSerializer{
					ID:   sessionId,
					Code: "1234",
				}, nil)
			},
			expected: result{
				response: serializers.SessionSerializer{
					ID:   sessionId,
					Code: "1234",
				},
				status: "201 Created",
				code:   http.StatusCreated,
			},
		},
		{
			name:   "Bad request",
			body:   strings.NewReader(`{"personal_code": "30303039914"}`),
			before: func() {},
			expected: result{
				error:  serializers.ErrorSerializer{Error: "empty country, should be 'EE', 'LV' or 'LT'"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name: "Unprocessable entity",
			body: strings.NewReader(`{"country": "EE", "personal_code": "30303039914"}`),
			before: func() {
				authentication.EXPECT().CreateSmartIdSession(ctx, dto.CreateSmartIdSessionRequest{
					Country:      "EE",
					PersonalCode: "30303039914",
				}).Return(nil, assert.AnError)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: assert.AnError.Error()},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			r := httptest.NewRequest(http.MethodPost, "/api/auth/smart_id", tt.body)
			w := httptest.NewRecorder()

			controller.CreateSession(w, r)

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
