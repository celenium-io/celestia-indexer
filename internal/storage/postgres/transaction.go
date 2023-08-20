package postgres

import (
	"context"

	models "github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// Transaction -
type Transaction struct {
	storage.Transaction
}

// BeginTransaction -
func BeginTransaction(ctx context.Context, tx storage.Transactable) (Transaction, error) {
	t, err := tx.BeginTransaction(ctx)
	return Transaction{t}, err
}

// SaveTransactions -
func (tx Transaction) SaveTransactions(ctx context.Context, txs ...models.Tx) error {
	switch len(txs) {
	case 0:
		return nil
	case 1:
		return tx.Add(ctx, &txs[0])
	default:
		arr := make([]any, len(txs))
		for i := range txs {
			arr[i] = &txs[i]
		}
		return tx.BulkSave(ctx, arr)
	}
}
