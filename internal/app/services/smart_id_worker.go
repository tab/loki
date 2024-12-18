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

const SmartIdWorkerName = "SmartId::Worker"

type SmartIdWorker interface {
	Start(ctx context.Context)
	Stop()
}

type smartIdWorker struct {
	cfg            *config.Config
	authentication Authentication
	certificate    Certificate
	sessions       Sessions
	users          Users
	queue          <-chan *SmartIdQueue
	wg             sync.WaitGroup
	log            *logger.Logger
}

func NewSmartIdWorker(
	cfg *config.Config,
	authentication Authentication,
	certificate Certificate,
	sessions Sessions,
	users Users,
	queue chan *SmartIdQueue,
	log *logger.Logger,
) SmartIdWorker {
	return &smartIdWorker{
		cfg:            cfg,
		authentication: authentication,
		certificate:    certificate,
		sessions:       sessions,
		users:          users,
		queue:          queue,
		log:            log,
	}
}

func (w *smartIdWorker) Start(ctx context.Context) {
	w.wg.Add(1)
	go w.run(ctx)
}

func (w *smartIdWorker) Stop() {
	w.wg.Wait()
}

func (w *smartIdWorker) run(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msgf("%s context cancelled, exiting", SmartIdWorkerName)
			return
		case req, ok := <-w.queue:
			if !ok {
				w.log.Info().Msgf("%s queue channel closed, exiting", SmartIdWorkerName)
				return
			}
			w.trace(ctx, req)
			w.perform(ctx, req)
		}
	}
}

func (w *smartIdWorker) trace(ctx context.Context, req *SmartIdQueue) {
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

func (w *smartIdWorker) perform(ctx context.Context, req *SmartIdQueue) {
	w.log.Info().Msgf("%s perform %s", SmartIdWorkerName, req.ID)

	for {
		response, err := w.authentication.GetSmartIdSessionStatus(ctx, req.ID)
		if err != nil {
			w.log.Error().Err(err).Msgf("%s failed to get session status", SmartIdWorkerName)
			return
		}
		if w.processSessionState(ctx, req, response) {
			return
		}
	}
}

func (w *smartIdWorker) processSessionState(ctx context.Context, req *SmartIdQueue, response *dto.SmartIdProviderSessionStatusResponse) bool {
	switch response.State {
	case models.SessionComplete:
		return w.handleSessionComplete(ctx, req, response)
	case models.SessionRunning:
		w.log.Warn().Msgf("%s session is still running", SmartIdWorkerName)
		return false
	default:
		w.log.Error().Msgf("%s unknown session state: %s", SmartIdWorkerName, response.State)
		return true
	}
}

func (w *smartIdWorker) handleSessionComplete(ctx context.Context, req *SmartIdQueue, response *dto.SmartIdProviderSessionStatusResponse) bool {
	if response.Result.EndResult == models.SessionResultOK {
		w.log.Info().Msgf("%s session is completed with OK result", SmartIdWorkerName)

		user, err := w.handleCreateUser(ctx, response)
		if err != nil {
			w.log.Error().Err(err).Msgf("%s failed to create user", SmartIdWorkerName)
			return true
		}

		_, err = w.sessions.Update(ctx, models.Session{
			ID:     req.ID,
			UserId: user.ID,
			Status: AuthenticationSuccess,
		})
		if err != nil {
			w.log.Error().Err(err).Msgf("%s failed to update session", SmartIdWorkerName)
			return true
		}
	} else {
		w.log.Info().Msgf("%s session is completed with error", SmartIdWorkerName)

		if _, err := w.sessions.Update(ctx, models.Session{
			ID:     req.ID,
			Status: AuthenticationError,
			Error:  response.Result.EndResult,
		}); err != nil {
			w.log.Error().Err(err).Msgf("%s failed to update session", SmartIdWorkerName)
		}
	}

	return true
}

func (w *smartIdWorker) handleCreateUser(ctx context.Context, response *dto.SmartIdProviderSessionStatusResponse) (*models.User, error) {
	cert, err := w.certificate.Extract(response.Cert.Value)
	if err != nil {
		w.log.Error().Err(err).Msgf("%s failed to extract user from certificate", SmartIdWorkerName)
		return nil, err
	}

	return w.users.Create(ctx, &models.User{
		IdentityNumber: cert.IdentityNumber,
		PersonalCode:   cert.PersonalCode,
		FirstName:      cert.FirstName,
		LastName:       cert.LastName,
	})
}
