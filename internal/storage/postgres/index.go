// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	celestialPg "github.com/celenium-io/celestial-module/pkg/storage/postgres"
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
			Model((*storage.Block)(nil)).
			Index("block_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Block)(nil)).
			Index("block_proposer_id_idx").
			Column("proposer_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Block)(nil)).
			Index("block_hash_idx").
			Column("hash").
			Using("HASH").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Block)(nil)).
			Index("block_data_hash_idx").
			Column("data_hash").
			Using("HASH").
			Exec(ctx); err != nil {
			return err
		}

		// BlockStats
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.BlockStats)(nil)).
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

		// Signer
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Signer)(nil)).
			Index("signer_tx_id_idx").
			Column("tx_id").
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
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Message)(nil)).
			Index("message_type_idx").
			Column("type").
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
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Namespace)(nil)).
			Index("namespace_last_action_idx").
			Column("last_message_time").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Namespace)(nil)).
			Index("namespace_pfb_count_idx").
			Column("pfb_count").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Namespace)(nil)).
			Index("namespace_size_idx").
			Column("size").
			Exec(ctx); err != nil {
			return err
		}

		// Validator
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Validator)(nil)).
			Index("validator_cons_address_idx").
			Column("cons_address").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Validator)(nil)).
			Index("validator_moniker_idx").
			ColumnExpr("moniker gin_trgm_ops").
			Using("GIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Validator)(nil)).
			Index("validator_jailed_idx").
			ColumnExpr("(1)").
			Where("NOT jailed").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Validator)(nil)).
			Index("validator_version_idx").
			Column("version").
			Exec(ctx); err != nil {
			return err
		}

		// Blob log
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.BlobLog)(nil)).
			Index("blob_log_namespace_id_idx").
			Column("namespace_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.BlobLog)(nil)).
			Index("blob_log_signer_id_idx").
			Column("signer_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.BlobLog)(nil)).
			Index("blob_log_tx_id_idx").
			Column("tx_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.BlobLog)(nil)).
			Index("blob_log_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.BlobLog)(nil)).
			Index("blob_log_commitment_idx").
			Column("commitment").
			Using("HASH").
			Exec(ctx); err != nil {
			return err
		}

		// Rollup
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Rollup)(nil)).
			Index("rollup_name_idx").
			ColumnExpr("name gin_trgm_ops").
			Using("GIN").
			Exec(ctx); err != nil {
			return err
		}

		// RollupProvider
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.RollupProvider)(nil)).
			Index("rollup_provider_rollup_id_idx").
			Column("rollup_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.RollupProvider)(nil)).
			Index("rollup_provider_namespace_id_idx").
			Column("namespace_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.RollupProvider)(nil)).
			Index("rollup_provider_address_id_idx").
			Column("address_id").
			Exec(ctx); err != nil {
			return err
		}

		// BlockSignature
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.BlockSignature)(nil)).
			Index("block_signature_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.BlockSignature)(nil)).
			Index("block_signature_validator_id_idx").
			Column("validator_id").
			Exec(ctx); err != nil {
			return err
		}

		// StakingLog
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.StakingLog)(nil)).
			Index("staking_log_validator_id_idx").
			Column("validator_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.StakingLog)(nil)).
			Index("staking_log_address_id_idx").
			Column("address_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.StakingLog)(nil)).
			Index("staking_log_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.StakingLog)(nil)).
			Index("staking_log_type_idx").
			Column("type").
			Exec(ctx); err != nil {
			return err
		}

		// Delegation
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Delegation)(nil)).
			Index("delegation_address_id_idx").
			Column("address_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Delegation)(nil)).
			Index("delegation_validator_id_idx").
			Column("validator_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Delegation)(nil)).
			Index("delegation_amount_idx").
			Column("amount").
			Exec(ctx); err != nil {
			return err
		}

		// Redelegation
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Redelegation)(nil)).
			Index("redelegation_address_id_idx").
			Column("address_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Redelegation)(nil)).
			Index("redelegation_amount_idx").
			Column("amount").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Redelegation)(nil)).
			Index("redelegation_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}

		// Undelegation
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Undelegation)(nil)).
			Index("undelegation_address_id_idx").
			Column("address_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Undelegation)(nil)).
			Index("undelegation_validator_id_idx").
			Column("validator_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Undelegation)(nil)).
			Index("undelegation_amount_idx").
			Column("amount").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Undelegation)(nil)).
			Index("undelegation_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}

		// Jail
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Jail)(nil)).
			Index("jail_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Jail)(nil)).
			Index("jail_validator_id_idx").
			Column("validator_id").
			Exec(ctx); err != nil {
			return err
		}

		// Vesting account
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.VestingAccount)(nil)).
			Index("vesting_account_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.VestingAccount)(nil)).
			Index("vesting_account_address_id_idx").
			Column("address_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.VestingAccount)(nil)).
			Index("vesting_account_end_time_idx").
			Column("end_time").
			Exec(ctx); err != nil {
			return err
		}

		// Grant
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Grant)(nil)).
			Index("grant_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Grant)(nil)).
			Index("grant_granter_id_idx").
			Column("granter_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Grant)(nil)).
			Index("grant_grantee_id_idx").
			Column("grantee_id").
			Exec(ctx); err != nil {
			return err
		}

		// Celestial
		if err := celestialPg.CreateIndex(ctx, tx); err != nil {
			return err
		}

		// Proposal
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Proposal)(nil)).
			Index("proposal_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Proposal)(nil)).
			Index("proposal_proposer_id_idx").
			Column("proposer_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Proposal)(nil)).
			Index("proposal_status_idx").
			Column("status").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Proposal)(nil)).
			Index("proposal_type_idx").
			Column("type").
			Exec(ctx); err != nil {
			return err
		}

		// Vote
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Vote)(nil)).
			Index("vote_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Vote)(nil)).
			Index("vote_proposal_id_idx").
			Column("proposal_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Vote)(nil)).
			Index("vote_voter_id_idx").
			Column("voter_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Vote)(nil)).
			Index("vote_validator_id_idx").
			Column("validator_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Vote)(nil)).
			Index("vote_option_idx").
			Column("option").
			Exec(ctx); err != nil {
			return err
		}

		// IBC Client
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcClient)(nil)).
			Index("ibc_client_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcClient)(nil)).
			Index("ibc_client_updated_at_idx").
			Column("updated_at").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcClient)(nil)).
			Index("ibc_client_chain_id_idx").
			Column("chain_id").
			Where("chain_id is not null").
			Exec(ctx); err != nil {
			return err
		}

		// IBC Connection
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcConnection)(nil)).
			Index("ibc_connection_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcConnection)(nil)).
			Index("ibc_connection_client_id_idx").
			Column("client_id").
			Exec(ctx); err != nil {
			return err
		}

		// IBC Channel
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcChannel)(nil)).
			Index("ibc_channel_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcChannel)(nil)).
			Index("ibc_channel_client_id_idx").
			Column("client_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcChannel)(nil)).
			Index("ibc_channel_connection_id_idx").
			Column("connection_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcChannel)(nil)).
			Index("ibc_channel_status_idx").
			Column("status").
			Exec(ctx); err != nil {
			return err
		}

		// IBC Transfer
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcTransfer)(nil)).
			Index("ibc_transfer_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcTransfer)(nil)).
			Index("ibc_transfer_receiver_id_idx").
			Column("receiver_id").
			Where("receiver_id is not null").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcTransfer)(nil)).
			Index("ibc_transfer_sender_id_idx").
			Column("sender_id").
			Where("sender_id is not null").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcTransfer)(nil)).
			Index("ibc_transfer_channel_id_idx").
			Column("channel_id").
			Where("channel_id is not null").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcTransfer)(nil)).
			Index("ibc_transfer_connection_id_idx").
			Column("connection_id").
			Where("connection_id is not null").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.IbcTransfer)(nil)).
			Index("ibc_transfer_tx_id_idx").
			Column("tx_id").
			Exec(ctx); err != nil {
			return err
		}

		// Hyperlane Gas Payment
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLGasPayment)(nil)).
			Index("hl_gas_payment_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLGasPayment)(nil)).
			Index("hl_gas_payment_transfer_id_idx").
			Column("transfer_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLGasPayment)(nil)).
			Index("hl_gas_payment_igp_id_idx").
			Column("igp_id").
			Where("igp_id is not null").
			Exec(ctx); err != nil {
			return err
		}

		// Hyperlane Mailbox
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLMailbox)(nil)).
			Index("hl_mailbox_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLMailbox)(nil)).
			Index("hl_mailbox_owner_id_idx").
			Column("owner_id").
			Where("owner_id is not null").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLMailbox)(nil)).
			Index("hl_mailbox_tx_id_idx").
			Column("tx_id").
			Exec(ctx); err != nil {
			return err
		}

		// Hyperlane Token
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLToken)(nil)).
			Index("hl_token_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLToken)(nil)).
			Index("hl_token_mailbox_id_idx").
			Column("mailbox_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLToken)(nil)).
			Index("hl_token_owner_id_idx").
			Column("owner_id").
			Where("owner_id is not null").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLToken)(nil)).
			Index("hl_token_type_idx").
			Column("type").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLToken)(nil)).
			Index("hl_token_tx_id_idx").
			Column("tx_id").
			Exec(ctx); err != nil {
			return err
		}

		// Hyperlane Transfer
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLTransfer)(nil)).
			Index("hl_transfer_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLTransfer)(nil)).
			Index("hl_transfer_mailbox_id_idx").
			Column("mailbox_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLTransfer)(nil)).
			Index("hl_transfer_token_id_idx").
			Column("token_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLTransfer)(nil)).
			Index("hl_transfer_relayer_id_idx").
			Column("relayer_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLTransfer)(nil)).
			Index("hl_transfer_tx_id_idx").
			Column("tx_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLTransfer)(nil)).
			Index("hl_transfer_address_id_idx").
			Column("address_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLTransfer)(nil)).
			Index("hl_transfer_counterparty_idx").
			Column("counterparty").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLTransfer)(nil)).
			Index("hl_transfer_type_idx").
			Column("type").
			Exec(ctx); err != nil {
			return err
		}

		// Signal Version
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLTransfer)(nil)).
			Index("signal_version_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.SignalVersion)(nil)).
			Index("signal_version_validator_id_idx").
			Column("validator_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.SignalVersion)(nil)).
			Index("signal_version_tx_id_idx").
			Column("tx_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.SignalVersion)(nil)).
			Index("signal_version_version_idx").
			Column("version").
			Exec(ctx); err != nil {
			return err
		}

		// Upgrade
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.HLTransfer)(nil)).
			Index("upgrade_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Upgrade)(nil)).
			Index("upgrade_signer_id_idx").
			Column("signer_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Upgrade)(nil)).
			Index("upgrade_tx_id_idx").
			Column("tx_id").
			Exec(ctx); err != nil {
			return err
		}

		// MsgValidator
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.MsgValidator)(nil)).
			Index("msg_validator_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}
