package serializers

import "github.com/google/uuid"

type UserSerializer struct {
	ID             uuid.UUID `json:"id"`
	IdentityNumber string    `json:"identity_number"`
	PersonalCode   string    `json:"personal_code"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	AccessToken    string    `json:"access_token,omitempty"`
	RefreshToken   string    `json:"refresh_token,omitempty"`
}
