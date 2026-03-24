// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	fwdTypes "github.com/celestiaorg/celestia-app/v7/x/forwarding/types"
)

// MsgForward
func MsgForward(ctx *context.Context, msgId uint64, m *fwdTypes.MsgForward) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgForward
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
		{t: storageTypes.MsgAddressTypeReceiver, address: m.ForwardAddr},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
