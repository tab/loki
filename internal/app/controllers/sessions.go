package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"loki/internal/app/errors"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

type SessionsController interface {
	GetStatus(w http.ResponseWriter, r *http.Request)
	Authenticate(w http.ResponseWriter, r *http.Request)
}

type sessionsController struct {
	authentication services.Authentication
	sessions       services.Sessions
}

func NewSessionsController(
	authentication services.Authentication,
	sessions services.Sessions,
) SessionsController {
	return &sessionsController{
		authentication: authentication,
		sessions:       sessions,
	}
}

func (c *sessionsController) GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	response, err := c.sessions.FindById(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.ErrSessionNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *sessionsController) Authenticate(w http.ResponseWriter, r *http.Request) {
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
