// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	upgrade "cosmossdk.io/x/upgrade/types"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

// MsgSoftwareUpgrade is the Msg/SoftwareUpgrade request type.
func MsgSoftwareUpgrade(ctx *context.Context, m *upgrade.MsgSoftwareUpgrade) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSoftwareUpgrade
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgSoftwareUpgrade is the Msg/SoftwareUpgrade request type.
func MsgCancelUpgrade(ctx *context.Context, m *upgrade.MsgCancelUpgrade) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCancelUpgrade
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
