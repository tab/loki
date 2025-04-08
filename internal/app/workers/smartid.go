package workers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tab/smartid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"loki/internal/app/models"
	"loki/internal/app/services"
	"loki/internal/config/logger"
)

type SmartIdWorker interface {
	Perform(ctx context.Context, id uuid.UUID, traceId string) *models.Session
}

type smartIdWorker struct {
	sessions services.Sessions
	users    services.Users
	worker   smartid.Worker
	log      *logger.Logger
}

func NewSmartIdWorker(
	sessions services.Sessions,
	users services.Users,
	worker smartid.Worker,
	log *logger.Logger,
) SmartIdWorker {
	return &smartIdWorker{
		sessions: sessions,
		users:    users,
		worker:   worker,
		log:      log,
	}
}

func (w *smartIdWorker) Perform(ctx context.Context, sessionId uuid.UUID, traceId string) *models.Session {
	w.log.Info().Msgf("%s perform %s", SmartIdWorkerName, sessionId)
	w.trace(ctx, traceId)

	resultCh := w.worker.Process(ctx, sessionId.String())

	result := <-resultCh
	if result.Err != nil {
		w.log.Error().Err(result.Err).Msgf("%s failed to get session status", SmartIdWorkerName)
		return w.updateSession(ctx, &models.UpdateSessionParams{
			ID:     sessionId,
			Status: Error,
			Error:  result.Err.Error(),
		})
	}

	user, err := w.users.Create(ctx, &models.User{
		IdentityNumber: result.Person.IdentityNumber,
		PersonalCode:   result.Person.PersonalCode,
		FirstName:      result.Person.FirstName,
		LastName:       result.Person.LastName,
	})

	if err != nil {
		w.log.Error().Err(err).Msgf("%s failed to create user", SmartIdWorkerName)
		return w.updateSession(ctx, &models.UpdateSessionParams{
			ID:     sessionId,
			Status: Error,
			Error:  err.Error(),
		})
	}

	return w.updateSession(ctx, &models.UpdateSessionParams{
		ID:     sessionId,
		UserId: user.ID,
		Status: Success,
	})
}

func (w *smartIdWorker) trace(ctx context.Context, traceId string) {
	tracer := otel.Tracer(TraceName)
	id, _ := trace.TraceIDFromHex(traceId)

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    id,
		Remote:     true,
		TraceFlags: trace.FlagsSampled,
	})

	operationName := fmt.Sprintf("%s perform {id}", SmartIdWorkerName)

	_, span := tracer.Start(
		trace.ContextWithSpanContext(ctx, spanCtx),
		operationName)
	defer span.End()
}

func (w *smartIdWorker) updateSession(ctx context.Context, params *models.UpdateSessionParams) *models.Session {
	session, err := w.sessions.Update(ctx, &models.UpdateSessionParams{
		ID:     params.ID,
		UserId: params.UserId,
		Status: params.Status,
		Error:  params.Error,
	})
	if err != nil {
		w.log.Error().Err(err).Msgf("%s failed to update session", SmartIdWorkerName)
	}

	return session
}
