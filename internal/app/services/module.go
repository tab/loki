package services

import (
	"github.com/google/uuid"
	"go.uber.org/fx"
)

const QueueSize = 50

type SmartIdQueue struct {
	ID      uuid.UUID
	TraceId string
}

type MobileIdQueue struct {
	ID      uuid.UUID
	TraceId string
}

var Module = fx.Options(
	fx.Provide(
		func() chan *SmartIdQueue {
			return make(chan *SmartIdQueue, QueueSize)
		},
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
	fx.Provide(NewSmartId),
	fx.Provide(NewSmartIdWorker),
	fx.Provide(NewMobileId),
	fx.Provide(NewMobileIdWorker),
)
