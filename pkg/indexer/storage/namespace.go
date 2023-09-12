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
) error {
	if len(namespaces) == 0 {
		return nil
	}

	data := make([]*storage.Namespace, 0, len(namespaces))
	for key := range namespaces {
		data = append(data, namespaces[key])
	}

	if err := tx.SaveNamespaces(ctx, data...); err != nil {
		return tx.HandleError(ctx, err)
	}

	return nil
}
