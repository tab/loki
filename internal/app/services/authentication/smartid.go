package authentication

import (
	"context"

	"github.com/tab/smartid"
	"go.opentelemetry.io/otel/trace"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/app/services"
	"loki/internal/app/workers"
	"loki/internal/config/logger"
)

type SmartIdProvider interface {
	CreateSession(ctx context.Context, params dto.CreateSmartIdSessionRequest) (*models.Session, error)
}

type smartIdProvider struct {
	client   smartid.Client
	sessions services.Sessions
	users    services.Users
	worker   workers.SmartIdWorker
	log      *logger.Logger
}

func NewSmartId(
	client smartid.Client,
	sessions services.Sessions,
	users services.Users,
	worker workers.SmartIdWorker,
	log *logger.Logger,
) SmartIdProvider {
	return &smartIdProvider{
		client:   client,
		sessions: sessions,
		users:    users,
		worker:   worker,
		log:      log,
	}
}

func (s *smartIdProvider) CreateSession(ctx context.Context, params dto.CreateSmartIdSessionRequest) (*models.Session, error) {
	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	identity := smartid.NewIdentity(smartid.TypePNO, params.Country, params.PersonalCode)
	result, err := s.client.CreateSession(ctx, identity)

	if err != nil {
		s.log.Error().Msgf("Failed to create Smart-ID session: %v", err)
		return nil, err
	}

	session, err := s.sessions.Create(ctx, &models.CreateSessionParams{
		SessionId: result.Id,
		Code:      result.Code,
	})
	if err != nil {
		s.log.Error().Msgf("Failed to save Smart-ID session: %v", err)
		return nil, err
	}

	go s.worker.Perform(workers.Ctx, session.ID, traceId)

	return &models.Session{
		ID:     session.ID,
		Code:   session.Code,
		Status: models.SessionRunning,
	}, nil
}
