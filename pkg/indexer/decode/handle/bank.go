// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// MsgSend represents a message to send coins from one account to another.
func MsgSend(ctx *context.Context, msgId uint64, m *cosmosBankTypes.MsgSend) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSend
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgMultiSend represents an arbitrary multi-in, multi-out send message.
func MsgMultiSend(ctx *context.Context, msgId uint64, m *cosmosBankTypes.MsgMultiSend) (storageTypes.MsgType, error) {
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

	err := createAddresses(ctx, aData, ctx.Block.Height, msgId)
	return msgType, err
}

func MsgSetSendEnabled(ctx *context.Context, msgId uint64, m *cosmosBankTypes.MsgSetSendEnabled) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSetSendEnabled
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

func MsgUpdateParamsBank(ctx *context.Context, msgId uint64, m *cosmosBankTypes.MsgUpdateParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
