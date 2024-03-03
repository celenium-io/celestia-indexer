// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	interchainAccounts "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
)

// MsgRegisterInterchainAccount defines the payload for Msg/MsgRegisterInterchainAccount
func MsgRegisterInterchainAccount(ctx *context.Context, m *interchainAccounts.MsgRegisterInterchainAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRegisterInterchainAccount
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.Owner},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgSendTx defines the payload for Msg/SendTx
func MsgSendTx(ctx *context.Context, m *interchainAccounts.MsgSendTx) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSendTx
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.Owner},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
