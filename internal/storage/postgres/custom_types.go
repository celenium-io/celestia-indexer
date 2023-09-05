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
	msgType = `DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'msg_type') THEN
			CREATE TYPE msg_type AS ENUM ('MsgPayForBlobs', 'MsgCreatePeriodicVestingAccount', 'MsgCreateVestingAccount', 'MsgSend', 'MsgUnjail', 'MsgUndelegate', 'MsgDelegate', 'MsgCreateValidator', 'MsgBeginRedelegate', 'MsgEditValidator', 'MsgWithdrawDelegatorReward', 'MsgWithdrawValidatorCommission', 'MsgGrantAllowance', 'MsgUnknown');
		END IF;
	END$$;`

	eventType = `DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_type') THEN
			CREATE TYPE event_type AS ENUM ('coin_received', 'coinbase', 'coin_spent', 'burn', 'mint', 'message', 'proposer_reward', 'rewards', 'commission', 'liveness', 'AttestationRequest', 'transfer', 'celestia.blob.v1.EventPayForBlobs', 'redelegate', 'withdraw_rewards', 'withdraw_commission', 'create_validator', 'delegate', 'edit_validator', 'unbond', 'tx', 'use_feegrant', 'revoke_feegrant', 'set_feegrant', 'update_feegrant', 'unknown');
		END IF;
	END$$;`

	statusType = `DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
			CREATE TYPE status AS ENUM ('success', 'failed');
		END IF;
	END$$;`

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
		if _, err := tx.ExecContext(ctx, msgType); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, eventType); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, statusType); err != nil {
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
