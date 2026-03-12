// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upDropAutoincrement, downDropAutoincrement)
}

func upDropAutoincrement(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."tx" ALTER COLUMN "id" DROP DEFAULT;`); err != nil {
		return errors.Wrap(err, "drop default")
	}
	if _, err := db.ExecContext(ctx, `DROP sequence tx_id_seq;`); err != nil {
		return errors.Wrap(err, "drop sequence tx_id_seq")
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."message" ALTER COLUMN "id" DROP DEFAULT;`); err != nil {
		return errors.Wrap(err, "drop default")
	}
	if _, err := db.ExecContext(ctx, `DROP sequence message_id_seq;`); err != nil {
		return errors.Wrap(err, "drop sequence message_id_seq")
	}
	return nil
}

func downDropAutoincrement(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE SEQUENCE tx_id_seq;`); err != nil {
		return errors.Wrap(err, "create sequence tx_id_seq")
	}
	// setval sets the current value so the next nextval() call returns MAX(id)+1.
	// COALESCE handles the empty-table case (starts from 1).
	if _, err := db.ExecContext(ctx, `SELECT setval('tx_id_seq', COALESCE((SELECT MAX(id) FROM public."tx"), 0));`); err != nil {
		return errors.Wrap(err, "set value tx_id_seq")
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."tx" ALTER COLUMN "id" SET DEFAULT nextval('tx_id_seq');`); err != nil {
		return errors.Wrap(err, "set default tx id")
	}
	if _, err := db.ExecContext(ctx, `ALTER SEQUENCE tx_id_seq OWNED BY public."tx"."id";`); err != nil {
		return errors.Wrap(err, "set owner tx_id_seq")
	}

	if _, err := db.ExecContext(ctx, `CREATE SEQUENCE message_id_seq;`); err != nil {
		return errors.Wrap(err, "create sequence message_id_seq")
	}
	if _, err := db.ExecContext(ctx, `SELECT setval('message_id_seq', COALESCE((SELECT MAX(id) FROM public."message"), 0));`); err != nil {
		return errors.Wrap(err, "set value message_id_seq")
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."message" ALTER COLUMN "id" SET DEFAULT nextval('message_id_seq');`); err != nil {
		return errors.Wrap(err, "set default message id")
	}
	if _, err := db.ExecContext(ctx, `ALTER SEQUENCE message_id_seq OWNED BY public."message"."id";`); err != nil {
		return errors.Wrap(err, "set owner message_id_seq")
	}

	return nil
}
