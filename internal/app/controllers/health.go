package controllers

import (
	"encoding/json"
	"net/http"

	"loki/internal/app/serializers"
)

type HealthController interface {
	HandleLiveness(w http.ResponseWriter, r *http.Request)
}

type healthController struct {
}

func NewHealthController() HealthController {
	return &healthController{}
}

// HandleLiveness handles application liveness check
func (h *healthController) HandleLiveness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serializers.HealthSerializer{Result: "alive"})
}
