// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddCreationTimeValidator, downAddCreationTimeValidator)
}

type msg struct {
	bun.BaseModel `bun:"message" comment:"Table with celestia messages."`

	Time time.Time         `bun:"time,pk,notnull"          comment:"The time of block"`
	Data types.PackedBytes `bun:"data,type:bytea,nullzero" comment:"Message data"`
}

func upAddCreationTimeValidator(ctx context.Context, db *bun.DB) error {
	if _, err := db.NewRaw("ALTER TABLE public.validator ADD COLUMN IF NOT EXISTS creation_time timestamptz NULL").Exec(ctx); err != nil {
		return err
	}
	if _, err := db.NewRaw("ALTER TABLE public.validator ALTER COLUMN creation_time SET STORAGE PLAIN").Exec(ctx); err != nil {
		return err
	}
	if _, err := db.NewRaw("COMMENT ON COLUMN public.validator.creation_time IS 'Creation time';").Exec(ctx); err != nil {
		return err
	}

	var msgs []msg
	if err := db.NewSelect().Model(&msgs).
		Column("msg.time", "data").
		Where("tx.status = 'success'").
		Where("type = 'MsgCreateValidator'").
		Join("left join tx on tx.id = tx_id").
		Scan(ctx); err != nil {
		return errors.Wrap(err, "receiving messages")
	}

	for i := range msgs {
		addr, ok := msgs[i].Data["ValidatorAddress"]
		if !ok {
			continue
		}
		valAddr, ok := addr.(string)
		if !ok {
			continue
		}
		if _, err := db.NewUpdate().Model((*storage.Validator)(nil)).Set("creation_time = ?", msgs[i].Time).Where("address = ?", valAddr).Exec(ctx); err != nil {
			return errors.Wrap(err, valAddr)
		}
	}

	return nil
}
func downAddCreationTimeValidator(ctx context.Context, db *bun.DB) error {
	return nil
}
