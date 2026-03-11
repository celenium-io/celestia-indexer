// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

func processGrants(addrToId map[string]uint64, grant *storage.Grant) error {
	granteeId, ok := addrToId[grant.Grantee.Address]
	if !ok {
		return errors.Wrapf(errCantFindAddress, "grantee: %s", grant.Grantee.Address)
	}
	grant.GranteeId = granteeId
	granterId, ok := addrToId[grant.Granter.Address]
	if !ok {
		return errors.Wrapf(errCantFindAddress, "granter: %s", grant.Granter.Address)
	}
	grant.GranterId = granterId
	return nil
}

func saveGrants(
	ctx context.Context,
	tx storage.Transaction,
	grants []*storage.Grant,
	addrToId map[string]uint64,
) error {
	if len(grants) == 0 {
		return nil
	}

	for key := range grants {
		if err := processGrants(addrToId, grants[key]); err != nil {
			return errors.Wrap(err, "process grant")
		}
	}

	if err := tx.SaveGrants(ctx, grants...); err != nil {
		return errors.Wrap(err, "saving grants")
	}

	return nil
}
