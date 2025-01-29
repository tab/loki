package serializers

import "github.com/google/uuid"

type ScopeSerializer struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
