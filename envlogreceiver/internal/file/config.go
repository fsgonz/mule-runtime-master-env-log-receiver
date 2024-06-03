package file

import (
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/statsconsumer"
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

// NewFileInputConfig creates a new BufferConfig with default values.
//
// Returns:
//   - *BufferConfig: A pointer to the newly created BufferConfig with default values.
func NewFileInputConfig() *BufferConfig {
	return NewFileInputConfigWithID(operatorType)
}

// NewFileInputConfigWithID creates a new BufferConfig with default values and a specific operator ID.
//
// Parameters:
//   - operatorID (string): The ID to be used for the operator.
//
// Returns:
//   - *BufferConfig: A pointer to the newly created BufferConfig with default values.
func NewFileInputConfigWithID(operatorID string) *BufferConfig {
	return &BufferConfig{
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
func (c BufferConfig) Build(set component.TelemetrySettings) (operator.Operator, error) {

	// Build the file statsconsumer with the specified configuration and emit function
	inputOperator, err := c.InputConfig.Build(set)
	if err != nil {
		return nil, err
	}

	// Build the file storage with the specified configuration and emit function

	input := &Input{
		InputOperator: inputOperator,
	}

	input.consumer, err = statsconsumer.Build(set, input.emit)

	if err != nil {
		return nil, err
	}

	return input, nil
}

// BufferConfig defines the configuration for the file input operator.
type BufferConfig struct {
	// InputConfig embeds the helper.InputConfig struct, which provides basic input operator configuration.
	helper.InputConfig `mapstructure:",squash"`

	// Config embeds the fileconsumer.Config struct, which provides configuration specific to file consumption.
	FileConsumerConfig `mapstructure:",squash"`

	// Config embeds the statsconsumer consuemr
	statsconsumer.StorageConsumerConfig `mapstructure:",squash"`

	id string
}

func (c BufferConfig) ID() string {
	return operatorType
}

func (c BufferConfig) Type() string {
	return c.id
}

func (c BufferConfig) SetID(ID string) {
	c.id = ID
}
