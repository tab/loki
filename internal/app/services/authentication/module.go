package authentication

import (
	"time"

	"github.com/tab/smartid"
	"go.uber.org/fx"

	"loki/internal/config"
	"loki/pkg/logger"
)

const (
	Concurrency = 5
	QueueSize   = 15
)

var Module = fx.Options(
	fx.Provide(
		func(cfg *config.Config, log *logger.Logger) (smartid.Client, error) {
			client := smartid.NewClient().
				WithRelyingPartyName(cfg.SmartId.RelyingPartyName).
				WithRelyingPartyUUID(cfg.SmartId.RelyingPartyUUID).
				WithCertificateLevel("QUALIFIED").
				WithHashType("SHA512").
				WithInteractionType("displayTextAndPIN").
				WithText(cfg.SmartId.Text).
				WithURL(cfg.SmartId.BaseURL).
				WithTimeout(60 * time.Second)
			if err := client.Validate(); err != nil {
				return nil, err
			}
			return client, nil
		},
	),
	fx.Provide(
		func(client smartid.Client) smartid.Worker {
			worker := smartid.NewWorker(client).
				WithConcurrency(Concurrency).
				WithQueueSize(QueueSize)

			return worker
		},
	),
	fx.Provide(NewSmartId),
)
