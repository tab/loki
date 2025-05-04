package authentication

import (
	"time"

	"github.com/tab/mobileid"
	"github.com/tab/smartid"
	"go.uber.org/fx"

	"loki/internal/config"
	"loki/internal/config/logger"
)

const (
	Concurrency = 5
	QueueSize   = 15
)

var Module = fx.Options(
	fx.Provide(
		func(cfg *config.Config, log *logger.Logger) (smartid.Client, error) {
			certManager, err := smartid.NewCertificateManager(cfg.CertPath)
			if err != nil {
				return nil, err
			}
			client := smartid.NewClient().
				WithRelyingPartyName(cfg.SmartId.RelyingPartyName).
				WithRelyingPartyUUID(cfg.SmartId.RelyingPartyUUID).
				WithCertificateLevel("QUALIFIED").
				WithHashType("SHA512").
				WithInteractionType("displayTextAndPIN").
				WithText(cfg.SmartId.Text).
				WithURL(cfg.SmartId.BaseURL).
				WithTimeout(60 * time.Second).
				WithTLSConfig(certManager.TLSConfig())
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

	fx.Provide(
		func(cfg *config.Config, log *logger.Logger) (mobileid.Client, error) {
			client := mobileid.NewClient().
				WithRelyingPartyName(cfg.MobileId.RelyingPartyName).
				WithRelyingPartyUUID(cfg.MobileId.RelyingPartyUUID).
				WithHashType("SHA512").
				WithText(cfg.MobileId.Text).
				WithTextFormat(cfg.MobileId.TextFormat).
				WithLanguage(cfg.MobileId.Language).
				WithURL(cfg.MobileId.BaseURL).
				WithTimeout(60 * time.Second)
			if err := client.Validate(); err != nil {
				return nil, err
			}
			return client, nil
		},
	),
	fx.Provide(
		func(client mobileid.Client) mobileid.Worker {
			worker := mobileid.NewWorker(client).
				WithConcurrency(Concurrency).
				WithQueueSize(QueueSize)

			return worker
		},
	),
	fx.Provide(NewMobileId),

)
