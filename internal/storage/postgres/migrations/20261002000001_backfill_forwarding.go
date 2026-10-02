// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upBackfillIsForwarding, downBackfillIsForwarding)
}

func upBackfillIsForwarding(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		UPDATE address SET is_forwarding = true
		WHERE NOT is_forwarding
		  AND id IN (SELECT DISTINCT address_id FROM forwarding)
	`)
	return err
}

// Restores the pre-fix state: the flag was only set when the address was created in the block of its first forwarding.
func downBackfillIsForwarding(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		UPDATE address a SET is_forwarding = false
		FROM (SELECT address_id, min(height) AS height FROM forwarding GROUP BY address_id) f
		WHERE f.address_id = a.id AND a.height <> f.height
	`)
	return err
}
