// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"
	"strconv"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upConstantsFormat, downConstantsFormat)
}

func upConstantsFormat(ctx context.Context, db *bun.DB) error {
	var constants []storage.Constant
	err := db.NewSelect().
		Model(&constants).
		Where("name = 'evidence_max_age_duration'").
		WhereOr("name = 'max_deposit_period'").
		WhereOr("name = 'voting_period'").
		WhereOr("name = 'downtime_jail_duration'").
		WhereOr("name = 'unbonding_time'").
		Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "receiving constants")
	}
	if len(constants) != 5 {
		return errors.Errorf("count of constnats is wrong")
	}

	for i := range constants {
		value, err := time.ParseDuration(constants[i].Value)
		if err != nil {
			return errors.Wrapf(err, "parsing constant: %s", constants[i].Name)
		}

		_, err = db.NewUpdate().
			Model((*storage.Constant)(nil)).
			Where("name = ?", constants[i].Name).
			Where("module = ?", constants[i].Module).
			Set("value = ?", strconv.FormatInt(value.Nanoseconds(), 10)).
			Exec(ctx)
		if err != nil {
			return errors.Wrapf(err, "update constant %s", constants[i].Name)
		}
	}
	return nil
}
func downConstantsFormat(ctx context.Context, db *bun.DB) error {
	return nil
}
