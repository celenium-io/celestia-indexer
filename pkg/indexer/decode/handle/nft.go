// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

// MsgSendNFT represents a message to send a nft from one account to another account.
func MsgSendNFT(ctx *context.Context, m *nft.MsgSend) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSendNFT
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.Sender},
		{t: storageTypes.MsgAddressTypeReceiver, address: m.Receiver},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
