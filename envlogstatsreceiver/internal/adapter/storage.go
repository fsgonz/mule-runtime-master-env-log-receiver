// Taken in part from "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
// so that the connector offers the same functionality as filelogreceiver related to retry on failure

package adapter

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/experimental/storage"
)

func GetStorageClient(ctx context.Context, host component.Host, storageID *component.ID, componentID component.ID) (storage.Client, error) {
	if storageID == nil {
		return nil, fmt.Errorf("An storage is mandatory for stats for component", componentID)
	}

	extension, ok := host.GetExtensions()[*storageID]
	if !ok {
		return nil, fmt.Errorf("storage extension '%s' not found", storageID)
	}

	storageExtension, ok := extension.(storage.Extension)
	if !ok {
		return nil, fmt.Errorf("non-storage extension '%s' found", storageID)
	}

	return storageExtension.GetClient(ctx, component.KindReceiver, componentID, "")

}

func (r *receiver) setStorageClient(ctx context.Context, host component.Host) error {
	client, err := GetStorageClient(ctx, host, r.storageID, r.id)
	if err != nil {
		return err
	}
	r.storageClient = client
	return nil
}
