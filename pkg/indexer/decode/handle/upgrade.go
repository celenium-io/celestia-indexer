// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	upgrade "cosmossdk.io/x/upgrade/types"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

// MsgSoftwareUpgrade is the Msg/SoftwareUpgrade request type.
func MsgSoftwareUpgrade(ctx *context.Context, msgId uint64, m *upgrade.MsgSoftwareUpgrade) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSoftwareUpgrade
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgSoftwareUpgrade is the Msg/SoftwareUpgrade request type.
func MsgCancelUpgrade(ctx *context.Context, msgId uint64, m *upgrade.MsgCancelUpgrade) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCancelUpgrade
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
