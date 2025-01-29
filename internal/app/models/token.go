package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	AccessTokenType  = "access_token"
	RefreshTokenType = "refresh_token"

	AccessTokenExp  = time.Minute * 30
	RefreshTokenExp = time.Hour * 24
)

type Token struct {
	ID        uuid.UUID
	UserId    uuid.UUID
	Type      string
	Value     string
	ExpiresAt time.Time
}
