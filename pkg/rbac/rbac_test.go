package rbac

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RBAC_HasPermission(t *testing.T) {
	tests := []struct {
		name               string
		claimPermissions   []string
		requiredPermission string
		expected           bool
	}{
		{
			name:               "Success",
			claimPermissions:   []string{"read:users", "write:users"},
			requiredPermission: "read:users",
			expected:           true,
		},
		{
			name:               "Fail",
			claimPermissions:   []string{"read:users", "write:users"},
			requiredPermission: "read:tokens",
			expected:           false,
		},
		{
			name:               "Empty",
			claimPermissions:   []string{},
			requiredPermission: "read:users",
			expected:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasPermission(tt.claimPermissions, tt.requiredPermission)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_RBAC_HasScope(t *testing.T) {
	tests := []struct {
		name       string
		claimScope []string
		expected   bool
	}{
		{
			name:       "Success",
			claimScope: []string{"sso-service"},
			expected:   true,
		},
		{
			name:       "Fail",
			claimScope: []string{"read:users", "write:users"},
			expected:   false,
		},
		{
			name:       "Empty",
			claimScope: []string{},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasScope(tt.claimScope)
			assert.Equal(t, tt.expected, result)
		})
	}
}
