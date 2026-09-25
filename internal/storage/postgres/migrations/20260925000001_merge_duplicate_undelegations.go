// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upMergeDuplicateUndelegations, downMergeDuplicateUndelegations)
}

// The SDK merges unbonding entries with the same creation height, so duplicates are summed
// into the row with the lowest id; createIndices then builds undelegation_unique_idx.
func upMergeDuplicateUndelegations(ctx context.Context, db *bun.DB) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// update and delete touch disjoint rows, so sharing one snapshot is safe
		if _, err := tx.ExecContext(ctx, `
			WITH dups AS (
				SELECT min(id) AS keep_id, sum(amount) AS amount, height, validator_id, address_id
				FROM undelegation
				-- GROUP BY treats nulls as equal but the delete join does not, and the unique index allows them
				WHERE validator_id IS NOT NULL AND address_id IS NOT NULL
				GROUP BY height, validator_id, address_id
				HAVING count(*) > 1
			), merged AS (
				UPDATE undelegation SET amount = dups.amount
				FROM dups
				WHERE undelegation.id = dups.keep_id
			)
			DELETE FROM undelegation
			USING dups
			WHERE undelegation.height = dups.height
				AND undelegation.validator_id = dups.validator_id
				AND undelegation.address_id = dups.address_id
				AND undelegation.id <> dups.keep_id
		`); err != nil {
			return errors.Wrap(err, "merge duplicate undelegations")
		}
		return nil
	})
}

// Merged rows cannot be split back, so down only drops the index that relies on the merge.
func downMergeDuplicateUndelegations(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `DROP INDEX IF EXISTS undelegation_unique_idx`)
	return errors.Wrap(err, "drop undelegation_unique_idx")
}
