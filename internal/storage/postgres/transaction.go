// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/uptrace/bun"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
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

func (tx Transaction) SaveNamespaces(ctx context.Context, namespaces ...*models.Namespace) (int64, error) {
	if len(namespaces) == 0 {
		return 0, nil
	}

	addedNamespaces := make([]addedNamespace, len(namespaces))
	for i := range namespaces {
		addedNamespaces[i].Namespace = namespaces[i]
	}

	_, err := tx.Tx().NewInsert().Model(&addedNamespaces).
		Column("version", "namespace_id", "pfb_count", "size", "first_height", "last_height", "last_message_time", "blobs_count").
		On("CONFLICT ON CONSTRAINT namespace_id_version_idx DO UPDATE").
		Set("size = EXCLUDED.size + added_namespace.size").
		Set("pfb_count = EXCLUDED.pfb_count + added_namespace.pfb_count").
		Set("last_height = EXCLUDED.last_height").
		Set("last_message_time = EXCLUDED.last_message_time").
		Set("blobs_count = EXCLUDED.blobs_count + added_namespace.blobs_count").
		Returning("xmax, id").
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
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

func (tx Transaction) SaveAddresses(ctx context.Context, addresses ...*models.Address) (int64, error) {
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

	var count int64
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

func (tx Transaction) SaveBlobLogs(ctx context.Context, logs ...models.BlobLog) error {
	if len(logs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&logs).Exec(ctx)
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

const doNotModify = "[do-not-modify]"

func (tx Transaction) SaveValidators(ctx context.Context, validators ...*models.Validator) error {
	if len(validators) == 0 {
		return nil
	}

	for i := range validators {
		query := tx.Tx().NewInsert().Model(validators[i]).
			On("CONFLICT ON CONSTRAINT address_validator DO UPDATE").
			Set("rate = EXCLUDED.rate").
			Set("min_self_delegation = EXCLUDED.min_self_delegation")

		if validators[i].Moniker != doNotModify {
			query.Set("moniker = EXCLUDED.moniker")
		}
		if validators[i].Website != doNotModify {
			query.Set("website = EXCLUDED.website")
		}
		if validators[i].Identity != doNotModify {
			query.Set("identity = EXCLUDED.identity")
		}
		if validators[i].Contacts != doNotModify {
			query.Set("contacts = EXCLUDED.contacts")
		}
		if validators[i].Details != doNotModify {
			query.Set("details = EXCLUDED.details")
		}
		if _, err := query.Returning("id").Exec(ctx); err != nil {
			return err
		}
	}

	return nil
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
	_, err = tx.Tx().NewDelete().Model((*models.Validator)(nil)).Where("height = ?", height).Exec(ctx)
	return
}

func (tx Transaction) RollbackBlobLog(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.BlobLog)(nil)).Where("height = ?", height).Exec(ctx)
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

func (tx Transaction) LastNamespaceMessage(ctx context.Context, nsId uint64) (msg models.NamespaceMessage, err error) {
	err = tx.Tx().NewSelect().
		Model(&msg).
		Where("namespace_id = ?", nsId).
		Order("msg_id desc").
		Limit(1).
		Scan(ctx)
	return
}

func (tx Transaction) GetProposerId(ctx context.Context, address string) (id uint64, err error) {
	err = tx.Tx().NewSelect().
		Model((*models.Validator)(nil)).
		Column("id").
		Where("cons_address = ?", address).
		Order("msg_id desc").
		Limit(1).
		Scan(ctx, &id)
	return
}

func (tx Transaction) SaveRollup(ctx context.Context, rollup *models.Rollup) error {
	if rollup == nil {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(rollup).Exec(ctx)
	return err
}

func (tx Transaction) UpdateRollup(ctx context.Context, rollup *models.Rollup) error {
	if rollup == nil || rollup.IsEmpty() {
		return nil
	}

	query := tx.Tx().NewUpdate().Model(rollup).WherePK()

	if rollup.Name != "" {
		query = query.Set("name = ?", rollup.Name)
	}
	if rollup.Description != "" {
		query = query.Set("description = ?", rollup.Description)
	}
	if rollup.Twitter != "" {
		query = query.Set("twitter = ?", rollup.Twitter)
	}
	if rollup.GitHub != "" {
		query = query.Set("github = ?", rollup.GitHub)
	}
	if rollup.Website != "" {
		query = query.Set("website = ?", rollup.Website)
	}
	if rollup.Logo != "" {
		query = query.Set("logo = ?", rollup.Logo)
	}

	_, err := query.Exec(ctx)
	return err
}

func (tx Transaction) SaveProviders(ctx context.Context, providers ...models.RollupProvider) error {
	if len(providers) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&providers).Exec(ctx)
	return err
}

func (tx Transaction) DeleteProviders(ctx context.Context, rollupId uint64) error {
	if rollupId == 0 {
		return nil
	}
	_, err := tx.Tx().NewDelete().
		Model((*models.RollupProvider)(nil)).
		Where("rollup_id = ?", rollupId).
		Exec(ctx)
	return err
}

func (tx Transaction) DeleteRollup(ctx context.Context, rollupId uint64) error {
	if rollupId == 0 {
		return nil
	}
	_, err := tx.Tx().NewDelete().
		Model((*models.Rollup)(nil)).
		Where("id = ?", rollupId).
		Exec(ctx)
	return err
}
