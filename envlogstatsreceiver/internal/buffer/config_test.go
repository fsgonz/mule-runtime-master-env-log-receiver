package buffer

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

// TestNewFileInputConfig tests the NewFileInputConfig function.
func TestNewFileInputConfig(t *testing.T) {
	config := NewBufferConfig()

	assert.NotNil(t, config, "NewFileInputConfig should not return nil")
	assert.Equal(t, config.InputConfig.ID(), operatorType, "The operator ID should match the operatorType")
	assert.IsType(t, helper.InputConfig{}, config.InputConfig, "InputConfig should be of type helper.InputConfig")
	assert.Equal(t, config.StatsConsumerConfig.PollInterval, 200*time.Millisecond)
}

// TestNewFileInputConfigWithID tests the NewFileInputConfigWithID function.
func TestNewFileInputConfigWithID(t *testing.T) {
	operatorID := "custom_operator_id"
	config := NewBufferConfigWithID(operatorID)

	assert.NotNil(t, config, "NewFileInputConfigWithID should not return nil")
	assert.Equal(t, config.InputConfig.ID(), operatorID, "The operator ID should match the provided operatorID")
	assert.Equal(t, config.InputConfig.Type(), operatorType, "The operator type should match the operatorType")
	assert.IsType(t, helper.InputConfig{}, config.InputConfig, "InputConfig should be of type helper.InputConfig")
	assert.Equal(t, config.StatsConsumerConfig.PollInterval, 200*time.Millisecond)
}

// TestBuild tests the Build method of the InputConfig struct.
func TestBuild(t *testing.T) {
	config := NewBufferConfig()
	logger := createMockTelemetrySettings()

	operator, err := config.Build(&logger)

	assert.NotNil(t, operator, "Build should return a valid operator")
	assert.NoError(t, err, "Build should not return an error")
}

// TestInput tests the Input method of the InputConfig struct.
func TestInput(t *testing.T) {
	config := NewBufferConfig()
	logger := createMockTelemetrySettings()

	_, err := config.Build(&logger)
	assert.NoError(t, err, "Build should not return an error")
}

func createMockTelemetrySettings() zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	return *logger.Sugar()
}
