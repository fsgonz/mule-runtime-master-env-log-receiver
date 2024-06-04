package statsconsumer

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/matcher"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/split"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/trim"
	"go.uber.org/zap"
	"time"
)

const (
	defaultMaxConcurrentFiles = 1024
	defaultEncoding           = "utf-8"
	defaultPollInterval       = 200 * time.Millisecond
	defaultFingerprintSize    = 1000
	DefaultMaxLogSize         = 1024 * 1024
	DefaultFlushPeriod        = 500 * time.Millisecond
)

type FileConsumerConfig struct {
	Path               string                     `mapstructure:"path,omitempty"`
	PollInterval       time.Duration              `mapstructure:"poll_interval,omitempty"`
	StartAt            string                     `mapstructure:"start_at,omitempty"`
	FingerprintSize    helper.ByteSize            `mapstructure:"fingerprint_size,omitempty"`
	MaxLogSize         helper.ByteSize            `mapstructure:"max_log_size,omitempty"`
	MaxConcurrentFiles int                        `mapstructure:"max_concurrent_files,omitempty"`
	MaxBatches         int                        `mapstructure:"max_batches,omitempty"`
	DeleteAfterRead    bool                       `mapstructure:"delete_after_read,omitempty"`
	SplitConfig        split.Config               `mapstructure:"multiline,omitempty"`
	TrimConfig         trim.Config                `mapstructure:",squash,omitempty"`
	Encoding           string                     `mapstructure:"encoding,omitempty"`
	FlushPeriod        time.Duration              `mapstructure:"force_flush_period,omitempty"`
	Header             *fileconsumer.HeaderConfig `mapstructure:"header,omitempty"`
}

func NewConsumerConfig() *FileConsumerConfig {
	return &FileConsumerConfig{
		PollInterval:       defaultPollInterval,
		MaxConcurrentFiles: defaultMaxConcurrentFiles,
		StartAt:            "end",
		FingerprintSize:    defaultFingerprintSize,
		MaxLogSize:         DefaultMaxLogSize,
		Encoding:           defaultEncoding,
		FlushPeriod:        DefaultFlushPeriod,
	}
}

func (c FileConsumerConfig) Build(logger *zap.SugaredLogger, emit func(ctx context.Context, token []byte, attrs map[string]any) error) (*fileconsumer.Manager, error) {
	criteria := matcher.Criteria{Include: []string{c.Path}}

	config := fileconsumer.Config{
		Criteria:           criteria,
		PollInterval:       c.PollInterval,
		MaxConcurrentFiles: c.MaxConcurrentFiles,
		MaxBatches:         c.MaxBatches,
		StartAt:            c.StartAt,
		FingerprintSize:    c.FingerprintSize,
		MaxLogSize:         c.MaxLogSize,
		Encoding:           c.Encoding,
		SplitConfig:        c.SplitConfig,
		TrimConfig:         c.TrimConfig,
		FlushPeriod:        c.FlushPeriod,
		Header:             c.Header,
		DeleteAfterRead:    c.DeleteAfterRead,
	}

	return config.Build(logger, emit)
}
