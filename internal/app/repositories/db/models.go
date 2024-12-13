// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type TokenType string

const (
	TokenTypeAccessToken  TokenType = "access_token"
	TokenTypeRefreshToken TokenType = "refresh_token"
)

func (e *TokenType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TokenType(s)
	case string:
		*e = TokenType(s)
	default:
		return fmt.Errorf("unsupported scan type for TokenType: %T", src)
	}
	return nil
}

type NullTokenType struct {
	TokenType TokenType
	Valid     bool // Valid is true if TokenType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTokenType) Scan(value interface{}) error {
	if value == nil {
		ns.TokenType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TokenType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTokenType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TokenType), nil
}

type Token struct {
	ID        uuid.UUID
	UserId    uuid.UUID
	Type      TokenType
	Value     string
	ExpiresAt pgtype.Timestamp
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type User struct {
	ID             uuid.UUID
	IdentityNumber string
	PersonalCode   string
	FirstName      string
	LastName       string
	CreatedAt      pgtype.Timestamp
	UpdatedAt      pgtype.Timestamp
}
