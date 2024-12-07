package serializers

import "github.com/google/uuid"

type SessionSerializer struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}
