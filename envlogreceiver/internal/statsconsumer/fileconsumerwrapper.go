package statsconsumer

import (
	"context"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogreceiver/internal/lumberjack"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"log"
)

type FileConsumerWrapper struct {
	manager fileconsumer.Manager
	emitter Emitter
	path    string
}

func (m *FileConsumerWrapper) Start(persister operator.Persister) error {
	metricsLogger := log.New(&lumberjack.Logger{
		Filename:   m.path,
		MaxSize:    100, // kilobytes
		MaxBackups: 20,
	}, "", 0)

	m.emitter.emit = emitToFile(metricsLogger)
	m.emitter.Start(persister)
	m.manager.Start(persister)
	return nil
}

func emitToFile(logger *log.Logger) func(ctx context.Context, token []byte, attrs map[string]any) error {
	return func(ctx context.Context, token []byte, attrs map[string]any) error {
		logger.Println(string(token))
		return nil
	}
}

// Stop will stop the file monitoring process
func (m *FileConsumerWrapper) Stop() error {
	m.manager.Stop()
	return nil
}
