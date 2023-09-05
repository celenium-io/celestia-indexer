package postgres

import (
	"context"
	"database/sql"

	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
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
			"tx_address_type",
			bun.Safe("tx_address_type"),
			bun.In(types.TxAddressTypeValues()),
		); err != nil {
			return err
		}
		return nil
	})
}
