package services

import (
	"github.com/google/uuid"
	"go.uber.org/fx"
)

const QueueSize = 50

type SmartIdQueue struct {
	ID uuid.UUID
}

var Module = fx.Options(
	fx.Provide(
		func() chan *SmartIdQueue {
			return make(chan *SmartIdQueue, QueueSize)
		},
	),
	fx.Provide(NewAuthentication),
	fx.Provide(NewSmartId),
	fx.Provide(NewSmartIdWorker),
)
