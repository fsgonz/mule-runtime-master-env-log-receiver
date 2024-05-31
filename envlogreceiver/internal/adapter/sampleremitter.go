package adapter

import (
	"context"
	"encoding/json"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/logsampler"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/lumberjack"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/stats/sampler"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/stats/scraper"
	"github.com/google/uuid"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"log"
	"os"
	"strconv"
	"strings"
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

type SamplerEmitter interface {
	Emit(context.Context)
}

type FileLoggerSamplerEmitter struct {
	metricsLogger log.Logger
	persister     operator.Persister
	sampler       sampler.Sampler
}

func (e FileLoggerSamplerEmitter) Emit(ctx context.Context) {
	jsonEntry := logEntry(ctx, e.persister, e.sampler)
	e.metricsLogger.Println(string(jsonEntry))

}

func SamplerEmitterFactory(uri string, persister operator.Persister, emitter *helper.LogEmitter, input helper.WriterOperator) (SamplerEmitter, error) {
	fileBasedSampler := sampler.NewFileBasedSampler("/proc/net/dev", scraper.NewLinuxNetworkDevicesFileScraper())

	metricsLogger := log.New(&lumberjack.Logger{
		Filename:   uri,
		MaxSize:    100, // kilobytes
		MaxBackups: 20,
	}, "", 0)

	return FileLoggerSamplerEmitter{
		*metricsLogger,
		persister,
		fileBasedSampler,
	}, nil
}

func logEntry(ctx context.Context, persister operator.Persister, sampler sampler.Sampler) []byte {
	byteSlice, _ := persister.Get(ctx, logsampler.LastCountKey)

	var last_count uint64 = 0

	if byteSlice != nil {
		// Parse the string to an integer
		counter, _ := strconv.ParseUint(string(byteSlice), 10, 64)
		last_count = counter
	}

	samp, _ := sampler.Sample()

	persister.Set(ctx, logsampler.LastCountKey, []byte(strconv.FormatUint(samp, 10)))

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
	return jsonEntry
}
