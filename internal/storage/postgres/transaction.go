package postgres

import (
	"context"

	models "github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

type Transaction struct {
	storage.Transaction
}

func BeginTransaction(ctx context.Context, tx storage.Transactable) (Transaction, error) {
	t, err := tx.BeginTransaction(ctx)
	return Transaction{t}, err
}

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

func (tx Transaction) SaveNamespaces(ctx context.Context, namespaces ...models.Namespace) error {
	if len(namespaces) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&namespaces).
		Column("version", "namespace_id", "pfd_count", "size").
		On("CONFLICT ON CONSTRAINT namespace_id_version_idx DO UPDATE").
		Set("size = EXCLUDED.size + namespace.size").
		Set("pfd_count = EXCLUDED.pfd_count + namespace.pfd_count").
		Returning("id").
		Exec(ctx)
	return err
}
