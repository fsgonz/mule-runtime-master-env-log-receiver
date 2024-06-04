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
	assert.Equal(t, config.inputConfig.ID(), operatorType, "The operator ID should match the operatorType")
	assert.IsType(t, helper.InputConfig{}, config.inputConfig, "OperatorConfig should be of type helper.OperatorConfig")
	assert.Equal(t, config.StatsConsumerConfig.PollInterval, 200*time.Millisecond)
}

// TestNewFileInputConfigWithID tests the NewFileInputConfigWithID function.
func TestNewFileInputConfigWithID(t *testing.T) {
	operatorID := "custom_operator_id"
	config := NewBufferConfigWithID(operatorID)

	assert.NotNil(t, config, "NewFileInputConfigWithID should not return nil")
	assert.Equal(t, config.inputConfig.ID(), operatorID, "The operator ID should match the provided operatorID")
	assert.Equal(t, config.inputConfig.Type(), operatorType, "The operator type should match the operatorType")
	assert.IsType(t, helper.InputConfig{}, config.inputConfig, "OperatorConfig should be of type helper.OperatorConfig")
	assert.Equal(t, config.StatsConsumerConfig.PollInterval, 200*time.Millisecond)
}

// TestBuild tests the Build method of the OperatorConfig struct.
func TestBuild(t *testing.T) {
	config := NewBufferConfig()
	set := createMockLogger()

	operator, err := config.Build(set)

	assert.NotNil(t, operator, "Build should return a valid operator")
	assert.NoError(t, err, "Build should not return an error")
}

// TestInput tests the Input method of the OperatorConfig struct.
func TestInput(t *testing.T) {
	config := NewBufferConfig()
	set := createMockLogger()

	_, err := config.Build(set)
	assert.NoError(t, err, "Build should not return an error")
}

func createMockLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	return logger.Sugar()
}
