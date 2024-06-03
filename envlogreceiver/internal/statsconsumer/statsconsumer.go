package statsconsumer

import (
	"context"
	"encoding/json"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/file"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/logsampler"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/stats/sampler"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/stats/scraper"
	"github.com/google/uuid"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/matcher"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type networkIOLogEntry struct {
	// Format is the schema version
	Format string `json:"format"`
	// Time is the time this entry was created in unix epoch milliseconds
	Time     int64                    `json:"time"`
	Events   []networkIOLogEntryEvent `json:"events"`
	Metadata map[string]string        `json:"metadata"`
}

type networkIOLogEntryEvent struct {
	ID string `json:"id"`
	// Timestamp is the time this entry was created in unix epoch milliseconds
	Timestamp  int64  `json:"timestamp"`
	RootOrgID  string `json:"root_org_id"`
	OrgID      string `json:"org_id"`
	EnvID      string `json:"env_id"`
	AssetID    string `json:"asset_id"`
	WorkerID   string `json:"worker_id"`
	UsageBytes uint64 `json:"usage_bytes"`
	Billable   bool   `json:"billable"`
}

type StatsConsumerConfig struct {
	PollInterval time.Duration `mapstructure:"poll_interval,omitempty"`
}

func Build(set component.TelemetrySettings, logSampler logsampler.LogSampler, emit func(ctx context.Context, Logtoken []byte, attrs map[string]any) error, fileConsumerConfig file.FileConsumerConfig) (file.StartStoppable, error) {
	if fileConsumerConfig.Path != "" {
		criteria := matcher.Criteria{Include: []string{fileConsumerConfig.Path}}

		config := fileconsumer.Config{
			criteria,
			fileConsumerConfig.Resolver,
			fileConsumerConfig.PollInterval,
			fileConsumerConfig.MaxConcurrentFiles,
			fileConsumerConfig.MaxBatches,
			fileConsumerConfig.StartAt,
			fileConsumerConfig.FingerprintSize,
			fileConsumerConfig.MaxLogSize,
			fileConsumerConfig.Encoding,
			fileConsumerConfig.SplitConfig,
			fileConsumerConfig.TrimConfig,
			fileConsumerConfig.FlushPeriod,
			fileConsumerConfig.Header,
			fileConsumerConfig.DeleteAfterRead,
		}

		return config.Build(set, emit)
	}

	return &Manager{
		set:          set,
		pollInterval: logSampler.PollInterval,
		emit:         emit,
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
	persister    operator.Persister
	emit         func(ctx context.Context, token []byte, attrs map[string]any) error
}

func (m *Manager) Start(persister operator.Persister) error {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.persister = persister

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
	byteSlice, _ := m.persister.Get(ctx, logsampler.LastCountKey)

	sampler := sampler.NewFileBasedSampler("/proc/net/dev", scraper.NewLinuxNetworkDevicesFileScraper())
	var last_count uint64 = 0

	if byteSlice != nil {
		// Parse the string to an integer
		counter, _ := strconv.ParseUint(string(byteSlice), 10, 64)
		last_count = counter
	}

	samp, _ := sampler.Sample()

	m.persister.Set(ctx, logsampler.LastCountKey, []byte(strconv.FormatUint(samp, 10)))

	orgID := os.Getenv(logsampler.OrgID)
	envID := os.Getenv(logsampler.EnvID)
	deploymentID := os.Getenv(logsampler.DeploymentID)
	rootOrgID := os.Getenv(logsampler.RootOrgID)
	billingEnabled := os.Getenv(logsampler.MuleBillingEnabled) == "true"
	workerID := "worker-" + strings.ReplaceAll(os.Getenv(logsampler.PodName), os.Getenv(logsampler.AppName)+"-", "")
	ts := time.Now().Unix() * 1000

	u, _ := uuid.NewRandom()

	evt := networkIOLogEntryEvent{
		ID:         u.String(),
		Timestamp:  ts,
		RootOrgID:  rootOrgID,
		OrgID:      orgID,
		EnvID:      envID,
		AssetID:    deploymentID,
		WorkerID:   workerID,
		UsageBytes: samp - last_count,
		Billable:   billingEnabled,
	}

	logEntry := networkIOLogEntry{
		Format: logsampler.Format,
		Time:   ts,
		Events: []networkIOLogEntryEvent{evt},
		Metadata: map[string]string{
			logsampler.SchemaID: logsampler.NetworkSchemaId,
		},
	}

	jsonEntry, _ := json.Marshal(logEntry)

	m.emit(ctx, jsonEntry, make(map[string]any))
}
