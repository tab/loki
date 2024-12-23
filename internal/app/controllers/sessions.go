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

	session, err := c.sessions.FindById(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.ErrSessionNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.SessionSerializer{
		ID:     session.ID,
		Code:   session.Code,
		Status: session.Status,
		Error:  session.Error,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (c *sessionsController) Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	user, err := c.authentication.Complete(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.UserSerializer{
		ID:             user.ID,
		IdentityNumber: user.IdentityNumber,
		PersonalCode:   user.PersonalCode,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AccessToken:    user.AccessToken,
		RefreshToken:   user.RefreshToken,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
