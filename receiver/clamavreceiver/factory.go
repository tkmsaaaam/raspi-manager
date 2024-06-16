package clamavreceiver

import (
	"context"

	"github.com/tkmsaaaam/raspi-manager/receiver/clamavreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability),
	)
}

func createMetricsReceiver(
	_ context.Context,
	settings receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	if consumer == nil {
		return nil, component.ErrDataTypeIsNotSupported
	}

	c := cfg.(*Config)
	err := c.Validate()
	if err != nil {
		return nil, err
	}
	return newClamavReceiver(c, settings, consumer)
}
