// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"

	signalTypes "github.com/celestiaorg/celestia-app/v7/x/signal/types"
)

// MsgSignalVersion -
func MsgSignalVersion(ctx *context.Context, status storageTypes.Status, m *signalTypes.MsgSignalVersion) (storageTypes.MsgType, []storage.AddressWithType, *storage.SignalVersion, error) {
	msgType := storageTypes.MsgSignalVersion
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height)
	if err != nil {
		return msgType, addresses, nil, err
	}

	if status != storageTypes.StatusSuccess {
		return msgType, addresses, nil, nil
	}

	validator := storage.EmptyValidator()
	validator.Address = m.ValidatorAddress
	validator.Version = m.Version

	signal := &storage.SignalVersion{
		Height:    ctx.Block.Height,
		Validator: &validator,
		Time:      ctx.Block.Time,
		Version:   m.Version,
	}
	ctx.AddValidator(*signal.Validator)
	ctx.AddUpgrade(storage.Upgrade{
		Version:      m.Version,
		SignalsCount: 1,
		Height:       ctx.Block.Height,
		Time:         ctx.Block.Time,
	})
	return msgType, addresses, signal, err
}

// MsgTryUpgrade -
func MsgTryUpgrade(ctx *context.Context, status storageTypes.Status, m *signalTypes.MsgTryUpgrade) (storageTypes.MsgType, []storage.AddressWithType, *storage.Upgrade, error) {
	msgType := storageTypes.MsgTryUpgrade
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	if err != nil {
		return msgType, addresses, nil, err
	}

	if status != storageTypes.StatusSuccess {
		return msgType, addresses, nil, nil
	}

	upgrade := &storage.Upgrade{
		Height:    ctx.Block.Height,
		Time:      ctx.Block.Time,
		EndTime:   ctx.Block.Time,
		EndHeight: ctx.Block.Height,
		Signer: &storage.Address{
			Address: m.Signer,
		},
	}
	return msgType, addresses, upgrade, err
}
