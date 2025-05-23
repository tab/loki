package errors

import "errors"

var (
	// ErrInvalidToken indicates that the provided token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrPrivateKeyNotFound indicates that the private key for signing JWT tokens could not be found
	ErrPrivateKeyNotFound = errors.New("failed to load JWT private key")

	// ErrPublicKeyNotFound indicates that the public key for verifying JWT tokens could not be found
	ErrPublicKeyNotFound = errors.New("failed to load JWT public key")

	// ErrInvalidSigningMethod indicates that an unsupported signing method was used
	ErrInvalidSigningMethod = errors.New("invalid signing method")

	// ErrEmptyCountry indicates that the country code is empty or invalid
	ErrEmptyCountry = errors.New("empty country, should be 'EE', 'LV' or 'LT'")

	// ErrEmptyIdentityNumber indicates that the identity number is empty or invalid
	ErrEmptyIdentityNumber = errors.New("empty identity number")

	// ErrEmptyPersonalCode indicates that the personal code is empty or invalid
	ErrEmptyPersonalCode = errors.New("empty personal code")

	// ErrEmptyFirstName indicates that the first name is empty or invalid
	ErrEmptyFirstName = errors.New("empty first name")

	// ErrEmptyLastName indicates that the last name is empty or invalid
	ErrEmptyLastName = errors.New("empty last name")

	// ErrEmptyPhoneNumber indicates that the phone number is empty or invalid
	ErrEmptyPhoneNumber = errors.New("empty phone number")

	// ErrEmptyLocale indicates that the locale is empty or invalid
	ErrEmptyLocale = errors.New("empty locale")

	// ErrEmptyName indicates that the name is empty or invalid
	ErrEmptyName = errors.New("empty name")

	// ErrEmptyDescription indicates that the description is empty or invalid
	ErrEmptyDescription = errors.New("empty description")

	// ErrInvalidIdentityNumber indicates that the provided identity number is invalid
	ErrInvalidIdentityNumber = errors.New("invalid identity number")

	// ErrInvalidCertificate indicates that the provided certificate is invalid
	ErrInvalidCertificate = errors.New("invalid certificate")

	// ErrSessionNotFound indicates that the requested session could not be found
	ErrSessionNotFound = errors.New("session not found")

	// ErrInvalidArguments indicates that the provided request arguments are invalid
	ErrInvalidArguments = errors.New("invalid arguments")

	// ErrFailedToFetchResults indicates that failed to fetch results
	ErrFailedToFetchResults = errors.New("failed to fetch results")

	// ErrRecordNotFound indicates that the requested record could not be found
	ErrRecordNotFound = errors.New("record not found")

	// ErrFailedToCreateRecord indicates that failed to create record
	ErrFailedToCreateRecord = errors.New("failed to create record")

	// ErrFailedToUpdateRecord indicates that failed to update record
	ErrFailedToUpdateRecord = errors.New("failed to update record")

	// ErrFailedToDeleteRecord indicates that failed to delete record
	ErrFailedToDeleteRecord = errors.New("failed to delete record")

	// ErrPermissionNotFound indicates that the requested permission could not be found
	ErrPermissionNotFound = errors.New("permission not found")

	// ErrRoleNotFound indicates that the requested role could not be found
	ErrRoleNotFound = errors.New("role not found")

	// ErrScopeNotFound indicates that the requested scope could not be found
	ErrScopeNotFound = errors.New("scope not found")

	// ErrUserNotFound indicates that the requested user could not be found
	ErrUserNotFound = errors.New("user not found")

	// ErrSmartIdProviderError indicates an error originating from the Smart-ID provider
	ErrSmartIdProviderError = errors.New("smart-id provider error")

	// ErrMobileIdProviderError indicates an error originating from the Mobile-ID provider
	ErrMobileIdProviderError = errors.New("mobile-id provider error")

	// ErrForbidden indicates that the user is not allowed to perform the requested action
	ErrForbidden = errors.New("access forbidden")

	// ErrUnauthorized indicates that the user is not authorized to perform the requested action
	ErrUnauthorized = errors.New("unauthorized")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
