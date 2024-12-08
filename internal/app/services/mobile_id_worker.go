package services

import (
	"context"
	"sync"

	"loki/internal/app/models"
	"loki/internal/config"
	"loki/pkg/logger"
)

type MobileIdWorker interface {
	Start(ctx context.Context)
	Stop()
}

type mobileIDWorker struct {
	cfg            *config.Config
	authentication Authentication
	queue          <-chan *MobileIdQueue
	wg             sync.WaitGroup
	log            *logger.Logger
}

func NewMobileIdWorker(
	cfg *config.Config,
	authentication Authentication,
	queue chan *MobileIdQueue,
	log *logger.Logger,
) MobileIdWorker {
	return &mobileIDWorker{
		cfg:            cfg,
		authentication: authentication,
		queue:          queue,
		log:            log,
	}
}

func (w *mobileIDWorker) Start(ctx context.Context) {
	w.log.Info().Msgf("MobileId::Worker starting in %s environment", w.cfg.AppEnv)

	w.wg.Add(1)
	go w.run(ctx)
}

func (w *mobileIDWorker) Stop() {
	w.log.Info().Msg("Stopping MobileId::Worker")
	w.wg.Wait()
}

func (w *mobileIDWorker) run(ctx context.Context) {
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

func (w *mobileIDWorker) perform(ctx context.Context, req *MobileIdQueue) {
	w.log.Info().Msgf("MobileId::Worker perform %s", req.ID)

	for {
		response, err := w.authentication.GetMobileIdSessionStatus(ctx, req.ID)
		if err != nil {
			w.log.Error().Err(err).Msg("MobileId::Worker failed to get session status")
			return
		}

		switch response.State {
		case models.SESSION_COMPLETE:
			w.log.Info().Msg("MobileId::Worker session is complete")

			switch response.Result {
			case models.SESSION_RESULT_OK:
				w.log.Info().Msg("MobileId::Worker session result is OK")

				_, err = w.authentication.UpdateSession(ctx, models.Session{
					ID:     req.ID,
					Status: models.SESSION_COMPLETE,
					Payload: models.SessionPayload{
						State:  response.State,
						Result: response.Result,
						Cert:   response.Cert,
					},
				})

				if err != nil {
					w.log.Error().Err(err).Msg("MobileId::Worker failed to update session")
					return
				}

			case models.SESSION_RESULT_NOT_MID_CLIENT:
				w.log.Info().Msg("User is not a Mobile ID client")
			case models.SESSION_RESULT_USER_CANCELLED:
				w.log.Info().Msg("User cancelled the authentication")
			case models.SESSION_RESULT_SIGNATURE_HASH_MISMATCH:
				w.log.Info().Msg("Signature hash mismatch")
			case models.SESSION_RESULT_PHONE_ABSENT:
				w.log.Info().Msg("Phone is absent")
			case models.SESSION_RESULT_DELIVERY_ERROR:
				w.log.Info().Msg("SMS delivery error")
			case models.SESSION_RESULT_SIM_ERROR:
				w.log.Info().Msg("SIM error")
			case models.SESSION_RESULT_TIMEOUT:
				w.log.Info().Msg("Timeout")
			default:
				w.log.Error().Msgf("Unknown session result: %s", response.Result)
			}

			return
		case models.SESSION_RUNNING:
			w.log.Warn().Msg("Session is still running")
			continue
		default:
			w.log.Error().Msgf("Unknown session state: %s", response)
			return
		}
	}
}
