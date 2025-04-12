package middlewares

import (
	"encoding/json"
	"net/http"

	"loki/internal/app/errors"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
	"loki/internal/config/logger"
	"loki/pkg/jwt"
	"loki/pkg/rbac"
)

type AuthorizationMiddleware interface {
	Authorize(next http.Handler) http.Handler
	Check(permission string) func(http.Handler) http.Handler
}

type authorizationMiddleware struct {
	jwt   jwt.Jwt
	users services.Users
	log   *logger.Logger
}

func NewAuthorizationMiddleware(jwt jwt.Jwt, users services.Users, log *logger.Logger) AuthorizationMiddleware {
	return &authorizationMiddleware{
		jwt:   jwt,
		users: users,
		log:   log,
	}
}

func (m *authorizationMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := extractBearerToken(r)
		if !ok {
			m.log.Error().Msg("Invalid authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claim, err := m.jwt.Decode(token)
		if err != nil {
			m.log.Error().Err(err).Msg("Failed to decode token")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
			return
		}

		user, err := m.users.FindByIdentityNumber(r.Context(), claim.ID)
		if err != nil {
			m.log.Error().Err(err).Msg("Failed to find user by identity number")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
			return
		}

		if !rbac.HasScope(claim.Scope) {
			m.log.Warn().Msgf("User %s does not have %s scope", claim.ID, rbac.SsoServiceType)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := NewContextModifier(r.Context()).
			WithCurrentUser(user).
			WithClaim(claim).
			Context()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *authorizationMiddleware) Check(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claim, ok := CurrentClaimFromContext(r.Context())
			if !ok {
				m.log.Error().Msg("No claims found in context")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrUnauthorized.Error()})
				return
			}

			if !rbac.HasPermission(claim.Permissions, permission) {
				m.log.Warn().Msgf("User %s does not have required permission: %s", claim.ID, permission)
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrForbidden.Error()})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
