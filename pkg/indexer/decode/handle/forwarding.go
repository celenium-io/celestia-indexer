// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	fwdTypes "github.com/celestiaorg/celestia-app/v7/x/forwarding/types"
)

// MsgForward
func MsgForward(ctx *context.Context, m *fwdTypes.MsgForward) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgForward
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
		{t: storageTypes.MsgAddressTypeReceiver, address: m.ForwardAddr},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
