package models

import "github.com/google/uuid"

const (
	AdminRoleType   = "admin"
	ManagerRoleType = "manager"
	UserRoleType    = "user"
)

type Role struct {
	ID          uuid.UUID
	Name        string
	Description string

	PermissionIDs []uuid.UUID
}

type RolePermission struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}
