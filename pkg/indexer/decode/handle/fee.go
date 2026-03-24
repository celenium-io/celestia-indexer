// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	minfeeTypes "github.com/celestiaorg/celestia-app/v7/x/minfee/types"
	fee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
)

// MsgRegisterPayee defines the request type for the RegisterPayee rpc
func MsgRegisterPayee(ctx *context.Context, msgId uint64, m *fee.MsgRegisterPayee) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgRegisterPayee
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeRelayer, address: m.Relayer},
		{t: storageTypes.MsgAddressTypePayee, address: m.Payee},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgRegisterCounterpartyPayee defines the request type for the RegisterCounterpartyPayee rpc
func MsgRegisterCounterpartyPayee(ctx *context.Context, msgId uint64, m *fee.MsgRegisterCounterpartyPayee) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgRegisterCounterpartyPayee
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeRelayer, address: m.Relayer},
		{t: storageTypes.MsgAddressTypePayee, address: m.CounterpartyPayee}, // the counterparty payee address
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgPayPacketFee defines the request type for the PayPacketFee rpc
// This Msg can be used to pay for a packet at the next sequence send & should be combined with the Msg that will be
// paid for
func MsgPayPacketFee(ctx *context.Context, msgId uint64, m *fee.MsgPayPacketFee) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgPayPacketFee
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgPayPacketFeeAsync defines the request type for the PayPacketFeeAsync rpc
// This Msg can be used to pay for a packet at a specified sequence (instead of the next sequence sends)
func MsgPayPacketFeeAsync() (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgPayPacketFeeAsync
	return msgType, nil
}

// MsgUpdateMinfeeParams defines a message for updating the minimum fee parameters.
func MsgUpdateMinfeeParams(ctx *context.Context, msgId uint64, m *minfeeTypes.MsgUpdateMinfeeParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateMinfeeParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
