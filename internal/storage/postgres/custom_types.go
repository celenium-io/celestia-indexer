// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestial-module/pkg/storage/postgres"
	"github.com/dipdup-net/go-lib/database"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

const (
	createTypeQuery = `DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = ?) THEN
			CREATE TYPE ? AS ENUM (?);
		END IF;
	END$$;`
)

func createTypes(ctx context.Context, conn *database.Bun) error {
	if err := postgres.CreateTypes(ctx, conn); err != nil {
		return err
	}

	log.Info().Msg("creating custom types...")
	return conn.DB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"event_type",
			bun.Safe("event_type"),
			bun.In(types.EventTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"msg_type",
			bun.Safe("msg_type"),
			bun.In(types.MsgTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"status",
			bun.Safe("status"),
			bun.In(types.StatusValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"msg_address_type",
			bun.Safe("msg_address_type"),
			bun.In(types.MsgAddressTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"module_name",
			bun.Safe("module_name"),
			bun.In(types.ModuleNameValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"staking_log_type",
			bun.Safe("staking_log_type"),
			bun.In(types.StakingLogTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"vesting_type",
			bun.Safe("vesting_type"),
			bun.In(types.VestingTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"rollup_type",
			bun.Safe("rollup_type"),
			bun.In(types.RollupTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"rollup_category",
			bun.Safe("rollup_category"),
			bun.In(types.RollupCategoryValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"proposal_status",
			bun.Safe("proposal_status"),
			bun.In(types.ProposalStatusValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"proposal_type",
			bun.Safe("proposal_type"),
			bun.In(types.ProposalTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"vote_option",
			bun.Safe("vote_option"),
			bun.In(types.VoteOptionValues()),
		); err != nil {
			return err
		}
		return nil
	})
}
