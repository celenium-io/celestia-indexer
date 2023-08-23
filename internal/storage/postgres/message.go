package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Message -
type Message struct {
	*postgres.Table[*storage.Message]
}

// NewMessage -
func NewMessage(db *database.Bun) *Message {
	return &Message{
		Table: postgres.NewTable[*storage.Message](db),
	}
}

// ByTxId -
func (m *Message) ByTxId(ctx context.Context, txId uint64) (messages []storage.Message, err error) {
	err = m.DB().NewSelect().Model(&messages).
		Where("tx_id = ?", txId).
		Scan(ctx)
	return
}
