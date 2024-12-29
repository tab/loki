package models

import "github.com/google/uuid"

const (
	ReadSelfType  = "read:self"
	WriteSelfType = "write:self"
)

type Permission struct {
	ID          uuid.UUID
	Name        string
	Description string
}
