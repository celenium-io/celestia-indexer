// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"encoding/json"

	"cosmossdk.io/x/feegrant"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/handle"
	"github.com/celestiaorg/celestia-app/v9/app"
	"github.com/celestiaorg/celestia-app/v9/app/encoding"
	"github.com/cosmos/cosmos-sdk/codec"
)

var cfg = encoding.MakeConfig(app.ModuleEncodingRegisters...)

func (module *Module) parseFeeGrants(
	feeGrantsRaw []json.RawMessage, block storage.Block, data *parsedData,
) error {
	cdc := codec.NewProtoCodec(cfg.InterfaceRegistry)

	ctx := decodeContext.NewContext()
	ctx.Block = &block
	for i := range feeGrantsRaw {
		var feeGrant feegrant.MsgGrantAllowance
		if err := cdc.UnmarshalJSON(feeGrantsRaw[i], &feeGrant); err != nil {
			return err
		}
		if _, err := handle.MsgGrantAllowance(ctx, storageTypes.StatusSuccess, 0, &feeGrant); err != nil {
			return err
		}
		if err := addAddress(data, block, feeGrant.Grantee); err != nil {
			return err
		}
		if err := addAddress(data, block, feeGrant.Granter); err != nil {
			return err
		}
	}
	data.grants = ctx.Grants.Values()
	return nil
}
