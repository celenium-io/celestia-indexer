package postgres

import (
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
