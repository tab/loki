package models

import "github.com/google/uuid"

const (
	SessionRunning  = "RUNNING"
	SessionComplete = "COMPLETE"
	SessionResultOK = "OK"
)

type Session struct {
	ID     uuid.UUID
	UserId uuid.UUID
	Code   string
	Status string
	Error  string
}

type CreateSessionParams struct {
	SessionId string
	Code      string
}
