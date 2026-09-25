// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/uptrace/bun"
)

type addedIbcClient struct {
	bun.BaseModel `bun:"ibc_client"`
	*models.IbcClient

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveIbcClients(ctx context.Context, clients ...*models.IbcClient) (int64, error) {
	if len(clients) == 0 {
		return 0, nil
	}

	count := int64(0)

	for i := range clients {
		add := addedIbcClient{
			IbcClient: clients[i],
		}

		query := tx.Tx().NewInsert().
			Column("id", "created_at", "updated_at", "height", "tx_id", "creator_id", "latest_revision_height", "latest_revision_number", "frozen_revision_height", "frozen_revision_number", "type").
			Column("trusting_period", "unbonding_period", "max_clock_drift", "trust_level_denominator", "trust_level_numerator", "connection_count", "chain_id").
			Model(&add).
			On("CONFLICT (id) DO UPDATE")

		if clients[i].ConnectionCount > 0 {
			query.Set("connection_count = added_ibc_client.connection_count + EXCLUDED.connection_count")
		}
		if clients[i].TrustingPeriod > 0 {
			query.Set("trusting_period = EXCLUDED.trusting_period")
		}
		if clients[i].UnbondingPeriod > 0 {
			query.Set("unbonding_period = EXCLUDED.unbonding_period")
		}
		if clients[i].MaxClockDrift > 0 {
			query.Set("max_clock_drift = EXCLUDED.max_clock_drift")
		}
		if clients[i].TrustLevelDenominator > 0 {
			query.Set("trust_level_denominator = EXCLUDED.trust_level_denominator")
		}
		if clients[i].TrustLevelNumerator > 0 {
			query.Set("trust_level_numerator = EXCLUDED.trust_level_numerator")
		}
		if clients[i].LatestRevisionHeight > 0 || clients[i].LatestRevisionNumber > 0 {
			// latest height only grows, compared as (revision, height)
			const isHigher = "(EXCLUDED.latest_revision_number, EXCLUDED.latest_revision_height) > (added_ibc_client.latest_revision_number, added_ibc_client.latest_revision_height)"
			query.Set("latest_revision_height = CASE WHEN " + isHigher + " THEN EXCLUDED.latest_revision_height ELSE added_ibc_client.latest_revision_height END")
			query.Set("latest_revision_number = CASE WHEN " + isHigher + " THEN EXCLUDED.latest_revision_number ELSE added_ibc_client.latest_revision_number END")
		}
		if clients[i].FrozenRevisionHeight > 0 {
			query.Set("frozen_revision_height = EXCLUDED.frozen_revision_height")
		}
		if clients[i].FrozenRevisionNumber > 0 {
			query.Set("frozen_revision_number = EXCLUDED.frozen_revision_number")
		}
		if clients[i].ChainId != "" {
			query.Set("chain_id = EXCLUDED.chain_id")
		}
		if !clients[i].UpdatedAt.IsZero() {
			query.Set("updated_at = EXCLUDED.updated_at")
		}

		if _, err := query.Returning("xmax, id").Exec(ctx); err != nil {
			return 0, err
		}

		if add.Xmax == 0 {
			count++
		}
	}

	return count, nil
}

// RecoverIbcClient mirrors ibc-go RecoverClient: the subject is unfrozen and takes the substitute's
// latest height, chain id and trusting period. Without a known substitute it is only unfrozen.
func (tx Transaction) RecoverIbcClient(ctx context.Context, subjectId, substituteId string, updatedAt time.Time) error {
	if substituteId == "" {
		_, err := tx.Tx().NewUpdate().
			Model((*models.IbcClient)(nil)).
			Set("frozen_revision_height = 0").
			Set("frozen_revision_number = 0").
			Set("updated_at = ?", updatedAt).
			Where("id = ?", subjectId).
			Exec(ctx)
		return err
	}

	_, err := tx.Tx().NewUpdate().
		TableExpr("ibc_client AS subject").
		TableExpr("ibc_client AS substitute").
		Set("frozen_revision_height = 0").
		Set("frozen_revision_number = 0").
		Set("latest_revision_height = substitute.latest_revision_height").
		Set("latest_revision_number = substitute.latest_revision_number").
		// chain id is known only after the first update_client, solomachine has no trusting period
		Set("chain_id = COALESCE(NULLIF(substitute.chain_id, ''), subject.chain_id)").
		Set("trusting_period = CASE WHEN substitute.trusting_period > 0 THEN substitute.trusting_period ELSE subject.trusting_period END").
		Set("updated_at = ?", updatedAt).
		Where("subject.id = ?", subjectId).
		Where("substitute.id = ?", substituteId).
		Exec(ctx)
	return err
}

func (tx Transaction) SaveIbcConnections(ctx context.Context, conns ...*models.IbcConnection) error {
	if len(conns) == 0 {
		return nil
	}

	for i := range conns {
		query := tx.Tx().NewInsert().
			Model(conns[i]).
			Column("connection_id", "client_id", "counterparty_connection_id", "counterparty_client_id", "created_at", "connected_at", "height", "connection_height", "create_tx_id", "connection_tx_id", "channels_count").
			On("CONFLICT (connection_id) DO UPDATE")

		if conns[i].ChannelsCount != 0 {
			query.Set("channels_count = ibc_connection.channels_count + EXCLUDED.channels_count")
		}
		if !conns[i].ConnectedAt.IsZero() {
			query.Set("connected_at = EXCLUDED.connected_at")
		}
		if conns[i].ConnectionTxId > 0 {
			query.Set("connection_tx_id = EXCLUDED.connection_tx_id")
		}
		if conns[i].ConnectionHeight > 0 {
			query.Set("connection_height = EXCLUDED.connection_height")
		}
		if conns[i].CounterpartyConnectionId != "" {
			query.Set("counterparty_connection_id = EXCLUDED.counterparty_connection_id")
		}

		if _, err := query.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (tx Transaction) SaveIbcChannels(ctx context.Context, channels ...*models.IbcChannel) error {
	if len(channels) == 0 {
		return nil
	}

	for i := range channels {
		query := tx.Tx().NewInsert().
			Model(channels[i]).
			Column("id", "connection_id", "client_id", "port_id", "counterparty_port_id", "counterparty_channel_id", "version", "created_at", "confirmed_at", "height", "confirmation_height", "create_tx_id", "confirmation_tx_id", "ordering", "creator_id", "status", "received", "sent", "transfers_count").
			On("CONFLICT (id) DO UPDATE")

		if !channels[i].ConfirmedAt.IsZero() {
			query.Set("confirmed_at = EXCLUDED.confirmed_at")
		}
		if channels[i].ConfirmationTxId > 0 {
			query.Set("confirmation_tx_id = EXCLUDED.confirmation_tx_id")
		}
		if channels[i].ConfirmationHeight > 0 {
			query.Set("confirmation_height = EXCLUDED.confirmation_height")
		}
		if channels[i].CounterpartyChannelId != "" {
			query.Set("counterparty_channel_id = EXCLUDED.counterparty_channel_id")
		}
		if channels[i].Status == storageTypes.IbcChannelStatusClosed || channels[i].Status == storageTypes.IbcChannelStatusOpened {
			query.Set("status = EXCLUDED.status")
		}
		if !channels[i].Received.IsZero() {
			query.Set("received = ibc_channel.received + EXCLUDED.received")
		}
		if !channels[i].Sent.IsZero() {
			query.Set("sent = ibc_channel.sent + EXCLUDED.sent")
		}
		if channels[i].TransfersCount > 0 {
			query.Set("transfers_count = ibc_channel.transfers_count + EXCLUDED.transfers_count")
		}
		if channels[i].Version != "" {
			query.Set("version = EXCLUDED.version")
		}

		if _, err := query.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (tx Transaction) RollbackIbcClients(ctx context.Context, height types.Level) (int64, error) {
	return tx.rollback(ctx, height, (*models.IbcClient)(nil))
}

func (tx Transaction) RollbackIbcConnections(ctx context.Context, height types.Level) (int64, error) {
	return tx.rollback(ctx, height, (*models.IbcConnection)(nil))
}

func (tx Transaction) RollbackIbcChannels(ctx context.Context, height types.Level) (int64, error) {
	return tx.rollback(ctx, height, (*models.IbcChannel)(nil))
}
