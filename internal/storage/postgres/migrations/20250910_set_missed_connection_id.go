// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upSetMissedConnectionId, downSetMissedConnectionId)
}

func upSetMissedConnectionId(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `update ibc_transfer 
			set connection_id = ibc_channel.connection_id
			from ibc_channel
			where ibc_transfer.connection_id = '' and ibc_channel.id = ibc_transfer.channel_id`); err != nil {
		return err
	}
	return nil
}
func downSetMissedConnectionId(ctx context.Context, db *bun.DB) error {
	return nil
}
