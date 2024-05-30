package file

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/attrs"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/matcher"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/split"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/trim"
	"go.opentelemetry.io/collector/component"
	"time"
)

type ConsumerConfig struct {
	Path               string `mapstructure:"path,omitempty"`
	attrs.Resolver     `mapstructure:",squash"`
	PollInterval       time.Duration              `mapstructure:"poll_interval,omitempty"`
	MaxConcurrentFiles int                        `mapstructure:"max_concurrent_files,omitempty"`
	MaxBatches         int                        `mapstructure:"max_batches,omitempty"`
	StartAt            string                     `mapstructure:"start_at,omitempty"`
	FingerprintSize    helper.ByteSize            `mapstructure:"fingerprint_size,omitempty"`
	MaxLogSize         helper.ByteSize            `mapstructure:"max_log_size,omitempty"`
	Encoding           string                     `mapstructure:"encoding,omitempty"`
	SplitConfig        split.Config               `mapstructure:"multiline,omitempty"`
	TrimConfig         trim.Config                `mapstructure:",squash,omitempty"`
	FlushPeriod        time.Duration              `mapstructure:"force_flush_period,omitempty"`
	Header             *fileconsumer.HeaderConfig `mapstructure:"header,omitempty"`
	DeleteAfterRead    bool                       `mapstructure:"delete_after_read,omitempty"`
}

func (c ConsumerConfig) Build(set component.TelemetrySettings, emit func(ctx context.Context, token []byte, attrs map[string]any) error) (*fileconsumer.Manager, error) {
	criteria := matcher.Criteria{Include: []string{c.Path}}

	config := fileconsumer.Config{
		criteria,
		c.Resolver,
		c.PollInterval,
		c.MaxConcurrentFiles,
		c.MaxBatches,
		c.StartAt,
		c.FingerprintSize,
		c.MaxLogSize,
		c.Encoding,
		c.SplitConfig,
		c.TrimConfig,
		c.FlushPeriod,
		c.Header,
		c.DeleteAfterRead,
	}

	return config.Build(set, emit)
}
