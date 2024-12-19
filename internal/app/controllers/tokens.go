package controllers

import (
	"encoding/json"
	"net/http"

	"loki/internal/app/models/dto"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

type TokensController interface {
	Refresh(w http.ResponseWriter, r *http.Request)
}

type tokensController struct {
	tokens services.Tokens
}

func NewTokensController(tokens services.Tokens) TokensController {
	return &tokensController{tokens: tokens}
}

func (c *tokensController) Refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.RefreshAccessTokenRequest
	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response, err := c.tokens.Refresh(r.Context(), params.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
