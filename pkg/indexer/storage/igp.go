// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/pkg/errors"
)

func (module *Module) saveIgp(
	ctx context.Context,
	tx storage.Transaction,
	dCtx *decodeContext.Context,
	addrToId map[string]uint64,
) error {
	for i := range dCtx.Igps {
		addressId, ok := addrToId[dCtx.Igps[i].Owner.Address]
		if !ok {
			return errors.Wrapf(errCantFindAddress, "owner address %s", dCtx.Igps[i].Owner.Address)
		}
		dCtx.Igps[i].OwnerId = addressId
	}

	if err := tx.SaveIgps(ctx, dCtx.Igps...); err != nil {
		return err
	}

	return nil
}
