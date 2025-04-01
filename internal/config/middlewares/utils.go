package middlewares

import (
	"context"

	"loki/internal/app/models"
	"loki/pkg/jwt"
)

const (
	Authorization = "Authorization"
	bearerScheme  = "Bearer "
)

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
