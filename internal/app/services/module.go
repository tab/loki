package services

import (
	"github.com/google/uuid"
	"go.uber.org/fx"
)

const QueueSize = 50

type MobileIdQueue struct {
	ID      uuid.UUID
	TraceId string
}

var Module = fx.Options(
	fx.Provide(
		func() chan *MobileIdQueue {
			return make(chan *MobileIdQueue, QueueSize)
		},
	),
	fx.Provide(NewAuthentication),
	fx.Provide(NewCertificate),
	fx.Provide(NewSessions),
	fx.Provide(NewPermissions),
	fx.Provide(NewRoles),
	fx.Provide(NewScopes),
	fx.Provide(NewTokens),
	fx.Provide(NewUsers),
	fx.Provide(NewMobileId),
	fx.Provide(NewMobileIdWorker),
)
