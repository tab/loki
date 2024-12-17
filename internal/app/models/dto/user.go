package dto

type CreateUserParams struct {
    IdentityNumber string
    AccessToken    CreateTokenParams
    RefreshToken   CreateTokenParams
}
