package envlogreceiver

import (
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/consumerretry"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/file"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/logsampler"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
)

import (
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/adapter"
)

const (
	operatorType = "stats_input"
)

// NewFactory creates a factory for receiver
func NewFactory() receiver.Factory {
	return adapter.NewFactory(ReceiverType{}, metadata.LogsStability)
}

func createDefaultConfig() *OtelNetStatsReceiverConfig {
	return &OtelNetStatsReceiverConfig{
		BaseConfig: adapter.BaseConfig{
			Operators:      []operator.Config{},
			RetryOnFailure: consumerretry.NewDefaultConfig(),
		},
		InputConfig:  helper.NewInputConfig(operatorType, operatorType),
		BufferConfig: *file.NewFileInputConfig(),
		LogSamplerConfig: logsampler.Config{
			LogSamplers: []logsampler.LogSampler{},
		},
	}
}

// OtelNetStatsReceiverConfig represents the configuration for the OpenTelemetry NetStats Logs Receiver.
type OtelNetStatsReceiverConfig struct {
	// InputConfig provides a basic implementation of an input operator config.
	InputConfig helper.InputConfig `mapstructure:",squash"`

	// InputConfig embeds the configuration for the network statistics input.
	BufferConfig file.BufferConfig `mapstructure:"buffer"`

	// BaseConfig embeds the base configuration for the logs receiver.
	adapter.BaseConfig `mapstructure:",squash"`

	// Log samplers
	LogSamplerConfig logsampler.Config `mapstructure:",squash"`
}

// ReceiverType implements stanza.LogReceiverType
// to create a net usage stats receiver
type ReceiverType struct{}

// InputConfig unmarshals the input operator
func (f ReceiverType) InputConfig(cfg component.Config) operator.Config {
	return operator.NewConfig(&cfg.(*OtelNetStatsReceiverConfig).BufferConfig)
}

func (f ReceiverType) LogSamplers(cfg component.Config) logsampler.Config {
	return cfg.(*OtelNetStatsReceiverConfig).LogSamplerConfig
}

// BaseConfig gets the base config from config, for now
func (f ReceiverType) BaseConfig(cfg component.Config) adapter.BaseConfig {
	return cfg.(*OtelNetStatsReceiverConfig).BaseConfig
}

// Type is the receiver type
func (f ReceiverType) Type() component.Type {
	return metadata.Type
}

func (f ReceiverType) BufferConfig(cfg component.Config) *file.BufferConfig {
	return &cfg.(*OtelNetStatsReceiverConfig).BufferConfig
}

// CreateDefaultConfig creates a config with type and version
func (f ReceiverType) CreateDefaultConfig() component.Config {
	return createDefaultConfig()
}
