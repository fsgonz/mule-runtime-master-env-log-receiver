package envlogreceiver

import (
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/consumerretry"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/file"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/logsampler"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
)

import (
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/adapter"
)

// NewFactory creates a factory for receiver
func NewFactory() receiver.Factory {
	return adapter.NewFactory(ReceiverType{}, metadata.LogsStability)
}

// ReceiverType implements stanza.LogReceiverType
// to create a net usage stats receiver
type ReceiverType struct{}

// Type is the receiver type
func (f ReceiverType) Type() component.Type {
	return metadata.Type
}

// CreateDefaultConfig creates a config with type and version
func (f ReceiverType) CreateDefaultConfig() component.Config {
	return createDefaultConfig()
}

func createDefaultConfig() *OtelNetStatsReceiverConfig {
	return &OtelNetStatsReceiverConfig{
		BaseConfig: adapter.BaseConfig{
			Operators:      []operator.Config{},
			RetryOnFailure: consumerretry.NewDefaultConfig(),
		},
		InputConfig: *file.NewFileInputConfig(),
		LogSamplerConfig: logsampler.Config{
			LogSamplers: []logsampler.LogSampler{},
		},
	}
}

// BaseConfig gets the base config from config, for now
func (f ReceiverType) BaseConfig(cfg component.Config) adapter.BaseConfig {
	return cfg.(*OtelNetStatsReceiverConfig).BaseConfig
}

// OtelNetStatsReceiverConfig represents the configuration for the OpenTelemetry NetStats Logs Receiver.
type OtelNetStatsReceiverConfig struct {
	// InputConfig embeds the configuration for the network statistics input.
	InputConfig file.InputConfig `mapstructure:"file"`

	// BaseConfig embeds the base configuration for the logs receiver.
	adapter.BaseConfig `mapstructure:",squash"`

	// Log samplers
	LogSamplerConfig logsampler.Config `mapstructure:",squash"`
}

// InputConfig unmarshals the input operator
func (f ReceiverType) InputConfig(cfg component.Config) operator.Config {
	return operator.NewConfig(&cfg.(*OtelNetStatsReceiverConfig).InputConfig)
}

func (f ReceiverType) LogSamplers(cfg component.Config) logsampler.Config {
	return cfg.(*OtelNetStatsReceiverConfig).LogSamplerConfig
}

func (f ReceiverType) Input(cfg component.Config) file.Input {
	return cfg.(*OtelNetStatsReceiverConfig).InputConfig.Input()
}
