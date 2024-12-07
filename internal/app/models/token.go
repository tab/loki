package models

import "time"

const (
	AccessTokenType  = "access_token"
	RefreshTokenType = "refresh_token"

	AccessTokenExp  = time.Minute * 30
	RefreshTokenExp = time.Hour * 24
)
