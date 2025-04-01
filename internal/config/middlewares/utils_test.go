package middlewares

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"loki/internal/app/models"
	"loki/pkg/jwt"
)

func Test_CurrentUserFromContext(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		user   *models.User
		exists bool
	}{
		{
			name: "Success",
			ctx: context.WithValue(context.Background(), CurrentUser{}, &models.User{
				ID:             uuid.MustParse("10000000-0000-0000-0000-000000000000"),
				IdentityNumber: "PNOEE-123456789",
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			}),
			user: &models.User{
				ID:             uuid.MustParse("10000000-0000-0000-0000-000000000000"),
				IdentityNumber: "PNOEE-123456789",
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			},
			exists: true,
		},
		{
			name:   "User does not exist",
			ctx:    context.Background(),
			user:   nil,
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, exists := CurrentUserFromContext(tt.ctx)
			assert.Equal(t, tt.exists, exists)

			if tt.exists {
				assert.Equal(t, tt.user.ID, user.ID)
				assert.Equal(t, tt.user.IdentityNumber, user.IdentityNumber)
				assert.Equal(t, tt.user.PersonalCode, user.PersonalCode)
				assert.Equal(t, tt.user.FirstName, user.FirstName)
				assert.Equal(t, tt.user.LastName, user.LastName)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}

func Test_CurrentClaimFromContext(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		claim  *jwt.Payload
		exists bool
	}{
		{
			name: "Success",
			ctx: context.WithValue(context.Background(), Claim{}, &jwt.Payload{
				ID:          "test-user",
				Permissions: []string{"read:users"},
			}),
			claim: &jwt.Payload{
				ID:          "test-user",
				Permissions: []string{"read:users"},
			},
			exists: true,
		},
		{
			name:   "Claim does not exist",
			ctx:    context.Background(),
			claim:  nil,
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claim, exists := CurrentClaimFromContext(tt.ctx)
			assert.Equal(t, tt.exists, exists)

			if tt.exists {
				assert.Equal(t, tt.claim.ID, claim.ID)
				assert.Equal(t, tt.claim.Permissions, claim.Permissions)
			} else {
				assert.Nil(t, claim)
			}
		})
	}
}
