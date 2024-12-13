package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/controllers"
	"loki/internal/config"
	"loki/internal/config/middlewares"
)

func Test_HealthCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:  "test",
		AppAddr: "localhost:8080",
	}

	mockAuthMiddleware := middlewares.NewMockAuthMiddleware(ctrl)
	mockSmartIdController := controllers.NewMockSmartIdController(ctrl)
	mockMobileIdController := controllers.NewMockMobileIdController(ctrl)
	mockSessionsController := controllers.NewMockSessionsController(ctrl)
	mockTokensController := controllers.NewMockTokensController(ctrl)
	mockUsersController := controllers.NewMockUsersController(ctrl)

	mockAuthMiddleware.EXPECT().
		Authenticate(gomock.Any()).
		AnyTimes().
		DoAndReturn(func(next http.Handler) http.Handler {
			return next
		})

	router := NewRouter(
		cfg,
		mockAuthMiddleware,
		mockSmartIdController,
		mockMobileIdController,
		mockSessionsController,
		mockTokensController,
		mockUsersController,
	)

	req := httptest.NewRequest(http.MethodHead, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
