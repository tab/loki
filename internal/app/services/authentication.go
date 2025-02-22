package services

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/config"
	"loki/pkg/logger"
)

const (
	AuthenticationSuccess = "SUCCESS"
	AuthenticationError   = "ERROR"

	AuthenticationTraceName = "authentication"
)

type Authentication interface {
	CreateMobileIdSession(ctx context.Context, params dto.CreateMobileIdSessionRequest) (*models.Session, error)
	GetMobileIdSessionStatus(ctx context.Context, id uuid.UUID) (*dto.MobileIdProviderSessionStatusResponse, error)
	Complete(ctx context.Context, id string) (*models.User, error)
}

type authentication struct {
	cfg              *config.Config
	mobileIdProvider MobileIdProvider
	mobileIdQueue    chan<- *MobileIdQueue
	sessions         Sessions
	tokens           Tokens
	log              *logger.Logger
}

func NewAuthentication(
	cfg *config.Config,
	mobileIdProvider MobileIdProvider,
	mobileIdQueue chan *MobileIdQueue,
	sessions Sessions,
	tokens Tokens,
	log *logger.Logger,
) Authentication {
	return &authentication{
		cfg:              cfg,
		mobileIdProvider: mobileIdProvider,
		mobileIdQueue:    mobileIdQueue,
		sessions:         sessions,
		tokens:           tokens,
		log:              log,
	}
}

func (a *authentication) CreateMobileIdSession(ctx context.Context, params dto.CreateMobileIdSessionRequest) (*models.Session, error) {
	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	result, err := a.mobileIdProvider.CreateSession(ctx, params)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to initiate MobileId authentication")
		return nil, err
	}

	session, err := a.sessions.Create(ctx, &models.CreateSessionParams{
		SessionId: result.ID,
		Code:      result.Code,
	})
	if err != nil {
		return nil, err
	}

	a.mobileIdQueue <- &MobileIdQueue{
		ID:      session.ID,
		TraceId: traceId,
	}

	return &models.Session{
		ID:   session.ID,
		Code: session.Code,
	}, nil
}

func (a *authentication) GetMobileIdSessionStatus(_ context.Context, id uuid.UUID) (*dto.MobileIdProviderSessionStatusResponse, error) {
	result, err := a.mobileIdProvider.GetSessionStatus(id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to get MobileId session status")
		return nil, err
	}

	return result, nil
}

func (a *authentication) Complete(ctx context.Context, sessionId string) (*models.User, error) {
	session, err := a.sessions.FindById(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	user, err := a.tokens.Create(ctx, session.UserId)
	if err != nil {
		return nil, err
	}

	err = a.sessions.Delete(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
