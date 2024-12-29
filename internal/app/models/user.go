package models

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID
	IdentityNumber string
	PersonalCode   string
	FirstName      string
	LastName       string
	AccessToken    string
	RefreshToken   string
}
