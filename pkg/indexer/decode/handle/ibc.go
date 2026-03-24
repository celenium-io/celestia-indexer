// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	ibcTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

// IBCTransfer defines a msg to transfer fungible tokens (i.e., Coins) between
// ICS20 enabled chains. See ICS Spec here:
// https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#data-structures
func IBCTransfer(ctx *context.Context, msgId uint64, m *ibcTypes.MsgTransfer) (storageTypes.MsgType, error) {
	msgType := storageTypes.IBCTransfer
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.Sender},
		// {t: storageTypes.MsgAddressTypeReceiver,
		// address: m.Receiver}, // TODO: is it data to do IBC Transfer on cosmos network?
	}, ctx.Block.Height, msgId)
	return msgType, err
}
