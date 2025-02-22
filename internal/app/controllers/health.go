package controllers

import (
	"encoding/json"
	"net/http"

	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

type HealthController interface {
	HandleLiveness(w http.ResponseWriter, r *http.Request)
	HandleReadiness(w http.ResponseWriter, r *http.Request)
}

type healthController struct {
	service services.HealthChecker
}

func NewHealthController(service services.HealthChecker) HealthController {
	return &healthController{service: service}
}

// HandleLiveness handles application liveness check
func (h *healthController) HandleLiveness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(serializers.HealthSerializer{Result: "alive"})
}

// HandleReadiness handles application readiness check
func (h *healthController) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := h.service.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: "unavailable"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(serializers.HealthSerializer{Result: "ready"})
}
