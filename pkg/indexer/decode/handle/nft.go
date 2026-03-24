// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"cosmossdk.io/x/nft"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

// MsgSendNFT represents a message to send a nft from one account to another account.
func MsgSendNFT(ctx *context.Context, msgId uint64, m *nft.MsgSend) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSendNFT
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.Sender},
		{t: storageTypes.MsgAddressTypeReceiver, address: m.Receiver},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
