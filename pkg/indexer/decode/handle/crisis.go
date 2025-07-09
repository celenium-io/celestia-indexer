// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// MsgVerifyInvariant represents a message to verify a particular invariance.
func MsgVerifyInvariant(ctx *context.Context, m *crisisTypes.MsgVerifyInvariant) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVerifyInvariant
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.Sender},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
