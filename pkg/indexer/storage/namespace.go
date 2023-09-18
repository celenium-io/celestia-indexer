package storage

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
)

func (module *Module) saveNamespaces(
	ctx context.Context,
	tx postgres.Transaction,
	namespaces map[string]*storage.Namespace,
) (uint64, error) {
	if len(namespaces) == 0 {
		return 0, nil
	}

	data := make([]*storage.Namespace, 0, len(namespaces))
	for key := range namespaces {
		data = append(data, namespaces[key])
	}

	totalNamespaces, err := tx.SaveNamespaces(ctx, data...)
	if err != nil {
		return 0, tx.HandleError(ctx, err)
	}

	return totalNamespaces, nil
}
