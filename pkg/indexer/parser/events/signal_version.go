// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func processSignalVersion(ctx *context.Context, c *Cursor, msg *storage.Message, data map[string]any) error {
	version, err := decoder.Uint64(data, "Version")
	if err != nil {
		return errors.Wrap(err, "get signal version in exec")
	}

	val := storage.EmptyValidator()
	val.Address, err = (types.PackedBytes)(data).GetString("ValidatorAddress")
	if err != nil {
		return err
	}

	val.Version = version

	signalVersion := &storage.SignalVersion{
		Height:    msg.Height,
		Time:      msg.Time,
		Version:   version,
		Validator: &val,
		TxId:      msg.TxId,
		MsgId:     msg.Id,
	}
	ctx.AddSignal(signalVersion)
	ctx.AddValidator(*signalVersion.Validator)
	ctx.AddUpgrade(storage.Upgrade{
		Version:      version,
		SignalsCount: 1,
		Height:       msg.Height,
		Time:         msg.Time,
	})
	c.Skip(1)
	return nil
}

func processTryUpgrade(ctx *context.Context, c *Cursor, msg *storage.Message, data map[string]any) error {
	signer, err := (types.PackedBytes)(data).GetString("Signer")
	if err != nil {
		return errors.Wrap(err, "get try upgrade signer in exec")
	}

	// the signer of an inner message is the granter, which MsgExec does not register
	address := &storage.Address{
		Address:    signer,
		Height:     msg.Height,
		LastHeight: msg.Height,
	}
	if err := ctx.AddAddress(address); err != nil {
		return err
	}

	ctx.TryUpgrade = &storage.Upgrade{
		Height:    msg.Height,
		Time:      msg.Time,
		EndTime:   msg.Time,
		EndHeight: msg.Height,
		Signer:    address,
		TxId:      msg.TxId,
		MsgId:     msg.Id,
	}
	c.Skip(1)
	return nil
}
