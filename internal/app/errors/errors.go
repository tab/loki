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

	ErrSessionNotFound = errors.New("session not found")
	ErrUserNotFound    = errors.New("user not found")

	ErrSmartIdProviderError  = errors.New("smart-id provider error")
	ErrMobileIdProviderError = errors.New("mobile-id provider error")

	ErrUnauthorized = errors.New("unauthorized")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
