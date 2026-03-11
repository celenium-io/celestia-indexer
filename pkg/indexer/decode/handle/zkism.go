// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	zkismTypes "github.com/celestiaorg/celestia-app/v7/x/zkism/types"
)

// MsgCreateInterchainSecurityModule
func MsgCreateInterchainSecurityModule(ctx *context.Context, msgId uint64, m *zkismTypes.MsgCreateInterchainSecurityModule) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateInterchainSecurityModule
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.Creator},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgUpdateInterchainSecurityModule
func MsgUpdateInterchainSecurityModule(ctx *context.Context, msgId uint64, m *zkismTypes.MsgUpdateInterchainSecurityModule) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateInterchainSecurityModule
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgSubmitMessages
func MsgSubmitMessages(ctx *context.Context, msgId uint64, m *zkismTypes.MsgSubmitMessages) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSubmitMessages
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
