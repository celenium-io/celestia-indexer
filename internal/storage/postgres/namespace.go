package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Namespace -
type Namespace struct {
	*postgres.Table[*storage.Namespace]
}

// NewNamespace -
func NewNamespace(db *database.Bun) *Namespace {
	return &Namespace{
		Table: postgres.NewTable[*storage.Namespace](db),
	}
}

// ByNamespaceId -
func (n *Namespace) ByNamespaceId(ctx context.Context, namespaceId []byte) (namespace []storage.Namespace, err error) {
	err = n.DB().NewSelect().Model(&namespace).
		Where("namespace_id = ?", namespaceId).
		Scan(ctx)
	return
}

// ByNamespaceIdAndVersion -
func (n *Namespace) ByNamespaceIdAndVersion(ctx context.Context, namespaceId []byte, version byte) (namespace storage.Namespace, err error) {
	err = n.DB().NewSelect().Model(&namespace).
		Where("namespace_id = ?", namespaceId).
		Where("version = ?", version).
		Scan(ctx)
	return
}
