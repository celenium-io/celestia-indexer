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

func (tx Transaction) SaveNamespaces(ctx context.Context, namespaces ...*models.Namespace) error {
	if len(namespaces) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&namespaces).
		Column("version", "namespace_id", "pfb_count", "size", "first_height").
		On("CONFLICT ON CONSTRAINT namespace_id_version_idx DO UPDATE").
		Set("size = EXCLUDED.size + namespace.size").
		Set("pfb_count = EXCLUDED.pfb_count + namespace.pfb_count").
		Returning("id").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveAddresses(ctx context.Context, addresses ...*models.Address) error {
	if len(addresses) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&addresses).
		Column("height", "balance", "hash").
		On("CONFLICT ON CONSTRAINT address_hash DO UPDATE").
		Set("balance = EXCLUDED.balance + address.balance").
		Returning("id").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveTxAddresses(ctx context.Context, addresses ...models.TxAddress) error {
	if len(addresses) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&addresses).Exec(ctx)
	return err
}

func (tx Transaction) LastBlock(ctx context.Context) (block models.Block, err error) {
	err = tx.Tx().NewSelect().Model(&block).Order("id desc").Limit(1).Scan(ctx)
	return
}

func (tx Transaction) State(ctx context.Context, name string) (state models.State, err error) {
	err = tx.Tx().NewSelect().Model(&state).Where("name = ?", name).Scan(ctx)
	return
}

func (tx Transaction) Namespace(ctx context.Context, id uint64) (ns models.Namespace, err error) {
	err = tx.Tx().NewSelect().Model(&ns).Where("id = ?", id).Scan(ctx)
	return
}

func (tx Transaction) RollbackBlock(ctx context.Context, height models.Level) (block models.Block, err error) {
	_, err = tx.Tx().NewDelete().Model(&block).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackAddresses(ctx context.Context, height models.Level) (address []models.Address, err error) {
	_, err = tx.Tx().NewDelete().Model(&address).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackTxs(ctx context.Context, height models.Level) (txs []models.Tx, err error) {
	_, err = tx.Tx().NewDelete().Model(&txs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackEvents(ctx context.Context, height models.Level) (events []models.Event, err error) {
	_, err = tx.Tx().NewDelete().Model(&events).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackMessages(ctx context.Context, height models.Level) (msgs []models.Message, err error) {
	_, err = tx.Tx().NewDelete().Model(&msgs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackNamespaceMessages(ctx context.Context, height models.Level) (msgs []models.NamespaceMessage, err error) {
	_, err = tx.Tx().NewDelete().Model(&msgs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackNamespaces(ctx context.Context, height models.Level) (ns []models.Namespace, err error) {
	_, err = tx.Tx().NewDelete().Model(&ns).Where("first_height = ?", height).Returning("*").Exec(ctx)
	return
}
