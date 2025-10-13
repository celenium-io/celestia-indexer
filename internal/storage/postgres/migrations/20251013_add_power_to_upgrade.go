// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddPowerToValidator, downAddPowerToValidator)
}

func upAddPowerToValidator(ctx context.Context, db *bun.DB) error {
	if _, err := db.NewRaw("ALTER TABLE public.upgrade ADD COLUMN IF NOT EXISTS voting_power numeric NULL;").Exec(ctx); err != nil {
		return err
	}
	if _, err := db.NewRaw("ALTER TABLE public.upgrade ALTER COLUMN voting_power SET STORAGE MAIN").Exec(ctx); err != nil {
		return err
	}
	if _, err := db.NewRaw("COMMENT ON COLUMN public.upgrade.voting_power IS 'Total voting power on upgrade block'").Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewRaw("ALTER TABLE public.upgrade ADD COLUMN IF NOT EXISTS voted_power numeric NULL;").Exec(ctx); err != nil {
		return err
	}
	if _, err := db.NewRaw("ALTER TABLE public.upgrade ALTER COLUMN voted_power SET STORAGE MAIN").Exec(ctx); err != nil {
		return err
	}
	if _, err := db.NewRaw("COMMENT ON COLUMN public.upgrade.voted_power IS 'Total voting power of upgraded validators'").Exec(ctx); err != nil {
		return err
	}
	return nil
}
func downAddPowerToValidator(ctx context.Context, db *bun.DB) error {
	return nil
}
