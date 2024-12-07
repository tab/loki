package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/controllers"
	"loki/internal/config"
)

func Test_HealthCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:  "test",
		AppAddr: "localhost:8080",
	}
	mockSmartIdController := controllers.NewMockSmartIdController(ctrl)
	router := NewRouter(cfg, mockSmartIdController)

	req := httptest.NewRequest(http.MethodHead, "/health", nil)
	w := httptest.NewRecorder()

	resp := w.Result()
	defer resp.Body.Close()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
