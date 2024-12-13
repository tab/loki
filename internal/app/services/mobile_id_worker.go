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
	queue          <-chan *MobileIdQueue
	wg             sync.WaitGroup
	log            *logger.Logger
}

func NewMobileIdWorker(
	cfg *config.Config,
	authentication Authentication,
	sessions Sessions,
	queue chan *MobileIdQueue,
	log *logger.Logger,
) MobileIdWorker {
	return &mobileIdWorker{
		cfg:            cfg,
		authentication: authentication,
		sessions:       sessions,
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

		switch response.State {
		case models.SESSION_COMPLETE:
			w.handleUpdateSession(ctx, req, response)
			return
		case models.SESSION_RUNNING:
			w.log.Warn().Msg("Session is still running")
			continue
		default:
			w.log.Error().Msgf("Unknown session state: %s", response.State)
			return
		}
	}
}

func (w *mobileIdWorker) handleUpdateSession(ctx context.Context, req *MobileIdQueue, response *dto.MobileIdProviderSessionStatusResponse) {
	status, message := buildMobileIdStatusAndMessage(response.Result)
	payload := models.SessionPayload{
		State:  response.State,
		Result: response.Result,
	}

	if response.Result == models.SESSION_RESULT_OK {
		w.log.Info().Msg("MobileId::Worker session result is OK")
		payload.Signature = response.Signature.Value
		payload.Cert = response.Cert
	} else if status == models.SESSION_ERROR {
		w.log.Info().Msgf("MobileId::Worker session error is %s", message)
	}

	if _, err := w.sessions.Update(ctx, models.Session{
		ID:      req.ID,
		Status:  status,
		Error:   message,
		Payload: payload,
	}); err != nil {
		w.log.Error().Err(err).Msg("MobileId::Worker failed to update session")
		return
	}
}

func buildMobileIdStatusAndMessage(endResult string) (status string, message string) {
	switch endResult {
	case models.SESSION_RESULT_OK:
		return models.SESSION_COMPLETE, ""
	case models.SESSION_RESULT_NOT_MID_CLIENT:
		return models.SESSION_ERROR, models.SESSION_RESULT_NOT_MID_CLIENT
	case models.SESSION_RESULT_USER_CANCELLED:
		return models.SESSION_ERROR, models.SESSION_RESULT_USER_CANCELLED
	case models.SESSION_RESULT_SIGNATURE_HASH_MISMATCH:
		return models.SESSION_ERROR, models.SESSION_RESULT_SIGNATURE_HASH_MISMATCH
	case models.SESSION_RESULT_PHONE_ABSENT:
		return models.SESSION_ERROR, models.SESSION_RESULT_PHONE_ABSENT
	case models.SESSION_RESULT_DELIVERY_ERROR:
		return models.SESSION_ERROR, models.SESSION_RESULT_DELIVERY_ERROR
	case models.SESSION_RESULT_SIM_ERROR:
		return models.SESSION_ERROR, models.SESSION_RESULT_SIM_ERROR
	case models.SESSION_RESULT_TIMEOUT:
		return models.SESSION_ERROR, models.SESSION_RESULT_TIMEOUT
	default:
		return models.SESSION_ERROR, models.SESSION_RESULT_UNKNOWN
	}
}
