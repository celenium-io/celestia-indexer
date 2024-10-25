// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

// MsgSignalVersion -
func MsgSignalVersion(ctx *context.Context, validatorAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSignalVersion
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: validatorAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgTryUpgrade -
func MsgTryUpgrade(ctx *context.Context, signer string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTryUpgrade
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
