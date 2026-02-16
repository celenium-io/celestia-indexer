// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	zkismTypes "github.com/celestiaorg/celestia-app/v7/x/zkism/types"
)

// MsgCreateInterchainSecurityModule
func MsgCreateInterchainSecurityModule(ctx *context.Context, m *zkismTypes.MsgCreateInterchainSecurityModule) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateInterchainSecurityModule
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.Creator},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgUpdateInterchainSecurityModule
func MsgUpdateInterchainSecurityModule(ctx *context.Context, m *zkismTypes.MsgUpdateInterchainSecurityModule) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateInterchainSecurityModule
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgSubmitMessages
func MsgSubmitMessages(ctx *context.Context, m *zkismTypes.MsgSubmitMessages) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSubmitMessages
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
