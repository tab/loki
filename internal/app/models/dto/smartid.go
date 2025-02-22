package dto

import (
	"encoding/json"
	"io"
	"strings"

	"loki/internal/app/errors"
)

type CreateSmartIdSessionRequest struct {
	Country      string `json:"country"`
	PersonalCode string `json:"personal_code"`
}

func (params *CreateSmartIdSessionRequest) Validate(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(params); err != nil {
		return err
	}

	params.PersonalCode = strings.TrimSpace(params.PersonalCode)
	params.Country = strings.TrimSpace(params.Country)
	params.Country = strings.ToUpper(params.Country)

	if params.Country == "" {
		return errors.ErrEmptyCountry
	}

	if params.PersonalCode == "" {
		return errors.ErrEmptyPersonalCode
	}

	return nil
}
