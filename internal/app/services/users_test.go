package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/models"
	"loki/internal/app/repositories"
	"loki/internal/app/serializers"
	"loki/pkg/logger"
)

func Test_Users_FindByIdentityNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	database := repositories.NewMockDatabase(ctrl)
	log := logger.NewLogger()
	service := NewUsers(database, log)

	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	identityNumber := "PNOEE-123456789"

	tests := []struct {
		name     string
		before   func()
		expected *serializers.UserSerializer
		error    error
	}{
		{
			name: "Success",
			before: func() {
				database.EXPECT().FindUserByIdentityNumber(ctx, identityNumber).Return(&models.User{
					ID:             id,
					IdentityNumber: identityNumber,
					PersonalCode:   "123456789",
					FirstName:      "John",
					LastName:       "Doe",
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
				}, nil)
			},
			expected: &serializers.UserSerializer{
				ID:             id,
				IdentityNumber: identityNumber,
				PersonalCode:   "123456789",
				FirstName:      "John",
				LastName:       "Doe",
			},
		},
		{
			name: "Error",
			before: func() {
				database.EXPECT().FindUserByIdentityNumber(ctx, identityNumber).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindByIdentityNumber(ctx, identityNumber)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
