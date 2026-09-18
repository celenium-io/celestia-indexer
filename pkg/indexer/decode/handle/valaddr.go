// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	valaddrTypes "github.com/celestiaorg/celestia-app/v10/x/valaddr/types"
)

// MsgSetFibreProviderInfo advertises the host a validator serves fibre shards
// from. The signer is a validator operator address (celestiavaloper...), and
// the module stores the host under the derived consensus address.
func MsgSetFibreProviderInfo(
	ctx *context.Context, status storageTypes.Status, msgId uint64, m *valaddrTypes.MsgSetFibreProviderInfo,
) (storageTypes.MsgType, []string, error) {
	msgType := storageTypes.MsgSetFibreProviderInfo
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.Signer},
	}, ctx.Block.Height, msgId)
	if err != nil {
		return msgType, nil, err
	}
	if status != storageTypes.StatusSuccess {
		return msgType, nil, nil
	}

	validatorAddress := m.GetSigner()
	validator := storage.Validator{
		Address:         validatorAddress,
		MessagesCount:   1,
		FibreHost:       &m.Host,
		FibreHostHeight: &ctx.Block.Height,
	}
	ctx.AddValidator(validator)

	return msgType, []string{validatorAddress}, nil
}
