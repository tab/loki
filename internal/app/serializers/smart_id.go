package serializers

import "github.com/google/uuid"

type SmartIdSessionSerializer struct {
	ID     uuid.UUID `json:"id"`
	Code   string    `json:"code,omitempty"`
	Status string    `json:"status,omitempty"`
}
