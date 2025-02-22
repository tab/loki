package dto

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki/internal/app/errors"
)

func Test_ValidateSmartIdParams(t *testing.T) {
	tests := []struct {
		name     string
		body     io.Reader
		expected error
	}{
		{
			name:     "Success",
			body:     strings.NewReader(`{"country": "EE", "personal_code": "30303039914"}`),
			expected: nil,
		},
		{
			name:     "Empty personal code",
			body:     strings.NewReader(`{"country": "EE"}`),
			expected: errors.ErrEmptyPersonalCode,
		},
		{
			name:     "Empty country",
			body:     strings.NewReader(`{"personal_code": "30303039914"}`),
			expected: errors.ErrEmptyCountry,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params CreateSmartIdSessionRequest
			err := params.Validate(tt.body)

			assert.Equal(t, tt.expected, err)
		})
	}
}
