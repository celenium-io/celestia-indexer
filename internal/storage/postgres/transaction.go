// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	pg "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

const copyThreshold = 20

func NewTxRepos(tx models.Transaction) models.TxRepos {
	db := tx.Tx()
	return models.TxRepos{
		Address:          NewAddress(db),
		Blocks:           NewBlocks(db),
		BondUpdates:      NewValidatorBondUpdate(db),
		Constants:        NewConstant(db),
		Delegation:       NewDelegation(db),
		HyperlaneMailbox: NewHLMailbox(db),
		HyperlaneToken:   NewHLToken(db),
		HyperlaneIgp:     NewHLIGP(db),
		IbcConnections:   NewIbcConnection(db),
		Namespace:        NewNamespace(db),
		Proposals:        NewProposal(db),
		SignalVersion:    NewSignalVersion(db),
		State:            NewState(db),
		Upgrades:         NewUpgrade(db),
		Validators:       NewValidator(db),
		Votes:            NewVote(db),
		ZkIsm:            NewZkISM(db),
	}
}

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

	_, err := tx.Tx().NewInsert().Model(&constants).
		Column("module", "name", "value").
		On("CONFLICT (module, name) DO UPDATE").
		Set("value = EXCLUDED.value").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveTransactions(ctx context.Context, txs ...models.Tx) error {
	return pg.SaveBulkWithCopy(ctx, tx, txs, copyThreshold)
}

func (tx Transaction) SaveEvents(ctx context.Context, events ...models.Event) error {
	return pg.SaveBulkWithCopy(ctx, tx, events, copyThreshold)
}

func (tx Transaction) SaveMessages(ctx context.Context, msgs ...*models.Message) error {
	return pg.SaveBulkWithCopy(ctx, tx, msgs, copyThreshold)
}

func (tx Transaction) SaveMsgAddresses(ctx context.Context, addresses ...*models.MsgAddress) error {
	return pg.SaveBulkWithCopy(ctx, tx, addresses, copyThreshold)
}

func (tx Transaction) Insert(ctx context.Context, data any) error {
	_, err := tx.Tx().NewInsert().Model(data).Exec(ctx)
	return err
}

func (tx Transaction) rollback[T storage.Model](ctx context.Context, height types.Level, model T) (int64, error) {
	result, err := tx.Tx().NewDelete().Model(model).
		Where("height = ?", height).
		Exec(ctx)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (tx Transaction) RollbackByHeight(ctx context.Context, height types.Level, models ...storage.Model) error {
	for i := range models {
		if _, err := tx.rollback(ctx, height, models[i]); err != nil {
			return err
		}
	}
	return nil
}

func (tx Transaction) RollbackBlockStats(ctx context.Context, height types.Level) (stats models.BlockStats, err error) {
	_, err = tx.Tx().NewDelete().Model(&stats).Where("height = ?", height).Returning("*").Exec(ctx)
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

func (tx Transaction) RollbackSigners(ctx context.Context, txIds []uint64) (err error) {
	_, err = tx.Tx().NewDelete().
		Model((*models.Signer)(nil)).
		Where("tx_id = ANY(?)", pgdialect.Array(txIds)).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackMessageAddresses(ctx context.Context, msgIds []uint64) (err error) {
	_, err = tx.Tx().NewDelete().
		Model((*models.MsgAddress)(nil)).
		Where("msg_id IN ?", bun.Tuple(msgIds)).
		Exec(ctx)
	return
}

func (tx Transaction) RetentionBlockSignatures(ctx context.Context, height types.Level) error {
	_, err := tx.Tx().NewDelete().Model((*models.BlockSignature)(nil)).
		Where("height <= ?", height).
		Exec(ctx)
	return err
}
