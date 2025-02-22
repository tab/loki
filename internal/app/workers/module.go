package workers

import (
	"context"

	"go.uber.org/fx"
)

var Ctx context.Context

const (
	Success = "SUCCESS"
	Error   = "ERROR"

	TraceName         = "authentication"
	SmartIdWorkerName = "SmartId::Worker"
)

var Module = fx.Options(
	fx.Provide(NewSmartIdWorker),
)
