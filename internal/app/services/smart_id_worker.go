package services

import (
	"context"
	"sync"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/config"
	"loki/pkg/logger"
)

type SmartIdWorker interface {
	Start(ctx context.Context)
	Stop()
}

type smartIdWorker struct {
	cfg            *config.Config
	authentication Authentication
	sessions       Sessions
	users          Users
	queue          <-chan *SmartIdQueue
	wg             sync.WaitGroup
	log            *logger.Logger
}

func NewSmartIdWorker(
	cfg *config.Config,
	authentication Authentication,
	sessions Sessions,
	users Users,
	queue chan *SmartIdQueue,
	log *logger.Logger,
) SmartIdWorker {
	return &smartIdWorker{
		cfg:            cfg,
		authentication: authentication,
		sessions:       sessions,
		users:          users,
		queue:          queue,
		log:            log,
	}
}

func (w *smartIdWorker) Start(ctx context.Context) {
	w.log.Info().Msgf("SmartId::Worker starting in %s environment", w.cfg.AppEnv)

	w.wg.Add(1)
	go w.run(ctx)
}

func (w *smartIdWorker) Stop() {
	w.log.Info().Msg("Stopping SmartId::Worker")
	w.wg.Wait()
}

func (w *smartIdWorker) run(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("SmartId::Worker context cancelled, exiting")
			return
		case req, ok := <-w.queue:
			if !ok {
				w.log.Info().Msg("SmartId::Worker queue channel closed, exiting")
				return
			}
			w.perform(ctx, req)
		}
	}
}

func (w *smartIdWorker) perform(ctx context.Context, req *SmartIdQueue) {
	w.log.Info().Msgf("SmartId::Worker perform %s", req.ID)

	for {
		response, err := w.authentication.GetSmartIdSessionStatus(ctx, req.ID)
		if err != nil {
			w.log.Error().Err(err).Msg("SmartId::Worker failed to get session status")
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
		w.log.Warn().Msg("Session is still running")
		return false
	default:
		w.log.Error().Msgf("Unknown session state: %s", response.State)
		return true
	}
}

func (w *smartIdWorker) handleSessionComplete(ctx context.Context, req *SmartIdQueue, response *dto.SmartIdProviderSessionStatusResponse) bool {
	if response.Result.EndResult == models.SessionResultOK {
		w.log.Info().Msg("SmartId::Worker session is completed with OK result")

		user, err := w.handleCreateUser(ctx, response)
		if err != nil {
			w.log.Error().Err(err).Msg("SmartId::Worker failed to create user")
			return true
		}

		_, err = w.sessions.Update(ctx, models.Session{
			ID:     req.ID,
			UserId: user.ID,
			Status: AuthenticationSuccess,
		})
		if err != nil {
			w.log.Error().Err(err).Msg("SmartId::Worker failed to update session")
			return true
		}
	} else {
		w.log.Info().Msg("SmartId::Worker session is completed with error")

		if _, err := w.sessions.Update(ctx, models.Session{
			ID:     req.ID,
			Status: AuthenticationError,
			Error:  response.Result.EndResult,
		}); err != nil {
			w.log.Error().Err(err).Msg("SmartId::Worker failed to update session")
		}
	}

	return true
}

func (w *smartIdWorker) handleCreateUser(ctx context.Context, response *dto.SmartIdProviderSessionStatusResponse) (*models.User, error) {
	cert, err := extractUserFromCertificate(response.Cert.Value)
	if err != nil {
		w.log.Error().Err(err).Msg("SmartId::Worker failed to extract user from certificate")
		return nil, err
	}

	return w.users.Create(ctx, &models.User{
		IdentityNumber: cert.IdentityNumber,
		PersonalCode:   cert.PersonalCode,
		FirstName:      cert.FirstName,
		LastName:       cert.LastName,
	})
}
