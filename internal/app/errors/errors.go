package errors

import "errors"

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidSigningMethod = errors.New("invalid signing method")

	ErrEmptyCountry      = errors.New("empty country, should be 'EE', 'LV' or 'LT'")
	ErrEmptyPersonalCode = errors.New("empty personal code")
	ErrEmptyPhoneNumber  = errors.New("empty phone number")
	ErrEmptyLocale       = errors.New("empty locale")

	ErrInvalidIdentityNumber = errors.New("invalid identity number")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
