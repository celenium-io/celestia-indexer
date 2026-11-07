// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upHlTransferNonceBigint, downHlTransferNonceBigint)
}

// Hyperlane nonce is uint32, but bun mapped it to signed integer: nonces above 2^31-1 overflowed.
func upHlTransferNonceBigint(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE hl_transfer ALTER COLUMN nonce TYPE bigint`)
	return errors.Wrap(err, "alter hl_transfer.nonce to bigint")
}

// Fails if a nonce above 2^31-1 is already stored.
func downHlTransferNonceBigint(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE hl_transfer ALTER COLUMN nonce TYPE integer`)
	return errors.Wrap(err, "alter hl_transfer.nonce to integer")
}
