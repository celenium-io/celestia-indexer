// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	interchainAccounts "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	interchainAccountsHost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
)

// MsgRegisterInterchainAccount defines the payload for Msg/MsgRegisterInterchainAccount
func MsgRegisterInterchainAccount(ctx *context.Context, msgId uint64, m *interchainAccounts.MsgRegisterInterchainAccount) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgRegisterInterchainAccount
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.Owner},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgSendTx defines the payload for Msg/SendTx
func MsgSendTx(ctx *context.Context, msgId uint64, m *interchainAccounts.MsgSendTx) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSendTx
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.Owner},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

func MsgModuleQuerySafe(ctx *context.Context, msgId uint64, m *interchainAccountsHost.MsgModuleQuerySafe) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgModuleQuerySafe
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

func MsgUpdateParamsIcaHost(ctx *context.Context, msgId uint64, m *interchainAccountsHost.MsgUpdateParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

func MsgUpdateParamsIcaController(ctx *context.Context, msgId uint64, m *interchainAccounts.MsgUpdateParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
