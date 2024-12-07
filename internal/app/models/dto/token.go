package dto

import (
	"time"
)

type CreateTokenParams struct {
	Type      string
	Value     string
	ExpiresAt time.Time
}
