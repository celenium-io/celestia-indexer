package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Tx -
type Tx struct {
	*postgres.Table[*storage.Tx]
}

// NewTx -
func NewTx(db *database.Bun) *Tx {
	return &Tx{
		Table: postgres.NewTable[*storage.Tx](db),
	}
}

func (tx *Tx) ByHash(ctx context.Context, hash []byte) (transaction storage.Tx, err error) {
	err = tx.DB().NewSelect().Model(&transaction).
		Where("hash = ?", hash).
		Scan(ctx)
	return
}
