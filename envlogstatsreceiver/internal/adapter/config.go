// Taken in part from "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
// so that the connector offers the same functionality as filelogreceiver related to retry on failure
package adapter

import (
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogstatsreceiver/internal/consumerretry"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"go.opentelemetry.io/collector/component"
)

// BaseConfig is the common configuration of a stanza-based receiver
type BaseConfig struct {
	Operators      []operator.Config    `mapstructure:"operators"`
	StorageID      *component.ID        `mapstructure:"storage"`
	RetryOnFailure consumerretry.Config `mapstructure:"retry_on_failure"`

	// currently not configurable by users, but available for benchmarking
	numWorkers    int
	maxBatchSize  uint
	flushInterval time.Duration
}
