package file

import (
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/storage"
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

// NewFileInputConfig creates a new BeufferConfig with default values.
//
// Returns:
//   - *BeufferConfig: A pointer to the newly created BeufferConfig with default values.
func NewFileInputConfig() *BeufferConfig {
	return NewFileInputConfigWithID(operatorType)
}

// NewFileInputConfigWithID creates a new BeufferConfig with default values and a specific operator ID.
//
// Parameters:
//   - operatorID (string): The ID to be used for the operator.
//
// Returns:
//   - *BeufferConfig: A pointer to the newly created BeufferConfig with default values.
func NewFileInputConfigWithID(operatorID string) *BeufferConfig {
	return &BeufferConfig{
		InputConfig:        helper.NewInputConfig(operatorID, operatorType),
		FileConsumerConfig: *NewConsumerConfig(),
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
func (c BeufferConfig) Build(set component.TelemetrySettings) (operator.Operator, error) {
	// Build the input operator from the configuration
	inputOperator, err := c.InputConfig.Build(set)
	if err != nil {
		return nil, err
	}

	// Build the file storage with the specified configuration and emit function

	input := &Input{
		InputOperator: inputOperator,
	}

	if c.FileConsumerConfig.Path != "" {
		c.input = input.InputOperator.WriterOperator
		input.consumer, err = c.FileConsumerConfig.Build(set, input.emit)
	} else {
		input.consumer, err = c.StorageConsumerConfig.Build(set, input.emit)
	}

	if err != nil {
		return nil, err
	}

	return input, nil
}

// Input returns the constructed Input operator from the BeufferConfig.
//
// Returns:
//   - Input: The input operator that was constructed from the BeufferConfig.
func (c BeufferConfig) Input() helper.WriterOperator {
	return c.input
}

// BeufferConfig defines the configuration for the file input operator.
type BeufferConfig struct {
	// InputConfig embeds the helper.InputConfig struct, which provides basic input operator configuration.
	helper.InputConfig `mapstructure:",squash"`

	// Config embeds the fileconsumer.Config struct, which provides configuration specific to file consumption.
	FileConsumerConfig `mapstructure:",squash"`

	// Config embeds the storage consuemr
	storage.StorageConsumerConfig `mapstructure:",squash"`

	// input holds the constructed Input operator.
	input helper.WriterOperator
}
