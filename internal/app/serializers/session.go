package serializers

import "github.com/google/uuid"

type SessionSerializer struct {
	ID     uuid.UUID `json:"id"`
	Code   string    `json:"code,omitempty"`
	Status string    `json:"status,omitempty"`
	Error  string    `json:"error,omitempty"`
}
