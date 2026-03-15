// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// MsgVerifyInvariant represents a message to verify a particular invariance.
func MsgVerifyInvariant(ctx *context.Context, msgId uint64, m *crisisTypes.MsgVerifyInvariant) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgVerifyInvariant
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.Sender},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
