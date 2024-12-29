package errors

import "errors"

var (
	// ErrInvalidToken indicates that the provided token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidSigningMethod indicates that an unsupported signing method was used
	ErrInvalidSigningMethod = errors.New("invalid signing method")

	// ErrEmptyCountry indicates that the country code is empty or invalid
	ErrEmptyCountry = errors.New("empty country, should be 'EE', 'LV' or 'LT'")

	// ErrEmptyPersonalCode indicates that the personal code is empty or invalid
	ErrEmptyPersonalCode = errors.New("empty personal code")

	// ErrEmptyPhoneNumber indicates that the phone number is empty or invalid
	ErrEmptyPhoneNumber = errors.New("empty phone number")

	// ErrEmptyLocale indicates that the locale is empty or invalid
	ErrEmptyLocale = errors.New("empty locale")

	// ErrInvalidIdentityNumber indicates that the provided identity number is invalid
	ErrInvalidIdentityNumber = errors.New("invalid identity number")

	// ErrInvalidCertificate indicates that the provided certificate is invalid
	ErrInvalidCertificate = errors.New("invalid certificate")

	// ErrSessionNotFound indicates that the requested session could not be found
	ErrSessionNotFound = errors.New("session not found")

	// ErrUserNotFound indicates that the requested user could not be found
	ErrUserNotFound = errors.New("user not found")

	// ErrSmartIdProviderError indicates an error originating from the Smart-ID provider
	ErrSmartIdProviderError = errors.New("smart-id provider error")

	// ErrMobileIdProviderError indicates an error originating from the Mobile-ID provider
	ErrMobileIdProviderError = errors.New("mobile-id provider error")

	// ErrUnauthorized indicates that the user is not authorized to perform the requested action
	ErrUnauthorized = errors.New("unauthorized")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
