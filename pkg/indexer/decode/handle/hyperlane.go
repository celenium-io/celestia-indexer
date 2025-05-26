// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	hyperlaneCore "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	hyperlaneWarp "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

// MsgCreateMailbox
func MsgCreateMailbox(ctx *context.Context, m *hyperlaneCore.MsgCreateMailbox) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateMailbox
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgProcessMessage
func MsgProcessMessage(ctx *context.Context, m *hyperlaneCore.MsgProcessMessage) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgProcessMessage
	addresses, err := createAddresses(ctx, addressesData{
		// {t: storageTypes.MsgAddressTypeRelayer, address: m.GetRelayer()},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgSetMailbox
func MsgSetMailbox(ctx *context.Context, m *hyperlaneCore.MsgSetMailbox) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSetMailbox
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetNewOwner()},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgCreateCollateralToken
func MsgCreateCollateralToken(ctx *context.Context, m *hyperlaneWarp.MsgCreateCollateralToken) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateCollateralToken
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgCreateSyntheticToken
func MsgCreateSyntheticToken(ctx *context.Context, m *hyperlaneWarp.MsgCreateSyntheticToken) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateSyntheticToken
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgSetToken
func MsgSetToken(ctx *context.Context, m *hyperlaneWarp.MsgSetToken) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSetToken
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetNewOwner()},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgEnrollRemoteRouter
func MsgEnrollRemoteRouter(ctx *context.Context, m *hyperlaneWarp.MsgEnrollRemoteRouter) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgEnrollRemoteRouter
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgUnrollRemoteRouter
func MsgUnrollRemoteRouter(ctx *context.Context, m *hyperlaneWarp.MsgUnrollRemoteRouter) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUnrollRemoteRouter
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgRemoteTransfer
func MsgRemoteTransfer(ctx *context.Context, m *hyperlaneWarp.MsgRemoteTransfer) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRemoteTransfer
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.GetSender()},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
