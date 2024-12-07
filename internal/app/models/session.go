package models

import (
	"github.com/google/uuid"

	"loki/internal/app/models/dto"
)

const (
	SESSION_RUNNING  = "RUNNING"
	SESSION_COMPLETE = "COMPLETE"

	SESSION_RESULT_OK                                              = "OK"
	SESSION_RESULT_USER_REFUSED                                    = "USER_REFUSED"
	SESSION_RESULT_USER_REFUSED_DISPLAYTEXTANDPIN                  = "USER_REFUSED_DISPLAYTEXTANDPIN"
	SESSION_RESULT_USER_REFUSED_VC_CHOICE                          = "USER_REFUSED_VC_CHOICE"
	SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE                = "USER_REFUSED_CONFIRMATIONMESSAGE"
	SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE = "USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE"
	SESSION_RESULT_USER_REFUSED_CERT_CHOICE                        = "USER_REFUSED_CERT_CHOICE"
	SESSION_RESULT_WRONG_VC                                        = "WRONG_VC"
	SESSION_RESULT_TIMEOUT                                         = "TIMEOUT"
)

type SessionPayload struct {
	State     string                         `json:"state"`
	Result    dto.SmartIdProviderResult      `json:"result"`
	Signature dto.SmartIdProviderSignature   `json:"signature"`
	Cert      dto.SmartIdProviderCertificate `json:"cert"`
}

type Session struct {
	ID           uuid.UUID
	PersonalCode string
	Code         string
	Status       string
	Payload      SessionPayload
}
