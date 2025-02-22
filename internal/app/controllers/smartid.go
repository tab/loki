package controllers

import (
	"encoding/json"
	"net/http"

	"loki/internal/app/models/dto"
	"loki/internal/app/serializers"
	"loki/internal/app/services/authentication"
)

type SmartIdController interface {
	CreateSession(w http.ResponseWriter, r *http.Request)
}

type smartIdController struct {
	provider authentication.SmartIdProvider
}

func NewSmartIdController(provider authentication.SmartIdProvider) SmartIdController {
	return &smartIdController{
		provider: provider,
	}
}

func (c *smartIdController) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.CreateSmartIdSessionRequest
	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	session, err := c.provider.CreateSession(r.Context(), dto.CreateSmartIdSessionRequest{
		Country:      params.Country,
		PersonalCode: params.PersonalCode,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.SessionSerializer{
		ID:   session.ID,
		Code: session.Code,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
