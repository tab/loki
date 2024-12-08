package services

import (
	"context"
	"sync"

	"loki/internal/app/models"
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
	queue          <-chan *SmartIdQueue
	wg             sync.WaitGroup
	log            *logger.Logger
}

func NewSmartIdWorker(
	cfg *config.Config,
	authentication Authentication,
	queue chan *SmartIdQueue,
	log *logger.Logger,
) SmartIdWorker {
	return &smartIdWorker{
		cfg:            cfg,
		authentication: authentication,
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

		switch response.State {
		case models.SESSION_COMPLETE:
			w.log.Info().Msg("SmartId::Worker session is complete")

			switch response.Result.EndResult {
			case models.SESSION_RESULT_OK:
				w.log.Info().Msg("SmartId::Worker session result is OK")

				_, err = w.authentication.UpdateSession(ctx, models.Session{
					ID:     req.ID,
					Status: models.SESSION_COMPLETE,
					Payload: models.SessionPayload{
						State:  response.State,
						Result: response.Result.EndResult,
						Cert:   response.Cert.Value,
					},
				})

				if err != nil {
					w.log.Error().Err(err).Msg("SmartId::Worker failed to update session")
					return
				}

			case models.SESSION_RESULT_USER_REFUSED:
				w.log.Info().Msg("User refused")
			case models.SESSION_RESULT_USER_REFUSED_DISPLAYTEXTANDPIN:
				w.log.Info().Msg("User refused display text and pin")
			case models.SESSION_RESULT_USER_REFUSED_VC_CHOICE:
				w.log.Info().Msg("User refused VC choice")
			case models.SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE:
				w.log.Info().Msg("User refused confirmation message")
			case models.SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE:
				w.log.Info().Msg("User refused confirmation message with VC choice")
			case models.SESSION_RESULT_USER_REFUSED_CERT_CHOICE:
				w.log.Info().Msg("User refused cert choice")
			case models.SESSION_RESULT_WRONG_VC:
				w.log.Info().Msg("Wrong VC")
			case models.SESSION_RESULT_TIMEOUT:
				w.log.Info().Msg("Timeout")
			default:
				w.log.Error().Msgf("Unknown session result: %s", response.Result.EndResult)
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
