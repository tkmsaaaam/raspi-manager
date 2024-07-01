package clamavreceiver

import (
	"errors"
	"time"

	"github.com/tkmsaaaam/raspi-manager/receiver/clamavreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	BufferInterval       string `mapstructure:"buffer_interval"`
	LogFilePath          string `mapstructure:"log_file_path"`
	MetricsBuilderConfig metadata.MetricsBuilderConfig
}

func createDefaultConfig() component.Config {
	return &Config{
		BufferInterval:       "10s",
		LogFilePath:          "/var/log/clamdscan.log",
		MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
	}
}

var _ component.Config = (*Config)(nil)

func (c Config) Validate() error {
	if _, err := time.ParseDuration(c.BufferInterval); err != nil {
		return err
	}
	if c.LogFilePath == "" {
		return errors.New("file path is not present")
	}
	return nil
}
