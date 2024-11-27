package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"loki/internal/config"

	"github.com/stretchr/testify/assert"
)

func Test_HealthCheck(t *testing.T) {
	cfg := &config.Config{
		AppEnv:  "test",
		AppAddr: "localhost:8080",
	}
	router := NewRouter(cfg)

	req := httptest.NewRequest(http.MethodHead, "/health", nil)
	w := httptest.NewRecorder()

	resp := w.Result()
	defer resp.Body.Close()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
