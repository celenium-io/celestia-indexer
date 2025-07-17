// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upChangeMsgTypes, downChangeMsgTypes)
}

var addedTypes = []types.MsgType{
	types.MsgCreateIgp,
	types.MsgSetIgpOwner,
	types.MsgSetDestinationGasConfig,
	types.MsgPayForGas,
	types.MsgClaim,
	types.MsgCreateMerkleTreeHook,
	types.MsgCreateNoopHook,
	types.MsgCreateMessageIdMultisigIsm,
	types.MsgCreateMerkleRootMultisigIsm,
	types.MsgCreateNoopIsm,
	types.MsgAnnounceValidator,
	types.MsgCreateRoutingIsm,
}

func upChangeMsgTypes(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE block
		ALTER COLUMN message_types
		TYPE bit varying(101)
		USING '000000000000'::bit varying(12) || message_types;`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE tx
		ALTER COLUMN message_types
		TYPE bit varying(101)
		USING '000000000000'::bit varying(12) || message_types;`); err != nil {
		return err
	}

	type msg struct {
		bun.BaseModel `bun:"message"`

		Height uint64        `bun:",notnull"`
		Type   types.MsgType `bun:",type:msg_type"`
		TxId   uint64        `bun:"tx_id"`
	}

	var msgs []msg
	if err := db.NewSelect().
		Column("tx_id", "height", "type").
		Table("message").
		Where("type IN (?)", bun.In(addedTypes)).
		Where("time > '2023-04-01T00:00:00Z'").
		Scan(ctx, &msgs); err != nil {
		return errors.Wrap(err, "get messages")
	}

	for i := range msgs {
		var block storage.Block
		if err := db.NewSelect().
			Model(&block).
			Column("id", "height", "message_types").
			Where("height = ?", msgs[i].Height).
			Scan(ctx); err != nil {
			return errors.Wrap(err, "get message block")
		}
		block.MessageTypes.SetByMsgType(msgs[i].Type)
		if _, err := db.NewUpdate().Model(&block).Where("id = ?", block.Id).Set("message_types = ?", block.MessageTypes).Exec(ctx); err != nil {
			return errors.Wrap(err, "get update block")
		}

		var tx storage.Tx
		if err := db.NewSelect().
			Model(&tx).
			Column("id", "message_types").
			Where("id = ?", msgs[i].TxId).
			Scan(ctx); err != nil {
			return errors.Wrap(err, "get message tx")
		}
		tx.MessageTypes.SetByMsgType(msgs[i].Type)
		if _, err := db.NewUpdate().Model(&tx).Where("id = ?", tx.Id).Set("message_types = ?", tx.MessageTypes).Exec(ctx); err != nil {
			return errors.Wrap(err, "get update tx")
		}
	}

	return nil
}
func downChangeMsgTypes(ctx context.Context, db *bun.DB) error {
	return nil
}
