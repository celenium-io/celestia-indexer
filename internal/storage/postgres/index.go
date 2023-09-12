package postgres

import (
	"context"
	"database/sql"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

func createIndices(ctx context.Context, conn *database.Bun) error {
	log.Info().Msg("creating indexes...")
	return conn.DB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// Address
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Address)(nil)).
			Index("address_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Address)(nil)).
			Index("address_hash_idx").
			Column("hash").
			Exec(ctx); err != nil {
			return err
		}

		// Block
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Address)(nil)).
			Index("block_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}

		// BlockStats
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Address)(nil)).
			Index("block_stats_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}

		// Tx
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Tx)(nil)).
			Index("tx_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Tx)(nil)).
			Index("tx_hash_idx").
			Column("hash").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Tx)(nil)).
			Index("tx_status_idx").
			Column("status").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Tx)(nil)).
			Index("tx_message_types_idx").
			Column("message_types").
			Exec(ctx); err != nil {
			return err
		}

		// Event
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Event)(nil)).
			Index("event_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Event)(nil)).
			Index("event_tx_id_idx").
			Column("tx_id").
			Where("tx_id IS NOT NULL").
			Exec(ctx); err != nil {
			return err
		}

		// Message
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Message)(nil)).
			Index("message_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Message)(nil)).
			Index("message_tx_id_idx").
			Column("tx_id").
			Where("tx_id IS NOT NULL").
			Exec(ctx); err != nil {
			return err
		}

		// Namespace
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Namespace)(nil)).
			Index("namespace_idx").
			Column("namespace_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Namespace)(nil)).
			Index("namespace_version_idx").
			Column("version").
			Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}
