// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	coreConnection "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
)

// MsgConnectionOpenInit defines the msg sent by an account on Chain A to initialize a connection with Chain B.
func MsgConnectionOpenInit(ctx *context.Context, m *coreConnection.MsgConnectionOpenInit) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgConnectionOpenInit
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgConnectionOpenTry defines a msg sent by a Relayer to try to open a connection on Chain B.
func MsgConnectionOpenTry(ctx *context.Context, m *coreConnection.MsgConnectionOpenTry) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgConnectionOpenTry
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgConnectionOpenAck defines a msg sent by a Relayer to Chain A to
// acknowledge the change of connection state to TRYOPEN on Chain B.
func MsgConnectionOpenAck(ctx *context.Context, m *coreConnection.MsgConnectionOpenAck) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgConnectionOpenAck
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgConnectionOpenConfirm defines a msg sent by a Relayer to Chain B to
// acknowledge the change of connection state to OPEN on Chain A.
func MsgConnectionOpenConfirm(ctx *context.Context, m *coreConnection.MsgConnectionOpenConfirm) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgConnectionOpenConfirm
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
