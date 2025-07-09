// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddRollupOtherTypeAndCategory, downAddRollupOtherTypeAndCategory)
}

func upAddRollupOtherTypeAndCategory(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TYPE rollup_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.RollupTypeOther.String(), types.RollupTypeSettled.String()); err != nil {
		return errors.Wrap(err, "add other rollup type")
	}
	return nil
}
func downAddRollupOtherTypeAndCategory(ctx context.Context, db *bun.DB) error {
	return nil
}
