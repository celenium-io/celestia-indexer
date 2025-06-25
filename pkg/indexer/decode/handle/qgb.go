// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/legacy"
)

// MsgRegisterEVMAddress registers an evm address to a validator.
func MsgRegisterEVMAddress(ctx *context.Context, m *legacy.MsgRegisterEVMAddress) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRegisterEVMAddress
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
		// TODO: think about EVM addresses
	}, ctx.Block.Height)
	return msgType, addresses, err
}
