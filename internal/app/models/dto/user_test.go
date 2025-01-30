package dto

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki/internal/app/errors"
)

func Test_Validate_UserRequest(t *testing.T) {
	tests := []struct {
		name     string
		body     io.Reader
		expected error
	}{
		{
			name:     "Success",
			body:     strings.NewReader(`{"identity_number": "PNOEE-123123123", "personal_code": "123123123", "first_name": "John", "last_name": "Doe"}`),
			expected: nil,
		},
		{
			name:     "Empty Identity Number",
			body:     strings.NewReader(`{"identity_number": "", "personal_code": "123123123", "first_name": "John", "last_name": "Doe"}`),
			expected: errors.ErrEmptyIdentityNumber,
		},
		{
			name:     "Empty Personal Code",
			body:     strings.NewReader(`{"identity_number": "PNOEE-123123123", "personal_code": "", "first_name": "John", "last_name": "Doe"}`),
			expected: errors.ErrEmptyPersonalCode,
		},
		{
			name:     "Empty First Name",
			body:     strings.NewReader(`{"identity_number": "PNOEE-123123123", "personal_code": "123123123", "first_name": "", "last_name": "Doe"}`),
			expected: errors.ErrEmptyFirstName,
		},
		{
			name:     "Empty Last Name",
			body:     strings.NewReader(`{"identity_number": "PNOEE-123123123", "personal_code": "123123123", "first_name": "John", "last_name": ""}`),
			expected: errors.ErrEmptyLastName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params UserRequest
			err := params.Validate(tt.body)

			assert.Equal(t, tt.expected, err)
		})
	}
}
