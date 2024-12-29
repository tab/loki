package dto

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki/internal/app/errors"
)

func Test_ValidateMobileIdParams(t *testing.T) {
	tests := []struct {
		name     string
		body     io.Reader
		expected error
	}{
		{
			name:     "Success",
			body:     strings.NewReader(`{"locale": "ENG", "phone_number": "+37268000769", "personal_code": "60001017869"}`),
			expected: nil,
		},
		{
			name:     "Empty personal code",
			body:     strings.NewReader(`{"locale": "ENG", "phone_number": "+37268000769"}`),
			expected: errors.ErrEmptyPersonalCode,
		},
		{
			name:     "Empty phone number",
			body:     strings.NewReader(`{"locale": "ENG", "personal_code": "60001017869"}`),
			expected: errors.ErrEmptyPhoneNumber,
		},
		{
			name:     "Empty locale",
			body:     strings.NewReader(`{"phone_number": "+37268000769", "personal_code": "60001017869"}`),
			expected: errors.ErrEmptyLocale,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params CreateMobileIdSessionRequest
			err := params.Validate(tt.body)

			assert.Equal(t, tt.expected, err)
		})
	}
}
