package dto

import (
	"encoding/json"
	"io"
	"strings"

	"loki/internal/app/errors"
)

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (params *RefreshAccessTokenRequest) Validate(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(params); err != nil {
		return err
	}

	params.RefreshToken = strings.TrimSpace(params.RefreshToken)
	if params.RefreshToken == "" {
		return errors.ErrInvalidToken
	}

	return nil
}
