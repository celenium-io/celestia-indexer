// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	qgbTypes "github.com/celestiaorg/celestia-app/x/qgb/types"
)

// MsgRegisterEVMAddress registers an evm address to a validator.
func MsgRegisterEVMAddress(ctx *context.Context, m *qgbTypes.MsgRegisterEVMAddress) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRegisterEVMAddress
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
		// TODO: think about EVM addresses
	}, ctx.Block.Height)
	return msgType, addresses, err
}
