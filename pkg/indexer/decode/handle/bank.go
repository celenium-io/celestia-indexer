// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// MsgSend represents a message to send coins from one account to another.
func MsgSend(ctx *context.Context, m *cosmosBankTypes.MsgSend) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSend
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgMultiSend represents an arbitrary multi-in, multi-out send message.
func MsgMultiSend(ctx *context.Context, m *cosmosBankTypes.MsgMultiSend) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgMultiSend
	aData := make(addressesData, len(m.Inputs)+len(m.Outputs))

	var i int64
	for _, input := range m.Inputs {
		aData[i] = addressData{t: storageTypes.MsgAddressTypeInput, address: input.Address}
		i++
	}
	for _, output := range m.Outputs {
		aData[i] = addressData{t: storageTypes.MsgAddressTypeOutput, address: output.Address}
		i++
	}

	addresses, err := createAddresses(ctx, aData, ctx.Block.Height)
	return msgType, addresses, err
}
