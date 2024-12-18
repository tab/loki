package services

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/config"
	"loki/pkg/logger"
)

const MobileIdWorkerName = "MobileId::Worker"

type MobileIdWorker interface {
	Start(ctx context.Context)
	Stop()
}

type mobileIdWorker struct {
	cfg            *config.Config
	authentication Authentication
	certificate    Certificate
	sessions       Sessions
	users          Users
	queue          <-chan *MobileIdQueue
	wg             sync.WaitGroup
	log            *logger.Logger
}

func NewMobileIdWorker(
	cfg *config.Config,
	authentication Authentication,
	certificate Certificate,
	sessions Sessions,
	users Users,
	queue chan *MobileIdQueue,
	log *logger.Logger,
) MobileIdWorker {
	return &mobileIdWorker{
		cfg:            cfg,
		authentication: authentication,
		certificate:    certificate,
		sessions:       sessions,
		users:          users,
		queue:          queue,
		log:            log,
	}
}

func (w *mobileIdWorker) Start(ctx context.Context) {
	w.wg.Add(1)
	go w.run(ctx)
}

func (w *mobileIdWorker) Stop() {
	w.wg.Wait()
}

func (w *mobileIdWorker) run(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msgf("%s context cancelled, exiting", MobileIdWorkerName)
			return
		case req, ok := <-w.queue:
			if !ok {
				w.log.Info().Msgf("%s queue channel closed, exiting", MobileIdWorkerName)
				return
			}
			w.trace(ctx, req)
			w.perform(ctx, req)
		}
	}
}

func (w *mobileIdWorker) trace(ctx context.Context, req *MobileIdQueue) {
	tracer := otel.Tracer(AuthenticationTraceName)
	traceId, _ := trace.TraceIDFromHex(req.TraceId)

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceId,
		Remote:     true,
		TraceFlags: trace.FlagsSampled,
	})

	operationName := fmt.Sprintf("%s perform %s", SmartIdWorkerName, req.ID)

	_, span := tracer.Start(
		trace.ContextWithSpanContext(ctx, spanCtx),
		operationName)
	defer span.End()
}

func (w *mobileIdWorker) perform(ctx context.Context, req *MobileIdQueue) {
	w.log.Info().Msgf("%s perform %s", MobileIdWorkerName, req.ID)

	for {
		response, err := w.authentication.GetMobileIdSessionStatus(ctx, req.ID)
		if err != nil {
			w.log.Error().Err(err).Msgf("%s failed to get session status", MobileIdWorkerName)
			return
		}
		if w.processSessionState(ctx, req, response) {
			return
		}
	}
}

func (w *mobileIdWorker) processSessionState(ctx context.Context, req *MobileIdQueue, response *dto.MobileIdProviderSessionStatusResponse) bool {
	switch response.State {
	case models.SessionComplete:
		return w.handleSessionComplete(ctx, req, response)
	case models.SessionRunning:
		w.log.Warn().Msgf("%s session is still running", MobileIdWorkerName)
		return false
	default:
		w.log.Error().Msgf("%s unknown session state: %s", MobileIdWorkerName, response.State)
		return true
	}
}

func (w *mobileIdWorker) handleSessionComplete(ctx context.Context, req *MobileIdQueue, response *dto.MobileIdProviderSessionStatusResponse) bool {
	if response.Result == models.SessionResultOK {
		w.log.Info().Msgf("%s session is completed with OK result", MobileIdWorkerName)

		user, err := w.handleCreateUser(ctx, response)
		if err != nil {
			w.log.Error().Err(err).Msgf("%s failed to create user", MobileIdWorkerName)
			return true
		}

		_, err = w.sessions.Update(ctx, models.Session{
			ID:     req.ID,
			UserId: user.ID,
			Status: AuthenticationSuccess,
		})
		if err != nil {
			w.log.Error().Err(err).Msgf("%s failed to update session", MobileIdWorkerName)
			return true
		}
	} else {
		w.log.Info().Msgf("%s session is completed with error", MobileIdWorkerName)

		if _, err := w.sessions.Update(ctx, models.Session{
			ID:     req.ID,
			Status: AuthenticationError,
			Error:  response.Result,
		}); err != nil {
			w.log.Error().Err(err).Msgf("%s failed to update session", MobileIdWorkerName)
		}
	}

	return true
}

func (w *mobileIdWorker) handleCreateUser(ctx context.Context, response *dto.MobileIdProviderSessionStatusResponse) (user *models.User, err error) {
	cert, err := w.certificate.Extract(response.Cert)
	if err != nil {
		w.log.Error().Err(err).Msgf("%s failed to extract user from certificate", MobileIdWorkerName)
		return nil, err
	}

	return w.users.Create(ctx, &models.User{
		IdentityNumber: cert.IdentityNumber,
		PersonalCode:   cert.PersonalCode,
		FirstName:      cert.FirstName,
		LastName:       cert.LastName,
	})
}
