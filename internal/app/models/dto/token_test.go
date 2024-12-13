package dto

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki/internal/app/errors"
)

func Test_Validate_RefreshAccessTokenParams(t *testing.T) {
	tests := []struct {
		name     string
		body     io.Reader
		expected error
	}{
		{
			name:     "Success",
			body:     strings.NewReader(`{"refresh_token": "aaa.bbb.ccc"}`),
			expected: nil,
		},
		{
			name:     "Empty refresh token",
			body:     strings.NewReader(`{"refresh_token": ""}`),
			expected: errors.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params RefreshAccessTokenRequest
			err := params.Validate(tt.body)

			assert.Equal(t, tt.expected, err)
		})
	}
}
