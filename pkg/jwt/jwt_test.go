package jwt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/errors"
	"loki/internal/config"
)

func Test_NewJWTService(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewJWTService(cfg)

	assert.NotNil(t, service)
}

func Test_JWTService_Generate(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewJWTService(cfg)

	UUID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	type result struct {
		header string
	}

	tests := []struct {
		name     string
		id       uuid.UUID
		expected result
	}{
		{
			name: "Success",
			id:   UUID,
			expected: result{
				header: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
		},
		{
			name: "Empty id",
			id:   uuid.UUID{},
			expected: result{
				header: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
		},
		{
			name: "Nil id",
			id:   uuid.Nil,
			expected: result{
				header: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Generate(tt.id)
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.Equal(t, tt.expected.header, token[:36])
		})
	}
}

func Test_JWTService_Verify(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewJWTService(cfg)

	tests := []struct {
		name     string
		uuid     uuid.UUID
		expected bool
	}{
		{
			name:     "Success",
			uuid:     uuid.New(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Generate(tt.uuid)
			assert.NoError(t, err)

			result, err := service.Verify(token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_JWTService_Verify_Mocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockJwt(ctrl)

	tests := []struct {
		name     string
		token    string
		before   func(token string)
		valid    bool
		expected error
	}{
		{
			name:  "Success",
			token: "eyAAA.BBB.CCC",
			before: func(token string) {
				service.EXPECT().Verify(token).Return(true, nil)
			},
			valid:    true,
			expected: nil,
		},
		{
			name:  "Invalid token",
			token: "eyAAA.BBB.CCC",
			before: func(token string) {
				service.EXPECT().Verify(token).Return(false, nil)
			},
			valid:    false,
			expected: nil,
		},
		{
			name:  "Invalid signing method",
			token: "eyAAA.BBB.CCC",
			before: func(token string) {
				service.EXPECT().Verify(token).Return(false, errors.ErrInvalidSigningMethod)
			},
			valid:    false,
			expected: errors.ErrInvalidSigningMethod,
		},
		{
			name:  "Invalid token",
			token: "invalid-token",
			before: func(token string) {
				service.EXPECT().Verify(token).Return(false, errors.ErrInvalidToken)
			},
			valid:    false,
			expected: errors.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(tt.token)

			result, err := service.Verify(tt.token)

			if tt.expected != nil {
				assert.ErrorIs(t, err, tt.expected)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.valid, result)
		})
	}
}
