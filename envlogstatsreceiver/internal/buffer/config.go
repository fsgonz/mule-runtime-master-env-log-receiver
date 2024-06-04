package buffer

import (
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogstatsreceiver/internal/logsampler"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogstatsreceiver/internal/statsconsumer"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"go.uber.org/zap"
)

const (
	operatorType = "stats_buffer"
)

func init() {
	operator.Register(operatorType, func() operator.Builder { return NewBufferConfig() })
}

// NewBufferConfig creates a new BufferConfig with default values.
//
// Returns:
//   - *BufferConfig: A pointer to the newly created BufferConfig with default values.
func NewBufferConfig() *BufferConfig {
	return NewBufferConfigWithID(operatorType)
}

// NewBufferConfigWithID creates a new BufferConfig with default values and a specific operator ID.
//
// Parameters:
//   - operatorID (string): The ID to be used for the operator.
//
// Returns:
//   - *BufferConfig: A pointer to the newly created BufferConfig with default values.
func NewBufferConfigWithID(operatorID string) *BufferConfig {
	return &BufferConfig{
		InputConfig:         helper.NewInputConfig(operatorID, operatorType),
		FileConsumerConfig:  *statsconsumer.NewConsumerConfig(),
		StatsConsumerConfig: statsconsumer.NewDefaultStatsConsumerConfig(),
	}
}

// Build constructs the buffer input operator from the supplied configuration.
//
// Parameters:
//   - set (component.TelemetrySettings): The telemetry settings to be used during the build process.
//
// Returns:
//   - operator.Operator: The constructed buffer input operator.
//   - error: An error that occurred during the build process, or nil if the build was successful.
func (c BufferConfig) Build(logger *zap.SugaredLogger) (operator.Operator, error) {

	// Build the buffer statsconsumer with the specified configuration and emit function
	inputOperator, err := c.InputConfig.Build(logger)
	if err != nil {
		return nil, err
	}

	// Build the buffer storage with the specified configuration and emit function

	input := &Input{
		InputOperator: inputOperator,
	}

	input.consumer, err = statsconsumer.Build(logger, c.logSampler, input.emit, c.FileConsumerConfig)

	if err != nil {
		return nil, err
	}

	return input, nil
}

// BufferConfig defines the configuration for the buffer input operator.
type BufferConfig struct {
	// InputConfig embeds the helper.InputConfig struct, which provides basic input operator configuration.
	helper.InputConfig `mapstructure:",squash"`

	// Config embeds the fileconsumer.Config struct, which provides configuration specific to buffer consumption.
	statsconsumer.FileConsumerConfig `mapstructure:",squash"`

	// Config embeds the statsconsumer consuemr
	statsconsumer.StatsConsumerConfig `mapstructure:",squash"`

	id string

	logSampler logsampler.LogSampler
}

func (c BufferConfig) ID() string {
	return operatorType
}

func (c *BufferConfig) Type() string {
	return c.id
}

func (c *BufferConfig) SetID(ID string) {
	c.id = ID
}

func (c *BufferConfig) SetLogSamplerConfig(LogSamplerConfig logsampler.Config) {
	c.logSampler = LogSamplerConfig.LogSamplers[0]
}
