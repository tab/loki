package middlewares

import (
	"context"

	"loki/internal/app/models"
	"loki/pkg/jwt"
)

type Claim struct{}
type Token struct{}
type TraceId struct{}
type CurrentUser struct{}

type Modifier interface {
	WithClaim(claims *jwt.Payload) Modifier
	WithToken(token string) Modifier
	WithTraceId(traceId string) Modifier
	WithCurrentUser(user *models.User) Modifier
	Context() context.Context
}

type modifier struct {
	ctx context.Context
}

func NewContextModifier(ctx context.Context) Modifier {
	return &modifier{ctx: ctx}
}

func (m *modifier) WithClaim(claims *jwt.Payload) Modifier {
	m.ctx = context.WithValue(m.ctx, Claim{}, claims)
	return m
}

func (m *modifier) WithToken(token string) Modifier {
	m.ctx = context.WithValue(m.ctx, Token{}, token)
	return m
}

func (m *modifier) WithTraceId(traceId string) Modifier {
	m.ctx = context.WithValue(m.ctx, TraceId{}, traceId)
	return m
}

func (m *modifier) WithCurrentUser(user *models.User) Modifier {
	m.ctx = context.WithValue(m.ctx, CurrentUser{}, user)
	return m
}

func (m *modifier) Context() context.Context {
	return m.ctx
}
