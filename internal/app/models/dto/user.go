package dto

type CreateUserTokensParams struct {
	IdentityNumber string
	AccessToken    CreateTokenParams
	RefreshToken   CreateTokenParams
}
