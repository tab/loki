package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"loki/internal/app/models"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

const bearerScheme = "Bearer "

type CurrentUser struct{}
type Claim struct{}

type AuthenticationMiddleware interface {
	Authenticate(next http.Handler) http.Handler
}

type authenticationMiddleware struct {
	jwt   jwt.Jwt
	users services.Users
	log   *logger.Logger
}

func NewAuthenticationMiddleware(jwt jwt.Jwt, users services.Users, log *logger.Logger) AuthenticationMiddleware {
	return &authenticationMiddleware{
		jwt:   jwt,
		users: users,
		log:   log,
	}
}

func (m *authenticationMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := extractBearerToken(r)
		if !ok {
			m.log.Error().Msg("Invalid authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := m.jwt.Decode(token)
		if err != nil {
			m.log.Error().Err(err).Msg("Failed to decode token")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
			return
		}

		user, err := m.users.FindByIdentityNumber(r.Context(), claims.ID)
		if err != nil {
			m.log.Error().Err(err).Msg("Failed to find user by identity number")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
			return
		}

		ctx := withCurrentUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CurrentUserFromContext(ctx context.Context) (*models.User, bool) {
	u := ctx.Value(CurrentUser{})
	if u == nil {
		return nil, false
	}

	user, ok := u.(*models.User)
	return user, ok
}

func CurrentClaimFromContext(ctx context.Context) (*jwt.Payload, bool) {
	c, ok := ctx.Value(Claim{}).(*jwt.Payload)
	return c, ok
}

func withCurrentUser(ctx context.Context, user *models.User) context.Context {
	ctx = context.WithValue(ctx, CurrentUser{}, user)
	return ctx
}

func extractBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}

	if len(authHeader) < len(bearerScheme) || !strings.EqualFold(authHeader[:len(bearerScheme)], bearerScheme) {
		return "", false
	}

	token := authHeader[len(bearerScheme):]
	if token == "" {
		return "", false
	}

	return token, true
}
