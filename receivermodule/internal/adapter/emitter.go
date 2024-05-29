// Taken in part from "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
// so that the connector offers the same functionality as filelogreceiver related to retry on failure

package adapter

import (
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
)

// Deprecated [v0.101.0] Use helper.LogEmitter directly instead
type LogEmitter = helper.LogEmitter

// Deprecated [v0.101.0] Use helper.NewLogEmitter directly instead
func NewLogEmitter(logger *zap.SugaredLogger, opts ...helper.EmitterOption) *LogEmitter {
	return helper.NewLogEmitter(component.TelemetrySettings{Logger: logger.Desugar()}, opts...)
}
