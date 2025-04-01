package interceptors

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"loki/pkg/jwt"
	"loki/pkg/logger"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/services"
	"loki/internal/config/middlewares"
)

func Test_AuthenticationInterceptor_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWT := jwt.NewMockJwt(ctrl)
	mockUsers := services.NewMockUsers(ctrl)
	log := logger.NewLogger()

	interceptor := NewAuthenticationInterceptor(mockJWT, mockUsers, log)

	userId := uuid.New()
	token := "valid-token"
	identityNumber := "PNOEE-1234567890"

	type result struct {
		code   codes.Code
		userId uuid.UUID
		error  bool
	}

	tests := []struct {
		name     string
		ctx      func() context.Context
		before   func()
		expected result
	}{
		{
			name: "Success",
			ctx: func() context.Context {
				md := metadata.New(map[string]string{
					"authorization": "Bearer " + token,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			before: func() {
				mockJWT.EXPECT().Decode(token).Return(&jwt.Payload{
					ID:          identityNumber,
					Permissions: []string{"read:users"},
					Roles:       []string{"admin"},
					Scope:       []string{"sso-service"},
				}, nil)
				mockUsers.EXPECT().FindByIdentityNumber(gomock.Any(), identityNumber).Return(&models.User{
					ID:             userId,
					IdentityNumber: identityNumber,
					FirstName:      "Test",
					LastName:       "User",
				}, nil)
			},
			expected: result{
				code:   codes.OK,
				userId: userId,
				error:  false,
			},
		},
		{
			name:   "Missing auth header",
			ctx:    context.Background,
			before: func() {},
			expected: result{
				code:   codes.Unauthenticated,
				userId: uuid.Nil,
				error:  true,
			},
		},
		{
			name: "Invalid auth scheme",
			ctx: func() context.Context {
				md := metadata.New(map[string]string{
					"authorization": "Basic " + token,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			before: func() {},
			expected: result{
				code:   codes.Unauthenticated,
				userId: uuid.Nil,
				error:  true,
			},
		},
		{
			name: "JWT decode error",
			ctx: func() context.Context {
				md := metadata.New(map[string]string{
					"authorization": "Bearer " + token,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			before: func() {
				mockJWT.EXPECT().Decode(token).Return(nil, errors.ErrInvalidToken)
			},
			expected: result{
				code:   codes.Unauthenticated,
				userId: uuid.Nil,
				error:  true,
			},
		},
		{
			name: "User not found",
			ctx: func() context.Context {
				md := metadata.New(map[string]string{
					"authorization": "Bearer " + token,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			before: func() {
				mockJWT.EXPECT().Decode(token).Return(&jwt.Payload{
					ID:          identityNumber,
					Permissions: []string{"read:users"},
					Roles:       []string{"admin"},
					Scope:       []string{"sso-service"},
				}, nil)
				mockUsers.EXPECT().FindByIdentityNumber(gomock.Any(), identityNumber).Return(nil, errors.ErrUserNotFound)
			},
			expected: result{
				code:   codes.Unauthenticated,
				userId: uuid.Nil,
				error:  true,
			},
		},
		{
			name: "Missing required scope",
			ctx: func() context.Context {
				md := metadata.New(map[string]string{
					"authorization": "Bearer " + token,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			before: func() {
				mockJWT.EXPECT().Decode(token).Return(&jwt.Payload{
					ID:          identityNumber,
					Permissions: []string{"read:users"},
					Roles:       []string{"admin"},
					Scope:       []string{"not-sso-service-scope"},
				}, nil)
				mockUsers.EXPECT().FindByIdentityNumber(gomock.Any(), identityNumber).Return(&models.User{
					ID:             userId,
					IdentityNumber: identityNumber,
					FirstName:      "Test",
					LastName:       "User",
				}, nil)
			},
			expected: result{
				code:   codes.PermissionDenied,
				userId: uuid.Nil,
				error:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			resultCtx, err := interceptor.Authenticate(tt.ctx())

			if tt.expected.error {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expected.code, st.Code())
			} else {
				assert.NoError(t, err)

				user, ok := middlewares.CurrentUserFromContext(resultCtx)
				assert.True(t, ok)
				assert.Equal(t, tt.expected.userId, user.ID)

				claim, ok := middlewares.CurrentClaimFromContext(resultCtx)
				assert.True(t, ok)
				assert.Equal(t, identityNumber, claim.ID)
			}
		})
	}
}
