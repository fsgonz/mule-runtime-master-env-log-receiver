package logsampler

import (
	"time"
)

// Config represents the configuration for log samplers.
type Config struct {
	LogSamplers []LogSampler `mapstructure:"log_samplers"`
}

// LogSamplerError represents an error that occurs during log sampling.
type LogSamplerError struct {
	Msg string
}

// Error returns the error message.
func (e *LogSamplerError) Error() string {
	return e.Msg
}

// LogSampler represents a log sampling configuration.
type LogSampler struct {
	Metric       string        `mapstructure:"metric"`
	PollInterval time.Duration `mapstructure:"poll_interval,omitempty"`
}

// Validate validates the configuration.
func (cfg *Config) Validate() error {
	if len(cfg.LogSamplers) > 1 {
		return &LogSamplerError{"No more than one sampler supported in this version"}
	}

	if len(cfg.LogSamplers) == 1 {
		logSampler := cfg.LogSamplers[0]

		if logSampler.Metric != MetricNetstats {
			return &LogSamplerError{"Incorrect metric in sampler. Possible Values: [" + MetricNetstats + "]"}
		}
		switch logSampler.Output {
		case OutputFileLogger, OutputPipelineEmitter:
			break
		default:
			return &LogSamplerError{"Incorrect output in sampler. Possible Values: [" + OutputFileLogger + ", " + OutputPipelineEmitter + "]"}
		}
	}
	return nil
}
