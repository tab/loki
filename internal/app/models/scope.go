package models

import "github.com/google/uuid"

const (
	SelfServiceType = "self-service"
	SsoServiceType  = "sso-service"
)

type Scope struct {
	ID          uuid.UUID
	Name        string
	Description string
}
