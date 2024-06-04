package envlogstatsreceiver

import (
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogstatsreceiver/internal/adapter"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogstatsreceiver/internal/buffer"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogstatsreceiver/internal/consumerretry"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogstatsreceiver/internal/logsampler"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogstatsreceiver/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
)

const (
	operatorType = "stats_input_type"
	operatorId   = "stats_input_type"
)

// NewFactory creates a new factory for the Environment Log Stats Receiver.
//
// This factory is used by the OpenTelemetry Collector to create instances of the receiver
// with the specified configuration.
//
// Returns:
// - A receiver.Factory configured to create Environment Log Stats Receivers.
func NewFactory() receiver.Factory {
	return adapter.NewFactory(ReceiverType{}, metadata.LogsStability)
}

// createDefaultConfig creates and returns a default configuration for the Environment Log Stats Receiver.
//
// The default configuration includes:
// - BaseConfig with empty operator configurations and a default retry-on-failure policy.
// - OperatorConfig initialized with default operator types.
// - BufferConfig with default buffer buffer settings.
// - LogSamplerConfig with no log samplers.
//
// Returns:
// - A pointer to an EnvLogStatsReceiverReceiverConfig struct initialized with default values.
func createDefaultConfig() *EnvLogStatsReceiverReceiverConfig {
	return &EnvLogStatsReceiverReceiverConfig{
		BaseConfig: adapter.BaseConfig{
			Operators:      []operator.Config{},
			RetryOnFailure: consumerretry.NewDefaultConfig(),
		},
		InputConfig:  helper.NewInputConfig(operatorId, operatorType),
		BufferConfig: *buffer.NewBufferConfig(),
		LogSamplerConfig: logsampler.Config{
			LogSamplers: []logsampler.LogSampler{},
		},
	}
}

// EnvLogStatsReceiverReceiverConfig represents the configuration for the Environment Log Stats Receiver.
type EnvLogStatsReceiverReceiverConfig struct {
	// InputConfig provides a basic implementation of an input operator config.
	InputConfig helper.InputConfig `mapstructure:",squash"`

	// BufferConfig embeds the configuration for the buffer config.
	BufferConfig buffer.BufferConfig `mapstructure:"buffer"`

	// BaseConfig embeds the base configuration for the logs receiver.
	adapter.BaseConfig `mapstructure:",squash"`

	// LogSamplerConfig embeds the configuration for log samplers.
	LogSamplerConfig logsampler.Config `mapstructure:",squash"`
}

// ReceiverType implements stanza.LogReceiverType to create a net usage stats receiver.
type ReceiverType struct{}

// InputConfig unmarshals the input operator configuration.
//
// Params:
// - cfg: The component configuration.
//
// Returns:
// - An operator.Config configured with the buffer settings from the provided configuration.
func (f ReceiverType) OperatorConfig(cfg component.Config) operator.Config {
	return operator.NewConfig(&cfg.(*EnvLogStatsReceiverReceiverConfig).BufferConfig)
}

func (f ReceiverType) InputConfig(cfg component.Config) helper.InputConfig {
	return cfg.(*EnvLogStatsReceiverReceiverConfig).InputConfig
}

// LogSamplers returns the log sampler configuration from the provided component configuration.
//
// Params:
// - cfg: The component configuration.
//
// Returns:
// - A logsampler.Config containing the log sampler settings.
func (f ReceiverType) LogSamplers(cfg component.Config) logsampler.Config {
	return cfg.(*EnvLogStatsReceiverReceiverConfig).LogSamplerConfig
}

// BaseConfig gets the base configuration from the provided component configuration.
//
// Params:
// - cfg: The component configuration.
//
// Returns:
// - An adapter.BaseConfig containing the base settings.
func (f ReceiverType) BaseConfig(cfg component.Config) adapter.BaseConfig {
	return cfg.(*EnvLogStatsReceiverReceiverConfig).BaseConfig
}

// Type returns the type of the receiver.
//
// Returns:
// - The component.Type representing the receiver type.
func (f ReceiverType) Type() component.Type {
	return metadata.Type
}

// BufferConfig returns the buffer configuration from the provided component configuration.
//
// Params:
// - cfg: The component configuration.
//
// Returns:
// - A pointer to the buffer.BufferConfig containing the buffer settings.
func (f ReceiverType) BufferConfig(cfg component.Config) *buffer.BufferConfig {
	return &cfg.(*EnvLogStatsReceiverReceiverConfig).BufferConfig
}

// CreateDefaultConfig creates and returns the default configuration for the receiver.
//
// Returns:
// - A component.Config initialized with default values.
func (f ReceiverType) CreateDefaultConfig() component.Config {
	return createDefaultConfig()
}
