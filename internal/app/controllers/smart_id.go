package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"loki/internal/app/models/dto"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

type SmartIdController interface {
	CreateSession(w http.ResponseWriter, r *http.Request)
	GetSessionStatus(w http.ResponseWriter, r *http.Request)
	CompleteSession(w http.ResponseWriter, r *http.Request)
}

type smartId struct {
	authentication services.Authentication
	provider       services.SmartIdProvider
}

func NewSmartIdController(
	authentication services.Authentication,
	provider services.SmartIdProvider,
) SmartIdController {
	return &smartId{
		authentication: authentication,
		provider:       provider,
	}
}

func (c *smartId) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.CreateSmartIdSessionRequest
	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response, err := c.authentication.CreateSmartIdSession(r.Context(), dto.CreateSmartIdSessionRequest{
		Country:      params.Country,
		PersonalCode: params.PersonalCode,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *smartId) GetSessionStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	response, err := c.authentication.FindSessionByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *smartId) CompleteSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	response, err := c.authentication.Complete(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
