package clamavreceiver

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

type clamavReceiver struct {
	consumer consumer.Metrics
	settings receiver.Settings
	cancel   context.CancelFunc
	config   *Config
	sh       *clamavHandler
	obsrecv  *receiverhelper.ObsReport
}

func newClamavReceiver(config *Config, settings receiver.Settings, consumer consumer.Metrics) (*clamavReceiver, error) {
	obsrecv, err := receiverhelper.NewObsReport(receiverhelper.ObsReportSettings{
		ReceiverID:             settings.ID,
		Transport:              "event",
		ReceiverCreateSettings: settings,
	})
	if err != nil {
		return nil, err
	}
	return &clamavReceiver{
		consumer: consumer,
		settings: settings,
		config:   config,
		obsrecv:  obsrecv,
	}, nil
}

func (r *clamavReceiver) Start(ctx context.Context, _ component.Host) error {
	ctx, r.cancel = context.WithCancel(ctx)

	var err error
	r.sh, err = newClamavHandler(r.consumer, r.config, r.settings, r.obsrecv)
	if err != nil {
		return err
	}
	if r.sh == nil {
		return errors.New("failed to create clamav handler")
	}
	if err := r.sh.run(ctx); err != nil {
		return err
	}
	return nil
}

func (r *clamavReceiver) Shutdown(ctx context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}
	return nil
}
