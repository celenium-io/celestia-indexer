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

func saveIgps(
	ctx context.Context,
	tx storage.HyperlaneTx,
	repo storage.IHLIGP,
	dCtx *decodeContext.Context,
	addrToId map[string]uint64,
) error {
	var (
		igps = make([]*storage.HLIGP, 0, dCtx.Igps.Len())
		ids  = make(map[string]uint64, dCtx.Igps.Len())
	)
	for igp := range dCtx.Igps.AllValues() {
		if igp.Owner == nil {
			igp.OwnerId = 0
		} else {
			addressId, ok := addrToId[igp.Owner.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "owner address %s", igp.Owner.Address)
			}
			igp.OwnerId = addressId
		}
		igps = append(igps, igp)
	}

	if err := tx.SaveHyperlaneIgps(ctx, igps...); err != nil {
		return err
	}

	for igpId, igp := range dCtx.Igps.All() {
		ids[igpId] = igp.Id
	}

	if dCtx.IgpConfigs.Len() > 0 {
		configs := make([]storage.HLIGPConfig, 0, dCtx.IgpConfigs.Len())

		for key, value := range dCtx.IgpConfigs.All() {
			id, ok := ids[key.IgpId]
			if !ok {
				hexAddress, err := util.DecodeHexAddress(key.IgpId)
				if err != nil {
					return errors.Wrapf(err, "decode igp id: %s", key.IgpId)
				}
				id, err = repo.IdByHash(ctx, hexAddress.Bytes())
				if err != nil {
					return errors.Wrapf(err, "can't find igp with this address %s", hexAddress)
				}
				ids[key.IgpId] = id
			}

			value.Id = id
			configs = append(configs, *value)
		}

		if err := tx.SaveHyperlaneIgpConfigs(ctx, configs...); err != nil {
			return err
		}
	}

	return nil
}
