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

func (w *smartIdWorker) handleUpdateSession(ctx context.Context, req *SmartIdQueue, response *dto.SmartIdProviderSessionStatusResponse) {
	status, message := buildSmartIdStatusAndMessage(response.Result.EndResult)
	payload := models.SessionPayload{
		State:  response.State,
		Result: response.Result.EndResult,
	}

	if response.Result.EndResult == models.SESSION_RESULT_OK {
		w.log.Info().Msg("SmartId::Worker session result is OK")
		payload.Signature = response.Signature.Value
		payload.Cert = response.Cert.Value
	} else if status == models.SESSION_ERROR {
		w.log.Info().Msgf("SmartId::Worker session error is %s", message)
	}

	if _, err := w.authentication.UpdateSession(ctx, models.Session{
		ID:      req.ID,
		Status:  status,
		Error:   message,
		Payload: payload,
	}); err != nil {
		w.log.Error().Err(err).Msg("SmartId::Worker failed to update session")
		return
	}
}

func buildSmartIdStatusAndMessage(endResult string) (status string, message string) {
	switch endResult {
	case models.SESSION_RESULT_OK:
		return models.SESSION_COMPLETE, ""
	case models.SESSION_RESULT_USER_REFUSED:
		return models.SESSION_ERROR, models.SESSION_RESULT_USER_REFUSED
	case models.SESSION_RESULT_USER_REFUSED_DISPLAYTEXTANDPIN:
		return models.SESSION_ERROR, models.SESSION_RESULT_USER_REFUSED_DISPLAYTEXTANDPIN
	case models.SESSION_RESULT_USER_REFUSED_VC_CHOICE:
		return models.SESSION_ERROR, models.SESSION_RESULT_USER_REFUSED_VC_CHOICE
	case models.SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE:
		return models.SESSION_ERROR, models.SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE
	case models.SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE:
		return models.SESSION_ERROR, models.SESSION_RESULT_USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE
	case models.SESSION_RESULT_USER_REFUSED_CERT_CHOICE:
		return models.SESSION_ERROR, models.SESSION_RESULT_USER_REFUSED_CERT_CHOICE
	case models.SESSION_RESULT_WRONG_VC:
		return models.SESSION_ERROR, models.SESSION_RESULT_WRONG_VC
	case models.SESSION_RESULT_TIMEOUT:
		return models.SESSION_ERROR, models.SESSION_RESULT_TIMEOUT
	default:
		return models.SESSION_ERROR, models.SESSION_RESULT_UNKNOWN
	}
}
