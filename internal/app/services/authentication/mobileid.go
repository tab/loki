package authentication

import (
	"context"

	"github.com/tab/mobileid"
	"go.opentelemetry.io/otel/trace"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/app/services"
	"loki/internal/app/workers"
	"loki/pkg/logger"
)

type MobileIdProvider interface {
	CreateSession(ctx context.Context, params dto.CreateMobileIdSessionRequest) (*models.Session, error)
}

type mobileIdProvider struct {
	client   mobileid.Client
	sessions services.Sessions
	users    services.Users
	worker   workers.MobileIdWorker
	log      *logger.Logger
}

func NewMobileId(
	client mobileid.Client,
	sessions services.Sessions,
	users services.Users,
	worker workers.MobileIdWorker,
	log *logger.Logger,
) MobileIdProvider {
	return &mobileIdProvider{
		client:   client,
		sessions: sessions,
		users:    users,
		worker:   worker,
		log:      log,
	}
}

func (s *mobileIdProvider) CreateSession(ctx context.Context, params dto.CreateMobileIdSessionRequest) (*models.Session, error) {
	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	result, err := s.client.CreateSession(ctx, params.PhoneNumber, params.PersonalCode)

	if err != nil {
		s.log.Error().Msgf("Failed to create Mobile-ID session: %v", err)
		return nil, err
	}

	session, err := s.sessions.Create(ctx, &models.CreateSessionParams{
		SessionId: result.Id,
		Code:      result.Code,
	})
	if err != nil {
		s.log.Error().Msgf("Failed to save Mobile-ID session: %v", err)
		return nil, err
	}

	go s.worker.Perform(workers.Ctx, session.ID, traceId)

	return &models.Session{
		ID:     session.ID,
		Code:   session.Code,
		Status: models.SessionRunning,
	}, nil
}
