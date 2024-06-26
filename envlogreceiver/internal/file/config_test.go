package file

import (
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
)

// TestNewFileInputConfig tests the NewFileInputConfig function.
func TestNewFileInputConfig(t *testing.T) {
	config := NewFileInputConfig()

	assert.NotNil(t, config, "NewFileInputConfig should not return nil")
	assert.Equal(t, config.InputConfig.ID(), operatorType, "The operator ID should match the operatorType")
	assert.IsType(t, helper.InputConfig{}, config.InputConfig, "InputConfig should be of type helper.InputConfig")
	assert.IsType(t, fileconsumer.Config{}, config.Config, "Config should be of type fileconsumer.Config")
}

// TestNewFileInputConfigWithID tests the NewFileInputConfigWithID function.
func TestNewFileInputConfigWithID(t *testing.T) {
	operatorID := "custom_operator_id"
	config := NewFileInputConfigWithID(operatorID)

	assert.NotNil(t, config, "NewFileInputConfigWithID should not return nil")
	assert.Equal(t, config.InputConfig.ID(), operatorID, "The operator ID should match the provided operatorID")
	assert.Equal(t, config.InputConfig.Type(), operatorType, "The operator type should match the operatorType")
	assert.IsType(t, helper.InputConfig{}, config.InputConfig, "InputConfig should be of type helper.InputConfig")
	assert.IsType(t, fileconsumer.Config{}, config.Config, "Config should be of type fileconsumer.Config")
}

// TestBuild tests the Build method of the InputConfig struct.
func TestBuild(t *testing.T) {
	config := NewFileInputConfig()
	config.Criteria.Include = append(config.Criteria.Include, "dummy")
	set := createMockTelemetrySettings()

	operator, err := config.Build(set)

	assert.NotNil(t, operator, "Build should return a valid operator")
	assert.NoError(t, err, "Build should not return an error")
}

// TestInput tests the Input method of the InputConfig struct.
func TestInput(t *testing.T) {
	config := NewFileInputConfig()
	config.Criteria.Include = append(config.Criteria.Include, "dummy")
	set := createMockTelemetrySettings()

	_, err := config.Build(set)
	assert.NoError(t, err, "Build should not return an error")

	input := config.Input()
	assert.NotNil(t, input, "Input should return a valid Input operator")
}

func createMockTelemetrySettings() component.TelemetrySettings {
	logger, _ := zap.NewDevelopment()
	return component.TelemetrySettings{
		Logger:         logger,
		TracerProvider: trace.NewNoopTracerProvider(),
		MetricsLevel:   configtelemetry.LevelNone,
		Resource:       pcommon.NewResource(),
		ReportStatus:   func(*component.StatusEvent) {},
	}
}
