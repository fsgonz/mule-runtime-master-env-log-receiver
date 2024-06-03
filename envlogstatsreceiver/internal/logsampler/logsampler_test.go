package logsampler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	t.Run("Empty configuration", func(t *testing.T) {
		cfg := &Config{}
		err := cfg.Validate()
		assert.NoError(t, err, "Empty configuration should pass validation")
	})

	t.Run("Valid configuration", func(t *testing.T) {
		cfg := &Config{
			LogSamplers: []LogSampler{
				{
					Metric: MetricNetstats,
				},
			},
		}
		err := cfg.Validate()
		assert.NoError(t, err, "Valid configuration should pass validation")
	})

	t.Run("Multiple log samplers", func(t *testing.T) {
		cfg := &Config{
			LogSamplers: []LogSampler{
				{
					Metric: MetricNetstats,
				},
				{
					Metric: MetricNetstats,
				},
			},
		}
		err := cfg.Validate()
		assert.Error(t, err, "Multiple log samplers should fail validation")
	})

	t.Run("Invalid metric", func(t *testing.T) {
		cfg := &Config{
			LogSamplers: []LogSampler{
				{
					Metric: "invalid_metric",
				},
			},
		}
		err := cfg.Validate()
		assert.Error(t, err, "Invalid metric should fail validation")
	})
}
