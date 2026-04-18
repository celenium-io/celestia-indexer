// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"

	signalTypes "github.com/celestiaorg/celestia-app/v8/x/signal/types"
)

// MsgSignalVersion -
func MsgSignalVersion(ctx *context.Context, status storageTypes.Status, txId, msgId uint64, m *signalTypes.MsgSignalVersion) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSignalVersion
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height, msgId)
	if err != nil {
		return msgType, err
	}

	if status != storageTypes.StatusSuccess {
		return msgType, nil
	}

	validator := storage.EmptyValidator()
	validator.Address = m.ValidatorAddress
	validator.Version = m.Version

	signal := &storage.SignalVersion{
		Height:    ctx.Block.Height,
		Validator: &validator,
		Time:      ctx.Block.Time,
		Version:   m.Version,
		TxId:      txId,
		MsgId:     msgId,
	}
	ctx.AddSignal(signal)
	ctx.AddValidator(validator)
	ctx.AddUpgrade(storage.Upgrade{
		Version:      m.Version,
		SignalsCount: 1,
		Height:       ctx.Block.Height,
		Time:         ctx.Block.Time,
	})
	return msgType, err
}

// MsgTryUpgrade -
func MsgTryUpgrade(ctx *context.Context, status storageTypes.Status, txId, msgId uint64, m *signalTypes.MsgTryUpgrade) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgTryUpgrade
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	if err != nil {
		return msgType, err
	}

	if status != storageTypes.StatusSuccess {
		return msgType, nil
	}

	ctx.TryUpgrade = &storage.Upgrade{
		Height:    ctx.Block.Height,
		Time:      ctx.Block.Time,
		EndTime:   ctx.Block.Time,
		EndHeight: ctx.Block.Height,
		Signer: &storage.Address{
			Address: m.Signer,
		},
		TxId:  txId,
		MsgId: msgId,
	}
	return msgType, err
}
