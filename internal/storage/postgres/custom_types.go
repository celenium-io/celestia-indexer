// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestial-module/pkg/storage/postgres"
	"github.com/dipdup-io/go-lib/database"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

const (
	createTypeQuery = `DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = ?) THEN
			CREATE TYPE ? AS ENUM ?;
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
			bun.Tuple(types.EventTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"msg_type",
			bun.Safe("msg_type"),
			bun.Tuple(types.MsgTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"status",
			bun.Safe("status"),
			bun.Tuple(types.StatusValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"msg_address_type",
			bun.Safe("msg_address_type"),
			bun.Tuple(types.MsgAddressTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"module_name",
			bun.Safe("module_name"),
			bun.Tuple(types.ModuleNameValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"staking_log_type",
			bun.Safe("staking_log_type"),
			bun.Tuple(types.StakingLogTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"vesting_type",
			bun.Safe("vesting_type"),
			bun.Tuple(types.VestingTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"rollup_type",
			bun.Safe("rollup_type"),
			bun.Tuple(types.RollupTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"rollup_category",
			bun.Safe("rollup_category"),
			bun.Tuple(types.RollupCategoryValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"proposal_status",
			bun.Safe("proposal_status"),
			bun.Tuple(types.ProposalStatusValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"proposal_type",
			bun.Safe("proposal_type"),
			bun.Tuple(types.ProposalTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"vote_option",
			bun.Safe("vote_option"),
			bun.Tuple(types.VoteOptionValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"ibc_channel_status",
			bun.Safe("ibc_channel_status"),
			bun.Tuple(types.IbcChannelStatusValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"hyperlane_transfer_type",
			bun.Safe("hyperlane_transfer_type"),
			bun.Tuple(types.HLTransferTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"hyperlane_token_type",
			bun.Safe("hyperlane_token_type"),
			bun.Tuple(types.HLTokenTypeValues()),
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			createTypeQuery,
			"upgrade_status",
			bun.Safe("upgrade_status"),
			bun.Tuple(types.UpgradeStatusValues()),
		); err != nil {
			return err
		}
		return nil
	})
}
