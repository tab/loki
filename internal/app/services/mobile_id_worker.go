package services

import (
	"context"
	"sync"

	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/config"
	"loki/pkg/logger"
)

type MobileIdWorker interface {
	Start(ctx context.Context)
	Stop()
}

type mobileIdWorker struct {
	cfg            *config.Config
	authentication Authentication
	sessions       Sessions
	users          Users
	queue          <-chan *MobileIdQueue
	wg             sync.WaitGroup
	log            *logger.Logger
}

func NewMobileIdWorker(
	cfg *config.Config,
	authentication Authentication,
	sessions Sessions,
	users Users,
	queue chan *MobileIdQueue,
	log *logger.Logger,
) MobileIdWorker {
	return &mobileIdWorker{
		cfg:            cfg,
		authentication: authentication,
		sessions:       sessions,
		users:          users,
		queue:          queue,
		log:            log,
	}
}

func (w *mobileIdWorker) Start(ctx context.Context) {
	w.log.Info().Msgf("MobileId::Worker starting in %s environment", w.cfg.AppEnv)

	w.wg.Add(1)
	go w.run(ctx)
}

func (w *mobileIdWorker) Stop() {
	w.log.Info().Msg("Stopping MobileId::Worker")
	w.wg.Wait()
}

func (w *mobileIdWorker) run(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("MobileId::Worker context cancelled, exiting")
			return
		case req, ok := <-w.queue:
			if !ok {
				w.log.Info().Msg("MobileId::Worker queue channel closed, exiting")
				return
			}
			w.perform(ctx, req)
		}
	}
}

func (w *mobileIdWorker) perform(ctx context.Context, req *MobileIdQueue) {
	w.log.Info().Msgf("MobileId::Worker perform %s", req.ID)

	for {
		response, err := w.authentication.GetMobileIdSessionStatus(ctx, req.ID)
		if err != nil {
			w.log.Error().Err(err).Msg("SmartId::Worker failed to get session status")
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
		w.log.Warn().Msg("Session is still running")
		return false
	default:
		w.log.Error().Msgf("Unknown session state: %s", response.State)
		return true
	}
}

func (w *mobileIdWorker) handleSessionComplete(ctx context.Context, req *MobileIdQueue, response *dto.MobileIdProviderSessionStatusResponse) bool {
	if response.Result == models.SessionResultOK {
		w.log.Info().Msg("MobileId::Worker session is completed with OK result")

		user, err := w.handleCreateUser(ctx, response)
		if err != nil {
			w.log.Error().Err(err).Msg("MobileId::Worker failed to create user")
			return true
		}

		_, err = w.sessions.Update(ctx, models.Session{
			ID:     req.ID,
			UserId: user.ID,
			Status: AuthenticationSuccess,
		})
		if err != nil {
			w.log.Error().Err(err).Msg("MobileId::Worker failed to update session")
			return true
		}
	} else {
		w.log.Info().Msg("MobileId::Worker session is completed with error")

		if _, err := w.sessions.Update(ctx, models.Session{
			ID:     req.ID,
			Status: AuthenticationError,
			Error:  response.Result,
		}); err != nil {
			w.log.Error().Err(err).Msg("MobileId::Worker failed to update session")
		}
	}

	return true
}

func (w *mobileIdWorker) handleCreateUser(ctx context.Context, response *dto.MobileIdProviderSessionStatusResponse) (user *models.User, err error) {
	cert, err := extractUserFromCertificate(response.Cert)
	if err != nil {
		w.log.Error().Err(err).Msg("MobileId::Worker failed to extract user from certificate")
		return nil, err
	}

	return w.users.Create(ctx, &models.User{
		IdentityNumber: cert.IdentityNumber,
		PersonalCode:   cert.PersonalCode,
		FirstName:      cert.FirstName,
		LastName:       cert.LastName,
	})
}
