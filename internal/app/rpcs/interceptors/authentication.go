package interceptors

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"loki/internal/app/services"
	"loki/internal/config/logger"
	"loki/internal/config/middlewares"
	"loki/pkg/jwt"
	"loki/pkg/rbac"
)

const bearerScheme = "Bearer"

type AuthenticationInterceptor interface {
	Authenticate(ctx context.Context) (context.Context, error)
}

type authenticationInterceptor struct {
	jwt   jwt.Jwt
	users services.Users
	log   *logger.Logger
}

func NewAuthenticationInterceptor(jwt jwt.Jwt, users services.Users, log *logger.Logger) AuthenticationInterceptor {
	return &authenticationInterceptor{
		jwt:   jwt,
		users: users,
		log:   log,
	}
}

func (i *authenticationInterceptor) Authenticate(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, bearerScheme)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	claims, err := i.jwt.Decode(token)
	if err != nil {
		i.log.Error().Err(err).Msg("Failed to decode token")
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	user, err := i.users.FindByIdentityNumber(ctx, claims.ID)
	if err != nil {
		i.log.Error().Err(err).Msg("Failed to find user by identity number")
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	if !rbac.HasScope(claims.Scope) {
		i.log.Error().Msgf("User %s does not have required scope: %s", claims.ID, rbac.SsoServiceType)
		return nil, status.Errorf(codes.PermissionDenied, "missing required scope")
	}

	ctx = middlewares.NewContextModifier(ctx).
		WithClaim(claims).
		WithCurrentUser(user).
		Context()

	return ctx, nil
}
