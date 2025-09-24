// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

func (module *Module) saveIgps(
	ctx context.Context,
	tx storage.Transaction,
	igps []*storage.HLIGP,
	addrToId map[string]uint64,
) error {
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

	return nil
}
