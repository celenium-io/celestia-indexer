package postgres

import (
	"context"
	"database/sql"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	models "github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

// Storage -
type Storage struct {
	*postgres.Storage

	Blocks    models.IBlock
	Tx        models.ITx
	Message   models.IMessage
	Event     models.IEvent
	Address   models.IAddress
	Namespace models.INamespace
	State     models.IState

	PartitionManager database.RangePartitionManager
}

// Create -
func Create(ctx context.Context, cfg config.Database) (Storage, error) {
	strg, err := postgres.Create(ctx, cfg, initDatabase)
	if err != nil {
		return Storage{}, err
	}

	s := Storage{
		Storage: strg,
		Blocks:  NewBlocks(strg.Connection()),
		Message: NewMessage(strg.Connection()),
		Event:   NewEvent(strg.Connection()),
		Address: NewAddress(strg.Connection()),
		Tx:      NewTx(strg.Connection()),
		State:   NewState(strg.Connection()),

		PartitionManager: database.NewPartitionManager(strg.Connection(), database.PartitionByMonth),
	}

	return s, nil
}

func initDatabase(ctx context.Context, conn *database.Bun) error {
	if err := createTypes(ctx, conn); err != nil {
		return errors.Wrap(err, "creating custom types")
	}

	if err := database.CreateTables(ctx, conn, models.Models...); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return err
	}

	if err := database.MakeComments(ctx, conn, models.Models...); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return errors.Wrap(err, "make comments")
	}

	return createIndices(ctx, conn)
}

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

		// Namespace
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Namespace)(nil)).
			Index("namespace_idx").
			Column("namespace_id").
			Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

func createTypes(ctx context.Context, conn *database.Bun) error {
	log.Info().Msg("creating custom types...")
	return conn.DB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(
			ctx,
			`DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'msg_type') THEN
					CREATE TYPE msg_type AS ENUM ('PayForBlobs', 'CreatePeriodicVestingAccount', 'CreateVestingAccount', 'Send', 'Unjail', 'Undelegate', 'Delegate', 'CreateValidator', 'BeginRedelegate', 'EditValidator', 'WithdrawDelegatorReward', 'WithdrawValidatorCommission', 'Unknown');
				END IF;
			END$$;`,
		); err != nil {
			return err
		}
		if _, err := tx.ExecContext(
			ctx,
			`DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_type') THEN
					CREATE TYPE event_type AS ENUM ('coin_received', 'coinbase', 'coin_spent', 'burn', 'mint', 'message', 'proposer_reward', 'rewards', 'commission', 'liveness', 'attestation_request', 'transfer', 'pay_for_blobs', 'redelegate', 'withdraw_rewards', 'withdraw_commission', 'create_validator', 'delegate', 'edit_validator', 'unbond', 'tx', 'unknown');
				END IF;
			END$$;`,
		); err != nil {
			return err
		}
		if _, err := tx.ExecContext(
			ctx,
			`DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
					CREATE TYPE status AS ENUM ('success', 'failed');
				END IF;
			END$$;`,
		); err != nil {
			return err
		}
		return nil
	})
}
