package dto

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki/internal/app/errors"
)

func Test_Validate_PermissionRequest(t *testing.T) {
	tests := []struct {
		name     string
		body     io.Reader
		expected error
	}{
		{
			name:     "Success",
			body:     strings.NewReader(`{"name": "users:read", "description": "Read users"}`),
			expected: nil,
		},
		{
			name:     "Empty name",
			body:     strings.NewReader(`{"name": "", "description": "Read users"}`),
			expected: errors.ErrEmptyName,
		},
		{
			name:     "Empty description",
			body:     strings.NewReader(`{"name": "users:read", "description": ""}`),
			expected: errors.ErrEmptyDescription,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params PermissionRequest
			err := params.Validate(tt.body)

			assert.Equal(t, tt.expected, err)
		})
	}
}
