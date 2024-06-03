package statsconsumer

import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"

type StartStoppable interface {
	Start(persister operator.Persister) error
	Stop() error
}
