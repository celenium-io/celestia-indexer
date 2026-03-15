// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	circuitTypes "cosmossdk.io/x/circuit/types"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

func MsgAuthorizeCircuitBreaker(ctx *context.Context, msgId uint64, m *circuitTypes.MsgAuthorizeCircuitBreaker) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgAuthorizeCircuitBreaker
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

func MsgResetCircuitBreaker(ctx *context.Context, msgId uint64, m *circuitTypes.MsgResetCircuitBreaker) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgResetCircuitBreaker
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

func MsgTripCircuitBreaker(ctx *context.Context, msgId uint64, m *circuitTypes.MsgTripCircuitBreaker) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgTripCircuitBreaker
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
