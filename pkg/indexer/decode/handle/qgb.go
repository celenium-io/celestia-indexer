// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/legacy"
)

// MsgRegisterEVMAddress registers an evm address to a validator.
func MsgRegisterEVMAddress(ctx *context.Context, msgId uint64, m *legacy.MsgRegisterEVMAddress) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgRegisterEVMAddress
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
		// TODO: think about EVM addresses
	}, ctx.Block.Height, msgId)
	return msgType, err
}
