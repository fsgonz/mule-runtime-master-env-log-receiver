package storage

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/matcher"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
	"sync"
	"time"
)

type StorageConsumerConfig struct {
	PollInterval time.Duration `mapstructure:"poll_interval,omitempty"`
}

func (c StorageConsumerConfig) Build(set component.TelemetrySettings, emit func(ctx context.Context, token []byte, attrs map[string]any) error) (*Manager, error) {
	return &Manager{
		set:          set,
		pollInterval: c.PollInterval,
	}, nil
}

type Manager struct {
	// Deprecated [v0.101.0]
	*zap.SugaredLogger

	set    component.TelemetrySettings
	wg     sync.WaitGroup
	cancel context.CancelFunc

	fileMatcher *matcher.Matcher

	pollInterval time.Duration
	emit         func(ctx context.Context, token []byte, attrs map[string]any) error
}

func (m *Manager) Start(persister operator.Persister) error {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	// Start polling goroutine
	m.startPoller(ctx)

	return nil
}

// Stop will stop the file monitoring process
func (m *Manager) Stop() error {
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
	m.wg.Wait()
	return nil
}

// startPoller kicks off a goroutine that will poll the filesystem periodically,
// checking if there are new files or new logs in the watched files
func (m *Manager) startPoller(ctx context.Context) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		globTicker := time.NewTicker(m.pollInterval)
		defer globTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-globTicker.C:
			}

			m.poll(ctx)
		}
	}()
}

// poll checks all the watched paths for new entries
func (m *Manager) poll(ctx context.Context) {
	m.emit(ctx, []byte("hola"), make(map[string]any))
}
