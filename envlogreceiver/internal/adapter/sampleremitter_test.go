package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/file"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/logsampler"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/entry"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

// Event represents the "events" array in the JSON.
type Event struct {
	ID         string `json:"id"`
	Timestamp  int64  `json:"timestamp"`
	RootOrgID  string `json:"root_org_id"`
	OrgID      string `json:"org_id"`
	EnvID      string `json:"env_id"`
	AssetID    string `json:"asset_id"`
	WorkerID   string `json:"worker_id"`
	UsageBytes uint64 `json:"usage_bytes"`
	Billable   bool   `json:"billable"`
}

// LogEntry represents the entire JSON structure.
type LogEntry struct {
	Format   string            `json:"format"`
	Time     int64             `json:"time"`
	Events   []Event           `json:"events"`
	Metadata map[string]string `json:"metadata"`
}

// MockEmitter is a mock implementation of helper.LogEmitter
type MockEmitter struct{}

func (m *MockEmitter) EmitLog(context.Context, *entry.Entry) error {
	return nil
}

// mockPersister is a mock implementation of operator.Persister
type mockPersister struct {
	data map[string][]byte
}

func (m *mockPersister) Get(ctx context.Context, key string) ([]byte, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return nil, nil
}

func (m *mockPersister) Set(ctx context.Context, key string, value []byte) error {
	m.data[key] = value
	return nil
}

// mockSampler is a mock implementation of sampler.Sampler
type mockSampler struct{}

func (m *mockSampler) Sample() (uint64, error) {
	return 100, nil
}

func TestSamplerEmitterFactory(t *testing.T) {
	t.Run("OutputFileLogger", func(t *testing.T) {
		// Prepare mock data
		mockPersister := &MockPersister{
			Data: make(map[string][]byte),
		}

		// Set up some mock data
		mockPersister.Data["test.key"] = []byte("value")

		mockEmitter := &helper.LogEmitter{}
		mockInput := &file.Input{}

		// Call SamplerEmitterFactory
		samplerEmitter, err := SamplerEmitterFactory(logsampler.OutputFileLogger, "test.log", mockPersister, mockEmitter, *mockInput)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, samplerEmitter)
		assert.IsType(t, FileLoggerSamplerEmitter{}, samplerEmitter)
	})

	t.Run("OutputPipelineEmitter", func(t *testing.T) {
		// Prepare mock data
		mockPersister := &MockPersister{
			Data: make(map[string][]byte),
		}

		// Set up some mock data
		mockPersister.Data["test.key"] = []byte("value")

		mockEmitter := &helper.LogEmitter{}
		mockInput := &file.Input{}

		// Call SamplerEmitterFactory
		samplerEmitter, err := SamplerEmitterFactory(logsampler.OutputPipelineEmitter, "test.log", mockPersister, mockEmitter, *mockInput)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, samplerEmitter)
		assert.IsType(t, PipelineConsumerSamplerEmitter{}, samplerEmitter)
	})

	t.Run("UnknownOutputType", func(t *testing.T) {
		// Prepare mock data
		mockPersister := &MockPersister{
			Data: make(map[string][]byte),
		}

		// Set up some mock data
		mockPersister.Data["test.key"] = []byte("value")

		mockEmitter := &helper.LogEmitter{}

		mockInput := &file.Input{}

		// Call SamplerEmitterFactory
		samplerEmitter, err := SamplerEmitterFactory("unknown_output", "test.log", mockPersister, mockEmitter, *mockInput)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, samplerEmitter)
	})
}

func TestLogEntry(t *testing.T) {
	// Prepare mock data
	mockPersister := &MockPersister{
		Data: make(map[string][]byte),
	}

	// Set up some mock data
	mockPersister.Data["test.key"] = []byte("value")

	mockSampler := &mockSampler{}

	// Call logEntry function
	jsonEntry := logEntry(context.Background(), mockPersister, mockSampler)

	var logEntry LogEntry
	json.Unmarshal([]byte(jsonEntry), &logEntry)

	assert.Equal(t, "v1", logEntry.Format)
	assert.Equal(t, "100", strconv.FormatUint(logEntry.Events[0].UsageBytes, 10))
	assert.Equal(t, logsampler.NetworkSchemaId, logEntry.Metadata[logsampler.SchemaID])
}

// MockPersister is a mock implementation of Persister interface for testing.
type MockPersister struct {
	Data map[string][]byte // Store data for testing
}

// Get retrieves the value associated with the given key.
func (m *MockPersister) Get(ctx context.Context, key string) ([]byte, error) {
	if val, ok := m.Data[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("key %s not found", key)
}

// Set sets the value associated with the given key.
func (m *MockPersister) Set(ctx context.Context, key string, value []byte) error {
	m.Data[key] = value
	return nil
}

// Delete deletes the value associated with the given key.
func (m *MockPersister) Delete(ctx context.Context, key string) error {
	delete(m.Data, key)
	return nil
}
