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

type SmartIdProviderSessionResponse struct {
	ID   string `json:"sessionID"`
	Code string `json:"code"`
}

type SmartIdProviderSessionStatusResponse struct {
	State               string                     `json:"state"`
	Result              SmartIdProviderResult      `json:"result"`
	Signature           SmartIdProviderSignature   `json:"signature"`
	Cert                SmartIdProviderCertificate `json:"cert"`
	InteractionFlowUsed string                     `json:"interactionFlowUsed"`
}

type SmartIdProviderResult struct {
	EndResult      string `json:"endResult"`
	DocumentNumber string `json:"documentNumber"`
}

type SmartIdProviderSignature struct {
	Value     string `json:"value"`
	Algorithm string `json:"algorithm"`
}

type SmartIdProviderCertificate struct {
	Value            string `json:"value"`
	CertificateLevel string `json:"certificateLevel"`
}

type ProviderCertificateExtract struct {
	IdentityNumber string `json:"identity_number"`
	PersonalCode   string `json:"personal_code"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
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
