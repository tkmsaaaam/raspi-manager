package clamavreceiver

import (
	"time"

	"github.com/tkmsaaaam/raspi-manager/receiver/clamavreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	BufferInterval       string `mapstructure:"buffer_interval"`
	logFilePath          string `mapstructure:"log_file_path"`
	MetricsBuilderConfig metadata.MetricsBuilderConfig
}

func createDefaultConfig() component.Config {
	return &Config{
		BufferInterval:       "10s",
		logFilePath:          "/var/log/clamdscan.log",
		MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
	}
}

var _ component.Config = (*Config)(nil)

func (c Config) Validate() error {
	if _, err := time.ParseDuration(c.BufferInterval); err != nil {
		return err
	}
	return nil
}
