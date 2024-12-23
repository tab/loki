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
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	user, err := c.tokens.Update(r.Context(), params.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.TokensSerializer{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
