package models

import "github.com/google/uuid"

const (
	SESSION_RUNNING  = "RUNNING"
	SESSION_COMPLETE = "COMPLETE"
	SESSION_ERROR    = "ERROR"

	SESSION_RESULT_OK                                              = "OK"
	SESSION_RESULT_USER_REFUSED                                    = "USER_REFUSED"
	SESSION_RESULT_USER_REFUSED_DISPLAYTEXTANDPIN                  = "USER_REFUSED_DISPLAYTEXTANDPIN"
	SESSION_RESULT_USER_REFUSED_VC_CHOICE                          = "USER_REFUSED_VC_CHOICE"
	SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE                = "USER_REFUSED_CONFIRMATIONMESSAGE"
	SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE = "USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE"
	SESSION_RESULT_USER_REFUSED_CERT_CHOICE                        = "USER_REFUSED_CERT_CHOICE"
	SESSION_RESULT_WRONG_VC                                        = "WRONG_VC"
	SESSION_RESULT_TIMEOUT                                         = "TIMEOUT"

	SESSION_RESULT_NOT_MID_CLIENT          = "NOT_MID_CLIENT"
	SESSION_RESULT_USER_CANCELLED          = "USER_CANCELLED"
	SESSION_RESULT_SIGNATURE_HASH_MISMATCH = "SIGNATURE_HASH_MISMATCH"
	SESSION_RESULT_PHONE_ABSENT            = "PHONE_ABSENT"
	SESSION_RESULT_DELIVERY_ERROR          = "DELIVERY_ERROR"
	SESSION_RESULT_SIM_ERROR               = "SIM_ERROR"

	SESSION_RESULT_UNKNOWN = "UNKNOWN"
)

type SessionPayload struct {
	State     string `json:"state"`
	Result    string `json:"result"`
	Signature string `json:"signature"`
	Cert      string `json:"cert"`
}

type Session struct {
	ID           uuid.UUID
	PersonalCode string
	Code         string
	Status       string
	Error        string
	Payload      SessionPayload
}

type CreateSessionParams struct {
	SessionId    string
	PersonalCode string
	Code         string
}
