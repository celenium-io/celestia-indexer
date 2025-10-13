// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upUpgrades, downUpgrades)
}

func upUpgrades(ctx context.Context, db *bun.DB) error {
	if _, err := db.Exec(`CREATE TABLE public.upgrade_new (
		height int8 NULL,
		end_height int8 NULL,
		signer_id int8 NULL,
		"time" timestamptz NOT NULL,
		"end_time" timestamptz NOT NULL,
		"version" int8 NULL,
		msg_id int8 NOT NULL,
		tx_id int8 NOT NULL,
		voting_power numeric NULL,
		voted_power numeric NULL,
		signals_count int4 NULL,
		CONSTRAINT upgrade_new_pkey PRIMARY KEY(version)
	);`); err != nil {
		return errors.Wrap(err, "create buffer table")
	}
	if _, err := db.Exec(` INSERT INTO upgrade_new (height, end_height, signer_id, time, end_time, version, msg_id, tx_id, voting_power, voted_power, signals_count)
    	SELECT height, height, signer_id, time, time, version, msg_id, tx_id, voting_power, voted_power, 0  FROM upgrade;`); err != nil {
		return errors.Wrap(err, "fill buffer table")
	}
	if _, err := db.Exec("drop table upgrade;"); err != nil {
		return errors.Wrap(err, "drop table upgrade")
	}

	if _, err := db.Exec(`ALTER TABLE upgrade_new RENAME TO upgrade;`); err != nil {
		return errors.Wrap(err, "rename table")
	}

	if _, err := db.Exec(`ALTER TABLE upgrade RENAME CONSTRAINT upgrade_new_pkey TO upgrade_pkey;`); err != nil {
		return errors.Wrap(err, "rename constraint")
	}

	var upgrades []storage.Upgrade
	if err := db.NewSelect().Model(&upgrades).Scan(ctx); err != nil {
		return errors.Wrap(err, "receiving upgrades")
	}

	for i := range upgrades {
		var signalsCount int
		if err := db.NewSelect().
			Model((*storage.SignalVersion)(nil)).
			ColumnExpr("count(*)").
			Where("version = ?", upgrades[i].Version).
			Scan(ctx, &signalsCount); err != nil {
			return errors.Wrapf(err, "get signals count for upgrade to version %d", upgrades[i].Version)
		}

		var firstSignal storage.SignalVersion
		if err := db.NewSelect().
			Model(&firstSignal).
			Where("version = ?", upgrades[i].Version).
			Limit(1).
			OrderExpr("time asc").
			Scan(ctx); err != nil {
			return errors.Wrapf(err, "get first signal for upgrade to version %d", upgrades[i].Version)
		}

		if _, err := db.NewUpdate().
			Model(&upgrades[i]).
			Where("version = ?", upgrades[i].Version).
			Set("signals_count = ?", signalsCount).
			Set("height = ?", firstSignal.Height).
			Set("time = ?", firstSignal.Time).
			Exec(ctx); err != nil {
			return errors.Wrapf(err, "update row with upgrade to version %d", upgrades[i].Version)
		}
	}

	var state storage.State
	if err := db.NewSelect().Model(&state).Limit(1).Scan(ctx); err != nil {
		return errors.Wrap(err, "receiving state")
	}

	currentUpgrade := state.Version + 1
	var signals []storage.SignalVersion
	if err := db.NewSelect().Model(&signals).Where("version = ?", currentUpgrade).OrderExpr("time asc").Scan(ctx); err != nil {
		return errors.Wrap(err, "receiving signals for current upgrade")
	}

	if len(signals) == 0 {
		return nil
	}

	votedPower := decimal.Zero
	for i := range signals {
		votedPower = votedPower.Add(signals[i].VotingPower)
	}

	upgrade := storage.Upgrade{
		Version:      currentUpgrade,
		Height:       signals[0].Height,
		Time:         signals[0].Time,
		SignalsCount: len(signals),
		VotedPower:   votedPower,
	}

	if _, err := db.NewInsert().Model(&upgrade).Exec(ctx); err != nil {
		return errors.Wrapf(err, "insert current upgrade to version %d", currentUpgrade)
	}

	return nil
}
func downUpgrades(ctx context.Context, db *bun.DB) error {
	return nil
}
