package dto

type CreateUserParams struct {
	IdentityNumber string
	PersonalCode   string
	FirstName      string
	LastName       string
	AccessToken    CreateTokenParams
	RefreshToken   CreateTokenParams
}
