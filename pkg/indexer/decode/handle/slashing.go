// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosSlashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// MsgUnjail defines the Msg/Unjail request type
func MsgUnjail(ctx *context.Context, status storageTypes.Status, msgId uint64, m *cosmosSlashingTypes.MsgUnjail) (storageTypes.MsgType, []string, error) {
	msgType := storageTypes.MsgUnjail
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddr},
	}, ctx.Block.Height, msgId)
	if err != nil {
		return msgType, nil, err
	}
	if status == storageTypes.StatusFailed {
		return msgType, nil, nil
	}
	validator := storage.EmptyValidator()
	validator.Address = m.ValidatorAddr
	validator.MessagesCount += 1
	ctx.AddValidator(validator)
	return msgType, []string{m.ValidatorAddr}, nil
}

func MsgUpdateParamsSlashing(ctx *context.Context, msgId uint64, m *cosmosSlashingTypes.MsgUpdateParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
