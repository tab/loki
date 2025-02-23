package dto

import (
	"encoding/json"
	"io"
	"strings"

	"loki/internal/app/errors"
)

type CreateMobileIdSessionRequest struct {
	PersonalCode string `json:"personal_code"`
	PhoneNumber  string `json:"phone_number"`
}

func (params *CreateMobileIdSessionRequest) Validate(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(params); err != nil {
		return err
	}

	params.PersonalCode = strings.TrimSpace(params.PersonalCode)
	if params.PersonalCode == "" {
		return errors.ErrEmptyPersonalCode
	}

	params.PhoneNumber = strings.TrimSpace(params.PhoneNumber)
	if params.PhoneNumber == "" {
		return errors.ErrEmptyPhoneNumber
	}

	return nil
}
