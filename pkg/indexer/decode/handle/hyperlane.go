// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	hyperlaneICS "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	hyperlanePostDispatch "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	hyperlaneCore "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	hyperlaneWarp "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

// MsgCreateMailbox
func MsgCreateMailbox(ctx *context.Context, msgId uint64, m *hyperlaneCore.MsgCreateMailbox) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateMailbox
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgProcessMessage
func MsgProcessMessage(ctx *context.Context, msgId uint64, m *hyperlaneCore.MsgProcessMessage) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgProcessMessage
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeRelayer, address: m.GetRelayer()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgSetMailbox
func MsgSetMailbox(ctx *context.Context, msgId uint64, m *hyperlaneCore.MsgSetMailbox) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSetMailbox

	items := addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}
	if m.GetNewOwner() != "" {
		items = append(items,
			addressData{t: storageTypes.MsgAddressTypeOwner, address: m.GetNewOwner()},
		)
	}
	err := createAddresses(ctx, items, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateCollateralToken
func MsgCreateCollateralToken(ctx *context.Context, msgId uint64, m *hyperlaneWarp.MsgCreateCollateralToken) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateCollateralToken
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateSyntheticToken
func MsgCreateSyntheticToken(ctx *context.Context, msgId uint64, m *hyperlaneWarp.MsgCreateSyntheticToken) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateSyntheticToken
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgSetToken
func MsgSetToken(ctx *context.Context, msgId uint64, m *hyperlaneWarp.MsgSetToken) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSetToken
	data := addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}
	if m.GetNewOwner() != "" {
		data = append(data, addressData{t: storageTypes.MsgAddressTypeOwner, address: m.GetNewOwner()})
	}
	err := createAddresses(ctx, data, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgEnrollRemoteRouter
func MsgEnrollRemoteRouter(ctx *context.Context, msgId uint64, m *hyperlaneWarp.MsgEnrollRemoteRouter) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgEnrollRemoteRouter
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgUnrollRemoteRouter
func MsgUnrollRemoteRouter(ctx *context.Context, msgId uint64, m *hyperlaneWarp.MsgUnrollRemoteRouter) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUnrollRemoteRouter
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgRemoteTransfer
func MsgRemoteTransfer(ctx *context.Context, msgId uint64, m *hyperlaneWarp.MsgRemoteTransfer) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgRemoteTransfer
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.GetSender()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgClaim
func MsgClaim(ctx *context.Context, msgId uint64, m *hyperlanePostDispatch.MsgClaim) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgClaim
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.GetSender()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateIgp
func MsgCreateIgp(ctx *context.Context, msgId uint64, m *hyperlanePostDispatch.MsgCreateIgp) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateIgp
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgSetIgpOwner
func MsgSetIgpOwner(ctx *context.Context, msgId uint64, m *hyperlanePostDispatch.MsgSetIgpOwner) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSetIgpOwner
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgPayForGas
func MsgPayForGas(ctx *context.Context, msgId uint64, m *hyperlanePostDispatch.MsgPayForGas) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgPayForGas
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.GetSender()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgSetDestinationGasConfig
func MsgSetDestinationGasConfig(ctx *context.Context, msgId uint64, m *hyperlanePostDispatch.MsgSetDestinationGasConfig) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSetDestinationGasConfig
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateMerkleTreeHook
func MsgCreateMerkleTreeHook(ctx *context.Context, msgId uint64, m *hyperlanePostDispatch.MsgCreateMerkleTreeHook) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateMerkleTreeHook
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateNoopHook
func MsgCreateNoopHook(ctx *context.Context, msgId uint64, m *hyperlanePostDispatch.MsgCreateNoopHook) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateNoopHook
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.GetOwner()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgAnnounceValidator
func MsgAnnounceValidator(ctx *context.Context, msgId uint64, m *hyperlaneICS.MsgAnnounceValidator) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgAnnounceValidator
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.GetCreator()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateMessageIdMultisigIsm
func MsgCreateMessageIdMultisigIsm(ctx *context.Context, msgId uint64, m *hyperlaneICS.MsgCreateMessageIdMultisigIsm) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateMessageIdMultisigIsm
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.GetCreator()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateMerkleRootMultisigIsm
func MsgCreateMerkleRootMultisigIsm(ctx *context.Context, msgId uint64, m *hyperlaneICS.MsgCreateMerkleRootMultisigIsm) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateMerkleRootMultisigIsm
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.GetCreator()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateNoopIsm
func MsgCreateNoopIsm(ctx *context.Context, msgId uint64, m *hyperlaneICS.MsgCreateNoopIsm) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateNoopIsm
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.GetCreator()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateRoutingIsm
func MsgCreateRoutingIsm(ctx *context.Context, msgId uint64, m *hyperlaneICS.MsgCreateRoutingIsm) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateRoutingIsm
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.GetCreator()},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
