package dto

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/google/uuid"

	"loki/internal/app/errors"
)

type UserRequest struct {
	IdentityNumber string      `json:"identity_number"`
	PersonalCode   string      `json:"personal_code"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	RoleIDs        []uuid.UUID `json:"role_ids,omitempty"`
	ScopeIDs       []uuid.UUID `json:"scope_ids,omitempty"`
}

func (params *UserRequest) Validate(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(params); err != nil {
		return err
	}

	params.IdentityNumber = strings.TrimSpace(params.IdentityNumber)
	if params.IdentityNumber == "" {
		return errors.ErrEmptyIdentityNumber
	}

	params.PersonalCode = strings.TrimSpace(params.PersonalCode)
	if params.PersonalCode == "" {
		return errors.ErrEmptyPersonalCode
	}

	params.FirstName = strings.TrimSpace(params.FirstName)
	if params.FirstName == "" {
		return errors.ErrEmptyFirstName
	}

	params.LastName = strings.TrimSpace(params.LastName)
	if params.LastName == "" {
		return errors.ErrEmptyLastName
	}

	return nil
}
