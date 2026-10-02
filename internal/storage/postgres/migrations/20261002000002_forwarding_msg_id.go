// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upForwardingMsgId, downForwardingMsgId)
}

func upForwardingMsgId(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE forwarding ADD COLUMN IF NOT EXISTS msg_id bigint NOT NULL DEFAULT 0`); err != nil {
		return err
	}

	// Pairs the n-th forwarding of a tx with its n-th MsgForward, only where the counts match.
	// MsgExec-wrapped forwards were saved with tx_id = 0 and stay unmatched.
	var pairs []struct {
		Id    uint64    `bun:"id"`
		Time  time.Time `bun:"time"`
		MsgId uint64    `bun:"msg_id"`
	}
	if err := db.NewRaw(`
		WITH f AS (
			SELECT id, time, tx_id,
				row_number() OVER (PARTITION BY tx_id ORDER BY id) AS rn,
				count(*) OVER (PARTITION BY tx_id) AS cnt
			FROM forwarding
			WHERE msg_id = 0 AND tx_id <> 0
		), m AS (
			SELECT id, tx_id,
				row_number() OVER (PARTITION BY tx_id ORDER BY position) AS rn,
				count(*) OVER (PARTITION BY tx_id) AS cnt
			FROM message
			WHERE type = 'MsgForward'
				AND tx_id IN (SELECT tx_id FROM f)
				AND time >= (SELECT min(time) FROM f)
				AND time <= (SELECT max(time) FROM f)
		)
		SELECT f.id, f.time, m.id AS msg_id
		FROM f
		JOIN m ON m.tx_id = f.tx_id AND m.rn = f.rn AND m.cnt = f.cnt
	`).Scan(ctx, &pairs); err != nil {
		return err
	}

	// forwarding is a hypertable: plain per-row updates, no subqueries in UPDATE
	for i := range pairs {
		if _, err := db.NewUpdate().
			Table("forwarding").
			Set("msg_id = ?", pairs[i].MsgId).
			Where("id = ?", pairs[i].Id).
			Where("time = ?", pairs[i].Time).
			Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

func downForwardingMsgId(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE forwarding DROP COLUMN IF EXISTS msg_id`)
	return err
}
