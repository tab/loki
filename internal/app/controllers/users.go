package controllers

import (
	"encoding/json"
	"net/http"

	"loki/internal/app/errors"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
	"loki/internal/config/middlewares"
)

type UsersController interface {
	Me(w http.ResponseWriter, r *http.Request)
}

type usersController struct {
	users services.Users
}

func NewUsersController(users services.Users) UsersController {
	return &usersController{
		users: users,
	}
}

func (c *usersController) Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, ok := middlewares.CurrentUserFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrUnauthorized.Error()})
		return
	}

	response := serializers.UserSerializer{
		ID:             user.ID,
		IdentityNumber: user.IdentityNumber,
		PersonalCode:   user.PersonalCode,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
