// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddPriceFeed, downPriceFeed)
}

func upAddPriceFeed(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TYPE module_name ADD VALUE IF NOT EXISTS ? AFTER ?`, types.ModuleNameBaseapp.String(), types.ModuleNameConsensus.String())
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `ALTER TYPE module_name ADD VALUE IF NOT EXISTS ? AFTER ?`, types.ModuleNameIcahost.String(), types.ModuleNameBaseapp.String())
	return err
}
func downPriceFeed(ctx context.Context, db *bun.DB) error {
	return nil
}
