// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddMissingHyperlaneMsgTypes, downAddMissingHyperlaneMsgTypes)
}

func upAddMissingHyperlaneMsgTypes(ctx context.Context, db *bun.DB) error {

	log.Info().Msg("add missing types...")

	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgSetRoutingIsmDomain.String(), types.MsgCreateRoutingIsm.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgRemoveRoutingIsmDomain.String(), types.MsgSetRoutingIsmDomain.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgUpdateRoutingIsmOwner.String(), types.MsgRemoveRoutingIsmDomain.String()); err != nil {
		return err
	}

	var addedTypes = []types.MsgType{
		types.MsgSetRoutingIsmDomain,
		types.MsgRemoveRoutingIsmDomain,
		types.MsgUpdateRoutingIsmOwner,
	}

	log.Info().Msg("migrate block mask column...")

	if _, err := db.ExecContext(ctx, `ALTER TABLE block
		ALTER COLUMN message_types
		TYPE bit varying(104)
		USING '000'::bit varying(12) || message_types;`); err != nil {
		return err
	}

	log.Info().Msg("migrate block tx column...")

	if _, err := db.ExecContext(ctx, `ALTER TABLE tx
		ALTER COLUMN message_types
		TYPE bit varying(104)
		USING '000'::bit varying(12) || message_types;`); err != nil {
		return err
	}

	type msg struct {
		bun.BaseModel `bun:"message"`

		Height uint64        `bun:",notnull"`
		Type   types.MsgType `bun:",type:msg_type"`
		TxId   uint64        `bun:"tx_id"`
	}

	log.Info().Msg("find messages...")
	var msgs []msg
	if err := db.NewSelect().
		Column("tx_id", "height", "type").
		Table("message").
		Where("type IN (?)", bun.In(addedTypes)).
		Where("time > '2023-04-01T00:00:00Z'").
		Scan(ctx, &msgs); err != nil {
		return errors.Wrap(err, "get messages")
	}

	log.Info().Msg("migrate tx and block masks for found messages...")
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

	log.Info().Msg("finished")
	return nil
}
func downAddMissingHyperlaneMsgTypes(ctx context.Context, db *bun.DB) error {
	return nil
}
