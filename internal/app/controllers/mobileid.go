package controllers

import (
	"encoding/json"
	"net/http"

	"loki/internal/app/models/dto"
	"loki/internal/app/serializers"
	"loki/internal/app/services/authentication"
)

type MobileIdController interface {
	CreateSession(w http.ResponseWriter, r *http.Request)
}

type mobileIdController struct {
	provider authentication.MobileIdProvider
}

func NewMobileIdController(
	provider authentication.MobileIdProvider,
) MobileIdController {
	return &mobileIdController{
		provider: provider,
	}
}

func (c *mobileIdController) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.CreateMobileIdSessionRequest
	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	session, err := c.provider.CreateSession(r.Context(), dto.CreateMobileIdSessionRequest{
		PersonalCode: params.PersonalCode,
		PhoneNumber:  params.PhoneNumber,
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
