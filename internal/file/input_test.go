package file

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/entry"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer"
	"github.com/stretchr/testify/assert"
)

func TestInputEmit(t *testing.T) {
	t.Run("Empty Token", func(t *testing.T) {
		// Mock dependencies
		consumer := &fileconsumer.Manager{}
		input := &Input{consumer: consumer}

		// Create a context with timeout to ensure the test doesn't hang
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// Added a mocked operator
		mockOutputOperator := new(MockOutputOperator)
		mockOutputOperator.On("Process", mock.Anything, mock.Anything).Return(nil)
		input.OutputOperators = append(input.OutputOperators, mockOutputOperator)

		// Test case: token is empty and the mocked output operator was not called
		err := input.Emit(ctx, nil, map[string]interface{}{"key": "value"})
		assert.NoError(t, err)
		mockOutputOperator.AssertNotCalled(t, "Process", mock.Anything, mock.Anything)
	})

	t.Run("Non-Empty Token", func(t *testing.T) {
		// Mock dependencies
		consumer := &fileconsumer.Manager{}
		input := &Input{consumer: consumer}

		// Create a context with timeout to ensure the test doesn't hang
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// Added a mocked operator
		mockOutputOperator := new(MockOutputOperator)
		mockOutputOperator.On("Process", mock.Anything, mock.Anything).Return(nil)
		input.OutputOperators = append(input.OutputOperators, mockOutputOperator)

		// Test case: token is not empty and the mocked output operator was called
		token := []byte("test_token")
		err := input.Emit(ctx, token, map[string]interface{}{"key": "value"})
		mockOutputOperator.AssertCalled(t, "Process", mock.Anything, mock.Anything)
		assert.NoError(t, err)
	})
}

// MockOutputOperator is a mock implementation of the Operator interface
type MockOutputOperator struct {
	mock.Mock
}

func (m *MockOutputOperator) ID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockOutputOperator) Type() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockOutputOperator) Start(persister operator.Persister) error {
	args := m.Called(persister)
	return args.Error(0)
}

func (m *MockOutputOperator) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockOutputOperator) CanOutput() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockOutputOperator) Outputs() []operator.Operator {
	args := m.Called()
	return args.Get(0).([]operator.Operator)
}

func (m *MockOutputOperator) GetOutputIDs() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockOutputOperator) SetOutputs(ops []operator.Operator) error {
	args := m.Called(ops)
	return args.Error(0)
}

func (m *MockOutputOperator) SetOutputIDs(ids []string) {
	m.Called(ids)
}

func (m *MockOutputOperator) CanProcess() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockOutputOperator) Process(ctx context.Context, entry *entry.Entry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockOutputOperator) Logger() *zap.Logger {
	args := m.Called()
	return args.Get(0).(*zap.Logger)
}
