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
	Locale       string `json:"locale"`
}

type MobileIdProviderSessionResponse struct {
	ID   string `json:"sessionID"`
	Code string `json:"code"`
}

type MobileIdProviderSessionStatusResponse struct {
	State     string                    `json:"state"`
	Result    string                    `json:"result"`
	Signature MobileIdProviderSignature `json:"signature"`
	Cert      string                    `json:"cert"`
	Time      string                    `json:"time"`
	TraceId   string                    `json:"traceId"`
}

type MobileIdProviderSignature struct {
	Value     string `json:"value"`
	Algorithm string `json:"algorithm"`
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

	params.Locale = strings.TrimSpace(params.Locale)
	if params.Locale == "" {
		return errors.ErrEmptyLocale
	}

	return nil
}
