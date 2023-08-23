package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Event -
type Event struct {
	*postgres.Table[*storage.Event]
}

// NewEvent -
func NewEvent(db *database.Bun) *Event {
	return &Event{
		Table: postgres.NewTable[*storage.Event](db),
	}
}

// ByTxId -
func (e *Event) ByTxId(ctx context.Context, txId uint64) (events []storage.Event, err error) {
	err = e.DB().NewSelect().Model(&events).
		Where("tx_id = ?", txId).
		Scan(ctx)
	return
}

// ByBlock -
func (e *Event) ByBlock(ctx context.Context, height uint64) (events []storage.Event, err error) {
	err = e.DB().NewSelect().Model(&events).
		Where("height = ?", height).
		Where("tx_id IS NULL").
		Scan(ctx)
	return
}
