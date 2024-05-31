package file

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"go.opentelemetry.io/collector/component"
)

const (
	operatorType = "file_input"
)

func init() {
	operator.Register(operatorType, func() operator.Builder { return NewFileInputConfig() })
}

// NewFileInputConfig creates a new InputConfig with default values.
//
// Returns:
//   - *InputConfig: A pointer to the newly created InputConfig with default values.
func NewFileInputConfig() *InputConfig {
	return NewFileInputConfigWithID(operatorType)
}

// NewFileInputConfigWithID creates a new InputConfig with default values and a specific operator ID.
//
// Parameters:
//   - operatorID (string): The ID to be used for the operator.
//
// Returns:
//   - *InputConfig: A pointer to the newly created InputConfig with default values.
func NewFileInputConfigWithID(operatorID string) *InputConfig {
	return &InputConfig{
		InputConfig:    helper.NewInputConfig(operatorID, operatorType),
		ConsumerConfig: *NewConsumerConfig(),
	}
}

// Build constructs the file input operator from the supplied configuration.
//
// Parameters:
//   - set (component.TelemetrySettings): The telemetry settings to be used during the build process.
//
// Returns:
//   - operator.Operator: The constructed file input operator.
//   - error: An error that occurred during the build process, or nil if the build was successful.
func (c InputConfig) Build(set component.TelemetrySettings) (operator.Operator, error) {
	// Build the input operator from the configuration
	inputOperator, err := c.InputConfig.Build(set)
	if err != nil {
		return nil, err
	}

	input := &Input{
		InputOperator: inputOperator,
	}

	c.input = input.InputOperator.WriterOperator

	// Build the file consumer with the specified configuration and emit function
	input.consumer, err = c.ConsumerConfig.Build(set, input.emit)
	if err != nil {
		return nil, err
	}

	return input, nil
}

// Input returns the constructed Input operator from the InputConfig.
//
// Returns:
//   - Input: The input operator that was constructed from the InputConfig.
func (c InputConfig) Input() helper.WriterOperator {
	return c.input
}

// InputConfig defines the configuration for the file input operator.
type InputConfig struct {
	// InputConfig embeds the helper.InputConfig struct, which provides basic input operator configuration.
	helper.InputConfig `mapstructure:",squash"`

	// Config embeds the fileconsumer.Config struct, which provides configuration specific to file consumption.
	ConsumerConfig `mapstructure:",squash"`

	// input holds the constructed Input operator.
	input helper.WriterOperator
}
