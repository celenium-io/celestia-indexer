// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	coreConnection "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
)

// MsgConnectionOpenInit defines the msg sent by an account on Chain A to initialize a connection with Chain B.
func MsgConnectionOpenInit(ctx *context.Context, msgId uint64, m *coreConnection.MsgConnectionOpenInit) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgConnectionOpenInit
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgConnectionOpenTry defines a msg sent by a Relayer to try to open a connection on Chain B.
func MsgConnectionOpenTry(ctx *context.Context, msgId uint64, m *coreConnection.MsgConnectionOpenTry) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgConnectionOpenTry
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgConnectionOpenAck defines a msg sent by a Relayer to Chain A to
// acknowledge the change of connection state to TRYOPEN on Chain B.
func MsgConnectionOpenAck(ctx *context.Context, msgId uint64, m *coreConnection.MsgConnectionOpenAck) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgConnectionOpenAck
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgConnectionOpenConfirm defines a msg sent by a Relayer to Chain B to
// acknowledge the change of connection state to OPEN on Chain A.
func MsgConnectionOpenConfirm(ctx *context.Context, msgId uint64, m *coreConnection.MsgConnectionOpenConfirm) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgConnectionOpenConfirm
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

func MsgUpdateParamsConnection(ctx *context.Context, msgId uint64, m *coreConnection.MsgUpdateParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
