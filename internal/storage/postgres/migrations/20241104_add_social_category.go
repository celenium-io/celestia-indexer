// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upSocialCategory, downSocialCategory)
}

func upSocialCategory(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TYPE rollup_category ADD VALUE IF NOT EXISTS 'social' AFTER 'nft'`); err != nil {
		return err
	}
	return nil
}

func downSocialCategory(_ context.Context, _ *bun.DB) error {
	return nil
}
