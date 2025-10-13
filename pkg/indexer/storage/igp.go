// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/pkg/errors"
)

func (module *Module) saveIgps(
	ctx context.Context,
	tx storage.Transaction,
	dCtx *decodeContext.Context,
	addrToId map[string]uint64,
) error {
	igps := dCtx.Igps.Values()
	for i := range igps {
		addressId, ok := addrToId[igps[i].Owner.Address]
		if !ok {
			return errors.Wrapf(errCantFindAddress, "owner address %s", igps[i].Owner.Address)
		}
		igps[i].OwnerId = addressId
	}

	if err := tx.SaveHyperlaneIgps(ctx, igps...); err != nil {
		return err
	}

	if dCtx.IgpConfigs.Len() > 0 {
		configs := make([]storage.HLIGPConfig, 0)
		if err := dCtx.IgpConfigs.Range(func(igpAddress string, value *storage.HLIGPConfig) (error, bool) {
			hexAddress, err := util.DecodeHexAddress(igpAddress)
			if err != nil {
				return errors.Wrap(err, "decode igp address"), false
			}

			igp, err := tx.HyperlaneIgp(ctx, hexAddress.Bytes())
			if err != nil {
				return errors.Wrapf(err, "can't find igp with this address %s", hexAddress), false
			}
			value.Id = igp.Id

			configs = append(configs, *value)
			return nil, false
		}); err != nil {
			return err
		}

		if err := tx.SaveHyperlaneIgpConfigs(ctx, configs...); err != nil {
			return err
		}
	}

	return nil
}
