// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	"github.com/vmihailenco/msgpack/v5"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

type Transaction struct {
	storage.Transaction
}

func BeginTransaction(ctx context.Context, tx storage.Transactable) (models.Transaction, error) {
	t, err := tx.BeginTransaction(ctx)
	return Transaction{t}, err
}

func (tx Transaction) SaveConstants(ctx context.Context, constants ...models.Constant) error {
	if len(constants) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&constants).
		Column("module", "name", "value").
		On("CONFLICT (module, name) DO UPDATE").
		Set("value = EXCLUDED.value").
		Exec(ctx)
	return err
}

func (tx Transaction) UpdateConstants(ctx context.Context, constants ...models.Constant) error {
	if len(constants) == 0 {
		return nil
	}

	values := tx.Tx().NewValues(&constants)

	_, err := tx.Tx().NewUpdate().
		With("_data", values).
		Model((*models.Constant)(nil)).
		TableExpr("_data").
		Set("value = _data.value").
		Where("constant.module = _data.module").
		Where("constant.name = _data.name").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveTransactions(ctx context.Context, txs ...models.Tx) error {
	switch len(txs) {
	case 0:
		return nil
	case 1:
		return tx.Add(ctx, &txs[0])
	default:
		arr := make([]any, len(txs))
		for i := range txs {
			arr[i] = &txs[i]
		}
		return tx.BulkSave(ctx, arr)
	}
}

type addedNamespace struct {
	bun.BaseModel `bun:"namespace"`
	*models.Namespace

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveNamespaces(ctx context.Context, namespaces ...*models.Namespace) (int64, error) {
	if len(namespaces) == 0 {
		return 0, nil
	}

	addedNamespaces := make([]addedNamespace, len(namespaces))
	for i := range namespaces {
		addedNamespaces[i].Namespace = namespaces[i]
	}

	_, err := tx.Tx().NewInsert().Model(&addedNamespaces).
		Column("version", "namespace_id", "pfb_count", "size", "first_height", "last_height", "last_message_time", "blobs_count").
		On("CONFLICT ON CONSTRAINT namespace_id_version_idx DO UPDATE").
		Set("size = EXCLUDED.size + added_namespace.size").
		Set("pfb_count = EXCLUDED.pfb_count + added_namespace.pfb_count").
		Set("last_height = EXCLUDED.last_height").
		Set("last_message_time = EXCLUDED.last_message_time").
		Set("blobs_count = EXCLUDED.blobs_count + added_namespace.blobs_count").
		Returning("xmax, id").
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	for i := range addedNamespaces {
		if addedNamespaces[i].Xmax == 0 {
			count++
		}
	}

	return count, err
}

type addedAddress struct {
	bun.BaseModel `bun:"address"`
	*models.Address

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveAddresses(ctx context.Context, addresses ...*models.Address) (int64, error) {
	if len(addresses) == 0 {
		return 0, nil
	}

	addr := make([]addedAddress, len(addresses))
	for i := range addresses {
		addr[i].Address = addresses[i]
	}

	_, err := tx.Tx().NewInsert().Model(&addr).
		Column("address", "height", "last_height", "hash", "name").
		On("CONFLICT ON CONSTRAINT address_idx DO UPDATE").
		Set("last_height = EXCLUDED.last_height").
		Returning("xmax, id").
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	for i := range addr {
		if addr[i].Xmax == 0 {
			count++
		}
	}

	return count, err
}

func (tx Transaction) SaveBalances(ctx context.Context, balances ...models.Balance) error {
	if len(balances) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&balances).
		Column("id", "currency", "spendable", "delegated", "unbonding").
		On("CONFLICT (id, currency) DO UPDATE").
		Set("spendable = EXCLUDED.spendable + balance.spendable").
		Set("delegated = EXCLUDED.delegated + balance.delegated").
		Set("unbonding = EXCLUDED.unbonding + balance.unbonding").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveEvents(ctx context.Context, events ...models.Event) error {
	switch {
	case len(events) == 0:
		return nil
	case len(events) < 20:
		_, err := tx.Tx().NewInsert().Model(&events).Exec(ctx)
		return err
	default:
		stmt, err := tx.Tx().PrepareContext(ctx,
			pq.CopyIn("event", "height", "time", "position", "type", "tx_id", "data"),
		)
		if err != nil {
			return err
		}

		for i := range events {
			var s []byte
			if len(events[i].Data) > 0 {
				if raw, err := msgpack.Marshal(events[i].Data); err == nil {
					s = raw
				}
			}

			if _, err := stmt.ExecContext(ctx, events[i].Height, events[i].Time, events[i].Position, events[i].Type, events[i].TxId, s); err != nil {
				return err
			}
		}

		if _, err := stmt.ExecContext(ctx); err != nil {
			return err
		}

		return stmt.Close()
	}
}

func (tx Transaction) SaveMessages(ctx context.Context, msgs ...*models.Message) error {
	if len(msgs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&msgs).Returning("id").Exec(ctx)
	return err
}

func (tx Transaction) SaveSigners(ctx context.Context, addresses ...models.Signer) error {
	if len(addresses) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&addresses).Exec(ctx)
	return err
}

func (tx Transaction) SaveBlobLogs(ctx context.Context, logs ...models.BlobLog) error {
	if len(logs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&logs).Exec(ctx)
	return err
}

func (tx Transaction) SaveMsgAddresses(ctx context.Context, addresses ...models.MsgAddress) error {
	if len(addresses) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&addresses).Exec(ctx)
	return err
}

func (tx Transaction) SaveMsgValidator(ctx context.Context, valMsgs ...models.MsgValidator) error {
	if len(valMsgs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&valMsgs).Exec(ctx)
	return err
}

func (tx Transaction) SaveNamespaceMessage(ctx context.Context, nsMsgs ...models.NamespaceMessage) error {
	if len(nsMsgs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&nsMsgs).Exec(ctx)
	return err
}

func (tx Transaction) SaveBlockSignatures(ctx context.Context, signs ...models.BlockSignature) error {
	if len(signs) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&signs).Exec(ctx)
	return err
}

func (tx Transaction) SaveVestingAccounts(ctx context.Context, accs ...*models.VestingAccount) error {
	if len(accs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&accs).Returning("id").Exec(ctx)
	return err
}

func (tx Transaction) SaveVestingPeriods(ctx context.Context, periods ...models.VestingPeriod) error {
	if len(periods) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&periods).Returning("id").Exec(ctx)
	return err
}

func (tx Transaction) SaveGrants(ctx context.Context, grants ...models.Grant) error {
	if len(grants) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().
		Model(&grants).
		Column("height", "time", "granter_id", "grantee_id", "authorization", "expiration", "revoked", "revoke_height", "params").
		On("CONFLICT ON CONSTRAINT grant_key DO UPDATE").
		Set("revoked = EXCLUDED.revoked").
		Set("revoke_height = EXCLUDED.revoke_height").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveSignals(ctx context.Context, signals ...*models.SignalVersion) error {
	if len(signals) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&signals).Exec(ctx)
	return err
}

func (tx Transaction) SaveUpgrades(ctx context.Context, upgrades ...*models.Upgrade) error {
	if len(upgrades) == 0 {
		return nil
	}

	for i := range upgrades {
		query := tx.Tx().NewInsert().Model(upgrades[i]).
			Column("version", "height", "time", "end_height", "end_time", "signer_id", "msg_id", "tx_id", "voting_power", "voted_power", "signals_count").
			On("CONFLICT (version) DO UPDATE")

		if upgrades[i].EndHeight > 0 {
			query = query.Set("end_height = EXCLUDED.end_height")
		}
		if !upgrades[i].EndTime.IsZero() {
			query = query.Set("end_time = EXCLUDED.end_time")
		}
		if upgrades[i].SignerId > 0 {
			query = query.Set("signer_id = EXCLUDED.signer_id")
		}
		if upgrades[i].MsgId > 0 {
			query = query.Set("msg_id = EXCLUDED.msg_id")
		}
		if upgrades[i].TxId > 0 {
			query = query.Set("tx_id = EXCLUDED.tx_id")
		}
		if !upgrades[i].VotingPower.IsZero() {
			query = query.Set("voting_power = EXCLUDED.voting_power")
		}
		if !upgrades[i].VotedPower.IsZero() {
			query = query.Set("voted_power = EXCLUDED.voted_power")
		}
		if upgrades[i].SignalsCount > 0 {
			query = query.Set("signals_count = EXCLUDED.signals_count + upgrade.signals_count")
		}
		if _, err := query.Exec(ctx); err != nil {
			return errors.Wrapf(err, "save upgrade %d", upgrades[i].Version)
		}
	}

	return nil
}

func (tx Transaction) SaveHyperlaneIgps(ctx context.Context, igps ...*models.HLIGP) error {
	if len(igps) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&igps).
		Column("id", "height", "time", "igp_id", "owner_id", "denom").
		On("CONFLICT (igp_id) DO UPDATE").
		Set("height = EXCLUDED.height").
		Set("time = EXCLUDED.time").
		Set("owner_id = EXCLUDED.owner_id").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveHyperlaneIgpConfigs(ctx context.Context, configs ...models.HLIGPConfig) error {
	if len(configs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&configs).
		Column("id", "height", "time", "gas_overhead", "gas_price", "remote_domain", "token_exchange_rate").
		On("CONFLICT (id, remote_domain) DO UPDATE").
		Set("height = EXCLUDED.height").
		Set("time = EXCLUDED.time").
		Set("gas_overhead = EXCLUDED.gas_overhead").
		Set("gas_price = EXCLUDED.gas_price").
		Set("token_exchange_rate = EXCLUDED.token_exchange_rate").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveHyperlaneGasPayments(ctx context.Context, payments ...*models.HLGasPayment) error {
	if len(payments) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&payments).Exec(ctx)
	return err
}

type addedValidator struct {
	bun.BaseModel `bun:"validator"`
	*models.Validator

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveValidators(ctx context.Context, validators ...*models.Validator) (int, error) {
	if len(validators) == 0 {
		return 0, nil
	}

	arr := make([]addedValidator, len(validators))
	for i := range validators {
		arr[i].Validator = validators[i]
	}

	query := tx.Tx().NewInsert().Model(&arr).
		Column("id", "delegator", "address", "cons_address", "moniker", "website", "identity", "contacts", "details", "rate", "max_rate", "max_change_rate", "min_self_delegation", "stake", "jailed", "commissions", "rewards", "height", "version", "messages_count", "creation_time").
		On("CONFLICT ON CONSTRAINT address_validator DO UPDATE").
		Set("rate = CASE WHEN EXCLUDED.rate > 0 THEN EXCLUDED.rate ELSE added_validator.rate END").
		Set("min_self_delegation = CASE WHEN EXCLUDED.min_self_delegation > 0 THEN EXCLUDED.min_self_delegation ELSE added_validator.min_self_delegation END").
		Set("stake = added_validator.stake + EXCLUDED.stake").
		Set("commissions = added_validator.commissions + EXCLUDED.commissions").
		Set("rewards = added_validator.rewards + EXCLUDED.rewards").
		Set("messages_count = added_validator.messages_count + EXCLUDED.messages_count").
		Set("moniker = CASE WHEN EXCLUDED.moniker != '[do-not-modify]' THEN EXCLUDED.moniker ELSE added_validator.moniker END").
		Set("website = CASE WHEN EXCLUDED.website != '[do-not-modify]' THEN EXCLUDED.website ELSE added_validator.website END").
		Set("identity = CASE WHEN EXCLUDED.identity != '[do-not-modify]' THEN EXCLUDED.identity ELSE added_validator.identity END").
		Set("contacts = CASE WHEN EXCLUDED.contacts != '[do-not-modify]' THEN EXCLUDED.contacts ELSE added_validator.contacts END").
		Set("details = CASE WHEN EXCLUDED.details != '[do-not-modify]' THEN EXCLUDED.details ELSE added_validator.details END").
		Set("jailed = CASE WHEN EXCLUDED.jailed IS NOT NULL THEN EXCLUDED.jailed ELSE added_validator.jailed END").
		Set("version = CASE WHEN EXCLUDED.version > 0 THEN EXCLUDED.version ELSE added_validator.version END").
		Returning("xmax, id")

	if _, err := query.Exec(ctx); err != nil {
		return 0, err
	}

	var count int
	for i := range arr {
		if arr[i].Xmax == 0 {
			count++
		}
	}

	return count, nil
}

func (tx Transaction) SaveUndelegations(ctx context.Context, undelegations ...models.Undelegation) error {
	if len(undelegations) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&undelegations).Exec(ctx)
	return err
}

func (tx Transaction) SaveRedelegations(ctx context.Context, redelegations ...models.Redelegation) error {
	if len(redelegations) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&redelegations).Exec(ctx)
	return err
}

func (tx Transaction) SaveStakingLogs(ctx context.Context, logs ...models.StakingLog) error {
	if len(logs) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&logs).Exec(ctx)
	return err
}

func (tx Transaction) SaveDelegations(ctx context.Context, delegations ...models.Delegation) error {
	if len(delegations) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&delegations).
		Column("id", "address_id", "validator_id", "amount").
		On("CONFLICT ON CONSTRAINT delegation_pair DO UPDATE").
		Set("amount = delegation.amount + EXCLUDED.amount").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveJails(ctx context.Context, jails ...models.Jail) error {
	if len(jails) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&jails).Exec(ctx)
	return err
}

func (tx Transaction) Jail(ctx context.Context, validators ...*models.Validator) error {
	if len(validators) == 0 {
		return nil
	}

	values := tx.Tx().NewValues(&validators)
	_, err := tx.Tx().NewUpdate().
		With("_data", values).
		Model((*models.Validator)(nil)).
		TableExpr("_data").
		Set("jailed = true").
		Set("stake = _data.stake + validator.stake").
		Where("validator.id = _data.id").
		Exec(ctx)
	return err
}

type addedProposal struct {
	bun.BaseModel `bun:"proposal"`
	*models.Proposal

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveProposals(ctx context.Context, proposals ...*models.Proposal) (int64, error) {
	if len(proposals) == 0 {
		return 0, nil
	}

	var count int64
	for i := range proposals {
		if proposals[i].Type == "" {
			proposals[i].Type = storageTypes.ProposalTypeText
		}
		if proposals[i].Status == "" {
			proposals[i].Status = storageTypes.ProposalStatusInactive
		}

		add := addedProposal{
			Proposal: proposals[i],
		}

		query := tx.Tx().NewInsert().
			Column("id", "proposer_id", "height", "created_at", "deposit_time", "activation_time", "status", "type", "title", "description", "deposit", "metadata", "changes", "yes", "no", "no_with_veto", "abstain").
			Column("yes_vals", "no_vals", "no_with_veto_vals", "abstain_vals", "yes_addrs", "no_addrs", "no_with_veto_addrs", "abstain_addrs", "votes_count", "voting_power", "yes_voting_power", "no_voting_power", "no_with_veto_voting_power", "abstain_voting_power").
			Column("total_voting_power", "quorum", "veto_quorum", "threshold", "min_deposit", "end_time", "error").
			Model(&add).
			On("CONFLICT (id) DO UPDATE").
			Set("votes_count = added_proposal.votes_count + EXCLUDED.votes_count").
			Set("yes = added_proposal.yes + EXCLUDED.yes").
			Set("no = added_proposal.no + EXCLUDED.no").
			Set("no_with_veto = added_proposal.no_with_veto + EXCLUDED.no_with_veto").
			Set("abstain = added_proposal.abstain + EXCLUDED.abstain").
			Set("yes_vals = added_proposal.yes_vals + EXCLUDED.yes_vals").
			Set("no_vals = added_proposal.no_vals + EXCLUDED.no_vals").
			Set("no_with_veto_vals = added_proposal.no_with_veto_vals + EXCLUDED.no_with_veto_vals").
			Set("abstain_vals = added_proposal.abstain_vals + EXCLUDED.abstain_vals").
			Set("yes_addrs = added_proposal.yes_addrs + EXCLUDED.yes_addrs").
			Set("no_addrs = added_proposal.no_addrs+ EXCLUDED.no_addrs").
			Set("no_with_veto_addrs = added_proposal.no_with_veto_addrs + EXCLUDED.no_with_veto_addrs").
			Set("abstain_addrs = added_proposal.abstain_addrs + EXCLUDED.abstain_addrs")

		if proposals[i].Deposit.IsPositive() {
			query.Set("deposit = added_proposal.deposit + EXCLUDED.deposit")
		}

		if !proposals[i].EmptyStatus() {
			query.Set("status = EXCLUDED.status")
		}

		if proposals[i].ActivationTime != nil {
			query.Set("activation_time = EXCLUDED.activation_time")
		}

		if proposals[i].EndTime != nil {
			query.Set("end_time = EXCLUDED.end_time")
		}

		if proposals[i].VotingPower.IsPositive() {
			query.Set("voting_power = EXCLUDED.voting_power")
		}
		if proposals[i].TotalVotingPower.IsPositive() {
			query.Set("total_voting_power = EXCLUDED.total_voting_power")
		}

		if proposals[i].Quorum != "" {
			query.Set("quorum = EXCLUDED.quorum")
		}
		if proposals[i].VetoQuorum != "" {
			query.Set("veto_quorum = EXCLUDED.veto_quorum")
		}
		if proposals[i].Threshold != "" {
			query.Set("threshold = EXCLUDED.threshold")
		}
		if proposals[i].MinDeposit != "" {
			query.Set("min_deposit = EXCLUDED.min_deposit")
		}
		if proposals[i].Error != "" {
			query.Set("error = EXCLUDED.error")
		}

		if proposals[i].YesVotingPower.IsPositive() {
			query.Set("yes_voting_power = EXCLUDED.yes_voting_power")
		}
		if proposals[i].NoVotingPower.IsPositive() {
			query.Set("no_voting_power = EXCLUDED.no_voting_power")
		}
		if proposals[i].NoWithVetoVotingPower.IsPositive() {
			query.Set("no_with_veto_voting_power = EXCLUDED.no_with_veto_voting_power")
		}
		if proposals[i].AbstainVotingPower.IsPositive() {
			query.Set("abstain_voting_power = EXCLUDED.abstain_voting_power")
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

var one = decimal.NewFromInt(1)

func (tx Transaction) SaveVotes(ctx context.Context, votes ...*models.Vote) (map[uint64]*models.VotesCount, error) {
	if len(votes) == 0 {
		return nil, nil
	}

	var votesCount = make(map[uint64]*models.VotesCount)
	for i := range votes {
		var existsVotes []models.Vote
		query := tx.Tx().NewSelect().
			Model(&existsVotes).
			Where("proposal_id = ?", votes[i].ProposalId)

		if votes[i].VoterId > 0 {
			query.Where("voter_id = ?", votes[i].VoterId)
		}
		if votes[i].ValidatorId != nil {
			query.Where("validator_id = ?", *votes[i].ValidatorId)
		}
		if err := query.Scan(ctx); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "receive existing votes")
		}
		ids := make([]uint64, len(existsVotes))
		total := votes[i].Weight.Copy()
		for j, vote := range existsVotes {
			total = total.Add(vote.Weight)
			ids[j] = vote.Id
		}
		if total.GreaterThan(one) {
			if _, err := tx.Tx().NewDelete().Model((*models.Vote)(nil)).Where("id IN (?)", bun.In(ids)).Exec(ctx); err != nil {
				return nil, errors.Wrap(err, "remove existing votes")
			}

			for _, vote := range existsVotes {
				if vc, ok := votesCount[vote.ProposalId]; ok {
					vc.Update(-1, vote)
				} else {
					var vc models.VotesCount
					vc.Update(-1, vote)
					votesCount[vote.ProposalId] = &vc
				}
			}
		}

		if vc, ok := votesCount[votes[i].ProposalId]; ok {
			vc.Update(1, *votes[i])
		} else {
			var vc models.VotesCount
			vc.Update(1, *votes[i])
			votesCount[votes[i].ProposalId] = &vc
		}
	}

	_, err := tx.Tx().NewInsert().Model(&votes).Exec(ctx)
	return votesCount, err
}

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
		if clients[i].LatestRevisionHeight > 0 {
			query.Set("latest_revision_height = EXCLUDED.latest_revision_height")
		}
		if clients[i].LatestRevisionNumber > 0 {
			query.Set("latest_revision_number = EXCLUDED.latest_revision_number")
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

		if _, err := query.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (tx Transaction) SaveIbcTransfers(ctx context.Context, transfers ...*models.IbcTransfer) error {
	if len(transfers) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&transfers).Exec(ctx)
	return err
}

func (tx Transaction) SaveHyperlaneMailbox(ctx context.Context, mailbox ...*models.HLMailbox) error {
	if len(mailbox) == 0 {
		return nil
	}

	for i := range mailbox {
		query := tx.Tx().NewInsert().
			Model(mailbox[i]).
			Column("height", "time", "tx_id", "mailbox", "internal_id", "owner_id", "default_ism", "default_hook", "required_hook", "domain", "sent_messages", "received_messages").
			On("CONFLICT (internal_id) DO UPDATE")

		if mailbox[i].Owner != nil {
			query.Set("owner_id = EXCLUDED.owner_id")
		}
		if mailbox[i].SentMessages > 0 {
			query.Set("sent_messages = hl_mailbox.sent_messages + EXCLUDED.sent_messages")
		}
		if mailbox[i].ReceivedMessages > 0 {
			query.Set("received_messages = hl_mailbox.received_messages + EXCLUDED.received_messages")
		}

		if _, err := query.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (tx Transaction) SaveHyperlaneTokens(ctx context.Context, tokens ...*models.HLToken) error {
	if len(tokens) == 0 {
		return nil
	}

	for i := range tokens {
		query := tx.Tx().NewInsert().
			Model(tokens[i]).
			Column("height", "time", "tx_id", "mailbox_id", "owner_id", "type", "denom", "token_id", "sent_transfers", "received_transfers", "sent", "received").
			On("CONFLICT (token_id) DO UPDATE")

		if tokens[i].Owner != nil {
			query.Set("owner_id = EXCLUDED.owner_id")
		}
		if tokens[i].SentTransfers > 0 {
			query.Set("sent_transfers = hl_token.sent_transfers + EXCLUDED.sent_transfers")
		}
		if tokens[i].ReceiveTransfers > 0 {
			query.Set("received_transfers = hl_token.received_transfers + EXCLUDED.received_transfers")
		}
		if !tokens[i].Sent.IsZero() {
			query.Set("sent = hl_token.sent + EXCLUDED.sent")
		}
		if !tokens[i].Received.IsZero() {
			query.Set("received = hl_token.received + EXCLUDED.received")
		}

		if _, err := query.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (tx Transaction) SaveHyperlaneTransfers(ctx context.Context, transfers ...*models.HLTransfer) error {
	if len(transfers) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&transfers).Exec(ctx)
	return err
}

func (tx Transaction) UpdateSlashedDelegations(ctx context.Context, validatorId uint64, burned decimal.Decimal) (balances []models.Balance, err error) {
	if validatorId == 0 || !burned.IsPositive() {
		return nil, nil
	}

	totalQuery := tx.Tx().NewSelect().
		Model((*models.Delegation)(nil)).
		ColumnExpr("sum(amount) as amount").
		Where("validator_id = ?", validatorId)

	burnedParts := tx.Tx().NewSelect().
		Table("total", "delegation").
		ColumnExpr("(delegation.amount * ? / total.amount) as amount", burned.String()).
		ColumnExpr("delegation.address_id as address_id").
		Where("validator_id = ?", validatorId)

	_, err = tx.Tx().NewUpdate().
		With("total", totalQuery).
		With("burned", burnedParts).
		Model((*models.Delegation)(nil)).
		TableExpr("burned").
		Set("amount = delegation.amount * burned.amount").
		Where("validator_id = ?", validatorId).
		Where("burned.address_id = delegation.address_id").
		Returning("delegation.address_id as id, 'utia' as currency, -burned.amount as delegated").
		Exec(ctx, &balances)
	return
}

func (tx Transaction) LastBlock(ctx context.Context) (block models.Block, err error) {
	err = tx.Tx().NewSelect().Model(&block).Order("id desc").Limit(1).Scan(ctx)
	return
}

func (tx Transaction) State(ctx context.Context, name string) (state models.State, err error) {
	err = tx.Tx().NewSelect().Model(&state).Where("name = ?", name).Scan(ctx)
	return
}

func (tx Transaction) Namespace(ctx context.Context, id uint64) (ns models.Namespace, err error) {
	err = tx.Tx().NewSelect().Model(&ns).Where("id = ?", id).Scan(ctx)
	return
}

func (tx Transaction) RollbackBlock(ctx context.Context, height types.Level) error {
	_, err := tx.Tx().NewDelete().
		Model((*models.Block)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return err
}

func (tx Transaction) RollbackBlockStats(ctx context.Context, height types.Level) (stats models.BlockStats, err error) {
	_, err = tx.Tx().NewDelete().Model(&stats).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackAddresses(ctx context.Context, height types.Level) (address []models.Address, err error) {
	_, err = tx.Tx().NewDelete().Model(&address).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackTxs(ctx context.Context, height types.Level) (txs []models.Tx, err error) {
	_, err = tx.Tx().NewDelete().Model(&txs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackEvents(ctx context.Context, height types.Level) (events []models.Event, err error) {
	_, err = tx.Tx().NewDelete().Model(&events).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackMessages(ctx context.Context, height types.Level) (msgs []models.Message, err error) {
	_, err = tx.Tx().NewDelete().Model(&msgs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackNamespaceMessages(ctx context.Context, height types.Level) (msgs []models.NamespaceMessage, err error) {
	_, err = tx.Tx().NewDelete().Model(&msgs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackNamespaces(ctx context.Context, height types.Level) (ns []models.Namespace, err error) {
	_, err = tx.Tx().NewDelete().Model(&ns).Where("first_height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackValidators(ctx context.Context, height types.Level) (validators []models.Validator, err error) {
	_, err = tx.Tx().NewDelete().Model(&validators).Where("height = ?", height).Returning("id").Exec(ctx)
	return
}

func (tx Transaction) RollbackBlobLog(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.BlobLog)(nil)).Where("height = ?", height).Exec(ctx)
	return
}

func (tx Transaction) RollbackGrants(ctx context.Context, height types.Level) (err error) {
	if _, err = tx.Tx().NewDelete().
		Model((*models.Grant)(nil)).
		Where("height = ?", height).
		Exec(ctx); err != nil {
		return err
	}

	_, err = tx.Tx().NewUpdate().
		Model((*models.Grant)(nil)).
		Where("revoke_height = ?", height).
		Set("revoked = false").
		Set("revoke_height = null").
		Exec(ctx)
	return
}

func (tx Transaction) RollbackUndelegations(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.Undelegation)(nil)).Where("height = ?", height).Exec(ctx)
	return
}

func (tx Transaction) RollbackRedelegations(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.Redelegation)(nil)).Where("height = ?", height).Exec(ctx)
	return
}

func (tx Transaction) RollbackVestingAccounts(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.VestingAccount)(nil)).Where("height = ?", height).Exec(ctx)
	return
}

func (tx Transaction) RollbackVestingPeriods(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.VestingPeriod)(nil)).Where("height = ?", height).Exec(ctx)
	return
}

func (tx Transaction) RollbackSigners(ctx context.Context, txIds []uint64) (err error) {
	_, err = tx.Tx().NewDelete().
		Model((*models.Signer)(nil)).
		Where("tx_id IN (?)", bun.In(txIds)).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackMessageAddresses(ctx context.Context, msgIds []uint64) (err error) {
	_, err = tx.Tx().NewDelete().
		Model((*models.MsgAddress)(nil)).
		Where("msg_id IN (?)", bun.In(msgIds)).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackMessageValidators(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.MsgValidator)(nil)).
		Where("height = ?", height).Exec(ctx)
	return
}

func (tx Transaction) RollbackBlockSignatures(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.BlockSignature)(nil)).
		Where("height = ?", height).Exec(ctx)
	return
}

func (tx Transaction) RollbackJails(ctx context.Context, height types.Level) (jails []models.Jail, err error) {
	_, err = tx.Tx().NewDelete().Model(&jails).
		Where("height = ?", height).
		Returning("id, validator_id").
		Exec(ctx)
	return
}

func (tx Transaction) RollbackStakingLogs(ctx context.Context, height types.Level) (logs []models.StakingLog, err error) {
	_, err = tx.Tx().NewDelete().Model(&logs).
		Where("height = ?", height).
		Returning("*").
		Exec(ctx)
	return
}

func (tx Transaction) RollbackProposals(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.Proposal)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackVotes(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.Vote)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackIbcClients(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.IbcClient)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackIbcConnections(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.IbcConnection)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackIbcChannels(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.IbcChannel)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackIbcTransfers(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.IbcChannel)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackHyperlaneMailbox(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.HLMailbox)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackHyperlaneTokens(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.HLToken)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackHyperlaneTransfers(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.HLTransfer)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackSignals(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.SignalVersion)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackUpgrades(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.Upgrade)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackHyperlaneIgps(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.HLIGP)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackHyperlaneIgpConfigs(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.HLIGPConfig)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) RollbackHyperlaneGasPayment(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.HLGasPayment)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return
}

func (tx Transaction) DeleteBalances(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := tx.Tx().NewDelete().
		Model((*models.Balance)(nil)).
		Where("id IN (?)", bun.In(ids)).
		Exec(ctx)
	return err
}

func (tx Transaction) DeleteDelegationsByValidator(ctx context.Context, ids ...uint64) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := tx.Tx().NewDelete().
		Model((*models.Delegation)(nil)).
		Where("validator_id IN (?)", bun.In(ids)).
		Exec(ctx)
	return err
}

func (tx Transaction) LastAddressAction(ctx context.Context, address []byte) (uint64, error) {
	var height uint64
	err := tx.Tx().NewSelect().
		Model((*models.MsgAddress)(nil)).
		ExcludeColumn("msg_id", "address_id", "type").
		Where("address.hash = ?", address).
		Order("msg_id desc").
		Relation("Msg", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("height")
		}).
		Relation("Address", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ExcludeColumn("*")
		}).
		Scan(ctx, &height)
	return height, err
}

func (tx Transaction) LastNamespaceMessage(ctx context.Context, nsId uint64) (msg models.NamespaceMessage, err error) {
	err = tx.Tx().NewSelect().
		Model(&msg).
		Where("namespace_id = ?", nsId).
		Order("msg_id desc").
		Limit(1).
		Scan(ctx)
	return
}

func (tx Transaction) GetProposerId(ctx context.Context, address string) (id uint64, err error) {
	err = tx.Tx().NewSelect().
		Model((*models.Validator)(nil)).
		Column("id").
		Where("cons_address = ?", address).
		Order("id desc").
		Limit(1).
		Scan(ctx, &id)
	return
}

func (tx Transaction) SaveRollup(ctx context.Context, rollup *models.Rollup) error {
	if rollup == nil {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(rollup).Exec(ctx)
	return err
}

func (tx Transaction) UpdateRollup(ctx context.Context, rollup *models.Rollup) error {
	if rollup == nil || (rollup.IsEmpty() && !rollup.Verified) {
		return nil
	}

	query := tx.Tx().NewUpdate().Model(rollup).WherePK()

	if rollup.Name != "" {
		query = query.Set("name = ?", rollup.Name)
	}
	if rollup.Slug != "" {
		query = query.Set("slug = ?", rollup.Slug)
	}
	if rollup.Description != "" {
		query = query.Set("description = ?", rollup.Description)
	}
	if rollup.Twitter != "" {
		query = query.Set("twitter = ?", rollup.Twitter)
	}
	if rollup.GitHub != "" {
		query = query.Set("github = ?", rollup.GitHub)
	}
	if rollup.Website != "" {
		query = query.Set("website = ?", rollup.Website)
	}
	if rollup.Logo != "" {
		query = query.Set("logo = ?", rollup.Logo)
	}
	if rollup.L2Beat != "" {
		query = query.Set("l2_beat = ?", rollup.L2Beat)
	}
	if rollup.Explorer != "" {
		query = query.Set("explorer = ?", rollup.Explorer)
	}
	if rollup.BridgeContract != "" {
		query = query.Set("bridge_contract = ?", rollup.BridgeContract)
	}
	if rollup.Stack != "" {
		query = query.Set("stack = ?", rollup.Stack)
	}
	if rollup.Links != nil {
		query = query.Set("links = ?", pq.Array(rollup.Links))
	}
	if rollup.Type != "" {
		query = query.Set("type = ?", rollup.Type)
	}
	if rollup.Category != "" {
		query = query.Set("category = ?", rollup.Category)
	}
	if rollup.Tags != nil {
		query = query.Set("tags = ?", pq.Array(rollup.Tags))
	}
	if rollup.Provider != "" {
		query = query.Set("provider = ?", rollup.Provider)
	}
	if rollup.Compression != "" {
		query = query.Set("compression = ?", rollup.Compression)
	}
	if rollup.VM != "" {
		query = query.Set("vm = ?", rollup.VM)
	}
	if rollup.DeFiLama != "" {
		query = query.Set("defi_lama = ?", rollup.DeFiLama)
	}
	if rollup.SettledOn != "" {
		query = query.Set("settled_on = ?", rollup.SettledOn)
	}
	if rollup.Color != "" {
		query = query.Set("color = ?", rollup.Color)
	}

	query = query.Set("verified = ?", rollup.Verified)

	_, err := query.Exec(ctx)
	return err
}

func (tx Transaction) SaveProviders(ctx context.Context, providers ...models.RollupProvider) error {
	if len(providers) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&providers).Exec(ctx)
	return err
}

func (tx Transaction) DeleteProviders(ctx context.Context, rollupId uint64) error {
	if rollupId == 0 {
		return nil
	}
	_, err := tx.Tx().NewDelete().
		Model((*models.RollupProvider)(nil)).
		Where("rollup_id = ?", rollupId).
		Exec(ctx)
	return err
}

func (tx Transaction) DeleteRollup(ctx context.Context, rollupId uint64) error {
	if rollupId == 0 {
		return nil
	}
	_, err := tx.Tx().NewDelete().
		Model((*models.Rollup)(nil)).
		Where("id = ?", rollupId).
		Exec(ctx)
	return err
}

func (tx Transaction) RetentionBlockSignatures(ctx context.Context, height types.Level) error {
	_, err := tx.Tx().NewDelete().Model((*models.BlockSignature)(nil)).
		Where("height <= ?", height).
		Exec(ctx)
	return err
}

func (tx Transaction) CancelUnbondings(ctx context.Context, cancellations ...models.Undelegation) error {
	if len(cancellations) == 0 {
		return nil
	}

	for i := range cancellations {
		if _, err := tx.Tx().NewDelete().
			Model(&cancellations[i]).
			Where("height = ?height").
			Where("amount = ?amount").
			Where("validator_id = ?validator_id").
			Where("address_id = ?address_id").
			Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (tx Transaction) RetentionCompletedUnbondings(ctx context.Context, blockTime time.Time) error {
	_, err := tx.Tx().NewDelete().Model((*models.Undelegation)(nil)).
		Where("completion_time < ?", blockTime).
		Exec(ctx)
	return err
}

func (tx Transaction) RetentionCompletedRedelegations(ctx context.Context, blockTime time.Time) error {
	_, err := tx.Tx().NewDelete().Model((*models.Redelegation)(nil)).
		Where("completion_time < ?", blockTime).
		Exec(ctx)
	return err
}

func (tx Transaction) UpdateValidators(ctx context.Context, validators ...*models.Validator) error {
	if len(validators) == 0 {
		return nil
	}

	values := tx.Tx().NewValues(&validators)

	_, err := tx.Tx().NewUpdate().
		With("_data", values).
		Model((*models.Validator)(nil)).
		TableExpr("_data").
		Set("stake = validator.stake + _data.stake").
		Set("jailed = _data.jailed").
		Set("commissions = validator.commissions + _data.commissions").
		Set("rewards = validator.rewards + _data.rewards").
		Where("validator.id = _data.id").
		Exec(ctx)
	return err
}

func (tx Transaction) Validator(ctx context.Context, id uint64) (val models.Validator, err error) {
	err = tx.Tx().NewSelect().Model(&val).Where("id = ?", id).Scan(ctx)
	return
}

func (tx Transaction) Delegation(ctx context.Context, validatorId, addressId uint64) (val models.Delegation, err error) {
	err = tx.Tx().NewSelect().Model(&val).
		Where("validator_id = ?", validatorId).
		Where("address_id = ?", addressId).
		Scan(ctx)
	return
}

func (tx Transaction) RefreshLeaderboard(ctx context.Context) error {
	_, err := tx.Tx().ExecContext(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	return err
}

func (tx Transaction) ActiveProposals(ctx context.Context) (proposals []models.Proposal, err error) {
	err = tx.Tx().NewSelect().Model(&proposals).
		Where("status = ?", storageTypes.ProposalStatusActive).
		Scan(ctx)
	return
}

func (tx Transaction) BondedValidators(ctx context.Context, limit int) (validators []models.Validator, err error) {
	err = tx.Tx().NewSelect().Model(&validators).
		Column("id", "stake", "version").
		OrderExpr("stake desc").
		Limit(limit).
		Scan(ctx)
	return
}

func (tx Transaction) ProposalVotes(ctx context.Context, proposalId uint64, limit, offset int) (votes []models.Vote, err error) {
	query := tx.Tx().NewSelect().Model(&votes).
		Where("proposal_id = ?", proposalId).
		OrderExpr("id asc")

	if limit < 1 {
		limit = 10
	}
	query = query.Limit(limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Scan(ctx)
	return
}

func (tx Transaction) AddressDelegations(ctx context.Context, addressId uint64) (val []models.Delegation, err error) {
	err = tx.Tx().NewSelect().Model(&val).
		Where("address_id = ?", addressId).
		Scan(ctx)
	return
}

func (tx Transaction) Proposal(ctx context.Context, id uint64) (proposal models.Proposal, err error) {
	err = tx.Tx().NewSelect().Model(&proposal).
		Where("id = ?", id).
		Column("id", "changes", "type").
		Scan(ctx)
	return
}

func (tx Transaction) IbcConnection(ctx context.Context, id string) (conn models.IbcConnection, err error) {
	err = tx.Tx().NewSelect().Model(&conn).
		Where("connection_id = ?", id).
		Column("client_id").
		Scan(ctx)
	return
}

func (tx Transaction) HyperlaneMailbox(ctx context.Context, internalId uint64) (mailbox models.HLMailbox, err error) {
	err = tx.Tx().NewSelect().Model(&mailbox).
		Where("internal_id = ?", internalId).
		Column("id").
		Scan(ctx)
	return
}

func (tx Transaction) HyperlaneToken(ctx context.Context, id []byte) (token models.HLToken, err error) {
	err = tx.Tx().NewSelect().Model(&token).
		Where("token_id = ?", id).
		Column("id").
		Scan(ctx)
	return
}

func (tx Transaction) UpdateSignalsAfterUpgrade(ctx context.Context, version uint64) error {
	_, err := tx.Tx().NewUpdate().Table("signal_version", "validator").
		SetColumn("voting_power", "validator.stake").
		Where("signal_version.version = ?", version).
		Where("validator.id = validator_id").
		Exec(ctx)
	return err
}

func (tx Transaction) HyperlaneIgp(ctx context.Context, id []byte) (igp models.HLIGP, err error) {
	err = tx.Tx().NewSelect().Model(&igp).
		Where("igp_id = ?", id).
		Scan(ctx)
	return
}

func (tx Transaction) HyperlaneIgpConfig(ctx context.Context, id uint64) (config models.HLIGPConfig, err error) {
	err = tx.Tx().NewSelect().Model(&config).
		Where("id = ?", id).
		Scan(ctx)
	return
}
