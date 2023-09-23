package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/uptrace/bun"

	models "github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

type Transaction struct {
	storage.Transaction
}

func BeginTransaction(ctx context.Context, tx storage.Transactable) (models.Transaction, error) {
	t, err := tx.BeginTransaction(ctx)
	return Transaction{t}, err
}

func (tx Transaction) SaveConstants(ctx context.Context, constants ...models.Constant) error {
	if len(constants) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&constants).Exec(ctx)
	return err
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

type addedNamespace struct {
	bun.BaseModel `bun:"namespace"`
	*models.Namespace

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveNamespaces(ctx context.Context, namespaces ...*models.Namespace) (uint64, error) {
	if len(namespaces) == 0 {
		return 0, nil
	}

	addedNamespaces := make([]addedNamespace, len(namespaces))
	for i := range namespaces {
		addedNamespaces[i].Namespace = namespaces[i]
	}

	_, err := tx.Tx().NewInsert().Model(&addedNamespaces).
		Column("version", "namespace_id", "pfb_count", "size", "first_height").
		On("CONFLICT ON CONSTRAINT namespace_id_version_idx DO UPDATE").
		Set("size = EXCLUDED.size + added_namespace.size").
		Set("pfb_count = EXCLUDED.pfb_count + added_namespace.pfb_count").
		Returning("id").
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	var count uint64
	for i := range addedNamespaces {
		if addedNamespaces[i].Xmax == 0 {
			count++
		}
	}

	return count, err
}

type addedAddress struct {
	bun.BaseModel `bun:"address"`
	*models.Address

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveAddresses(ctx context.Context, addresses ...*models.Address) (uint64, error) {
	if len(addresses) == 0 {
		return 0, nil
	}

	addr := make([]addedAddress, len(addresses))
	for i := range addresses {
		addr[i].Address = addresses[i]
	}

	_, err := tx.Tx().NewInsert().Model(&addr).
		Column("address", "height", "last_height", "hash").
		On("CONFLICT ON CONSTRAINT address_idx DO UPDATE").
		Set("last_height = EXCLUDED.last_height").
		Returning("xmax, id").
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	var count uint64
	for i := range addr {
		if addr[i].Xmax == 0 {
			count++
		}
	}

	return count, err
}

func (tx Transaction) SaveBalances(ctx context.Context, balances ...models.Balance) error {
	if len(balances) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&balances).
		Column("id", "currency", "total").
		On("CONFLICT (id, currency) DO UPDATE").
		Set("total = EXCLUDED.total + balance.total").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveEvents(ctx context.Context, events ...models.Event) error {
	switch {
	case len(events) == 0:
		return nil
	case len(events) < 20:
		data := make([]any, len(events))
		for i := range events {
			data[i] = &events[i]
		}
		return tx.BulkSave(ctx, data)
	default:
		copiable := make([]storage.Copiable, len(events))
		for i := range events {
			copiable[i] = events[i]
		}
		return tx.CopyFrom(ctx, "event", copiable)
	}
}

func (tx Transaction) SaveMessages(ctx context.Context, msgs ...*models.Message) error {
	if len(msgs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&msgs).Returning("id").Exec(ctx)
	return err
}

func (tx Transaction) SaveSigners(ctx context.Context, addresses ...models.Signer) error {
	if len(addresses) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&addresses).Exec(ctx)
	return err
}

func (tx Transaction) SaveMsgAddresses(ctx context.Context, addresses ...models.MsgAddress) error {
	if len(addresses) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&addresses).Exec(ctx)
	return err
}

func (tx Transaction) SaveNamespaceMessage(ctx context.Context, nsMsgs ...models.NamespaceMessage) error {
	if len(nsMsgs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&nsMsgs).Exec(ctx)
	return err
}

func (tx Transaction) SaveValidators(ctx context.Context, validators ...*models.Validator) error {
	if len(validators) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&validators).
		On("CONFLICT ON CONSTRAINT address_validator DO UPDATE").
		Set("moniker = EXCLUDED.moniker").
		Set("website = EXCLUDED.website").
		Set("identity = EXCLUDED.identity").
		Set("contacts = EXCLUDED.contacts").
		Set("details = EXCLUDED.details").
		Set("rate = EXCLUDED.rate").
		Set("min_self_delegation = EXCLUDED.min_self_delegation").
		Returning("id").
		Exec(ctx)
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

func (tx Transaction) RollbackBlock(ctx context.Context, height types.Level) error {
	_, err := tx.Tx().NewDelete().
		Model((*models.Block)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return err
}

func (tx Transaction) RollbackBlockStats(ctx context.Context, height types.Level) (stats models.BlockStats, err error) {
	_, err = tx.Tx().NewDelete().Model(&stats).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackAddresses(ctx context.Context, height types.Level) (address []models.Address, err error) {
	_, err = tx.Tx().NewDelete().Model(&address).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackTxs(ctx context.Context, height types.Level) (txs []models.Tx, err error) {
	_, err = tx.Tx().NewDelete().Model(&txs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackEvents(ctx context.Context, height types.Level) (events []models.Event, err error) {
	_, err = tx.Tx().NewDelete().Model(&events).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackMessages(ctx context.Context, height types.Level) (msgs []models.Message, err error) {
	_, err = tx.Tx().NewDelete().Model(&msgs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackNamespaceMessages(ctx context.Context, height types.Level) (msgs []models.NamespaceMessage, err error) {
	_, err = tx.Tx().NewDelete().Model(&msgs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackNamespaces(ctx context.Context, height types.Level) (ns []models.Namespace, err error) {
	_, err = tx.Tx().NewDelete().Model(&ns).Where("first_height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackValidators(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.Validator)(nil)).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackSigners(ctx context.Context, txIds []uint64) (err error) {
	_, err = tx.Tx().NewDelete().
		Model((*models.Signer)(nil)).
		Where("tx_id IN (?)", bun.In(txIds)).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackMessageAddresses(ctx context.Context, msgIds []uint64) (err error) {
	_, err = tx.Tx().NewDelete().
		Model((*models.MsgAddress)(nil)).
		Where("msg_id IN (?)", bun.In(msgIds)).
		Exec(ctx)
	return
}

func (tx Transaction) DeleteBalances(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := tx.Tx().NewDelete().
		Model((*models.Balance)(nil)).
		Where("id IN (?)", bun.In(ids)).
		Exec(ctx)
	return err
}

func (tx Transaction) LastAddressAction(ctx context.Context, address []byte) (uint64, error) {
	var height uint64
	err := tx.Tx().NewSelect().
		Model((*models.MsgAddress)(nil)).
		ExcludeColumn("msg_id", "address_id", "type").
		Where("address.hash = ?", address).
		Order("msg_id desc").
		Relation("Msg", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("height")
		}).
		Relation("Address", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ExcludeColumn("*")
		}).
		Scan(ctx, &height)
	return height, err
}
