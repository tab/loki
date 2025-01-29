package dto

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/google/uuid"

	"loki/internal/app/errors"
)

type RoleRequest struct {
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids,omitempty"`
}

func (params *RoleRequest) Validate(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(params); err != nil {
		return err
	}

	params.Name = strings.TrimSpace(params.Name)
	if params.Name == "" {
		return errors.ErrEmptyName
	}

	params.Description = strings.TrimSpace(params.Description)
	if params.Description == "" {
		return errors.ErrEmptyDescription
	}

	return nil
}
