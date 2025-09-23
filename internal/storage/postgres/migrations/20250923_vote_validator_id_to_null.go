// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upResetVoteValidatorIdToNull, downResetVoteValidatorIdToNull)
}

func upResetVoteValidatorIdToNull(ctx context.Context, db *bun.DB) error {
	_, err := db.NewUpdate().Model((*storage.Vote)(nil)).
		Set("validator_id = NULL").
		Where("validator_id = 0").
		Exec(ctx)
	return err
}
func downResetVoteValidatorIdToNull(ctx context.Context, db *bun.DB) error {
	return nil
}
