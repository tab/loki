package serializers

import "github.com/google/uuid"

type RoleSerializer struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`

	PermissionIDs []uuid.UUID `json:"permission_ids,omitempty"`
}
