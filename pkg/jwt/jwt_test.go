package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki/internal/app/errors"
	"loki/internal/config"
)

func Test_NewJWT(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewJWT(cfg)

	assert.NotNil(t, service)
}

func Test_JWT_Generate(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewJWT(cfg)

	type result struct {
		header string
	}

	tests := []struct {
		name     string
		payload  Payload
		expected result
	}{
		{
			name: "Success",
			payload: Payload{
				ID: "PNOEE-30303039914",
			},
			expected: result{
				header: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
		},
		{
			name: "Empty id",
			payload: Payload{
				ID: "",
			},
			expected: result{
				header: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Generate(tt.payload, time.Minute*30)
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.Equal(t, tt.expected.header, token[:36])
		})
	}
}

func Test_JWT_Verify(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewJWT(cfg)

	tests := []struct {
		name     string
		payload  Payload
		expected bool
	}{
		{
			name: "Success",
			payload: Payload{
				ID: "PNOEE-30303039914",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Generate(tt.payload, time.Minute*30)
			assert.NoError(t, err)

			result, err := service.Verify(token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_JWT_Verify_Mocked(t *testing.T) {
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

func Test_JWT_Decode(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewJWT(cfg)

	tests := []struct {
		name     string
		payload  Payload
		expected *Payload
	}{
		{
			name: "Success",
			payload: Payload{
				ID: "PNOEE-30303039914",
			},
			expected: &Payload{
				ID: "PNOEE-30303039914",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Generate(tt.payload, time.Minute*30)
			assert.NoError(t, err)

			result, err := service.Decode(token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_JWT_Decode_Mocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockJwt(ctrl)

	tests := []struct {
		name     string
		token    string
		before   func(token string)
		valid    bool
		expected *Payload
	}{
		{
			name:  "Success",
			token: "eyAAA.BBB.CCC",
			before: func(token string) {
				service.EXPECT().Decode(token).Return(&Payload{
					ID: "PNOEE-30303039914",
				}, nil)
			},
			expected: &Payload{
				ID: "PNOEE-30303039914",
			},
		},
		{
			name:  "Invalid token",
			token: "eyAAA.BBB.CCC",
			before: func(token string) {
				service.EXPECT().Decode(token).Return(nil, assert.AnError)
			},
			expected: nil,
		},
		{
			name:  "Invalid signing method",
			token: "eyAAA.BBB.CCC",
			before: func(token string) {
				service.EXPECT().Decode(token).Return(nil, errors.ErrInvalidSigningMethod)
			},
			expected: nil,
		},
		{
			name:  "Invalid token",
			token: "invalid-token",
			before: func(token string) {
				service.EXPECT().Decode(token).Return(nil, errors.ErrInvalidToken)
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(tt.token)

			result, err := service.Decode(tt.token)

			if tt.expected != nil {
				assert.Equal(t, tt.expected, result)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
