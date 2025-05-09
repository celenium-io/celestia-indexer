// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/lib/pq"
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
		Column("address", "height", "last_height", "hash").
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
		Column("id", "delegator", "address", "cons_address", "moniker", "website", "identity", "contacts", "details", "rate", "max_rate", "max_change_rate", "min_self_delegation", "stake", "jailed", "commissions", "rewards", "height").
		On("CONFLICT ON CONSTRAINT address_validator DO UPDATE").
		Set("rate = CASE WHEN EXCLUDED.rate > 0 THEN EXCLUDED.rate ELSE added_validator.rate END").
		Set("min_self_delegation = CASE WHEN EXCLUDED.min_self_delegation > 0 THEN EXCLUDED.min_self_delegation ELSE added_validator.min_self_delegation END").
		Set("stake = added_validator.stake + EXCLUDED.stake").
		Set("commissions = added_validator.commissions + EXCLUDED.commissions").
		Set("rewards = added_validator.rewards + EXCLUDED.rewards").
		Set("moniker = CASE WHEN EXCLUDED.moniker != '[do-not-modify]' THEN EXCLUDED.moniker ELSE added_validator.moniker END").
		Set("website = CASE WHEN EXCLUDED.website != '[do-not-modify]' THEN EXCLUDED.website ELSE added_validator.website END").
		Set("identity = CASE WHEN EXCLUDED.identity != '[do-not-modify]' THEN EXCLUDED.identity ELSE added_validator.identity END").
		Set("contacts = CASE WHEN EXCLUDED.contacts != '[do-not-modify]' THEN EXCLUDED.contacts ELSE added_validator.contacts END").
		Set("details = CASE WHEN EXCLUDED.details != '[do-not-modify]' THEN EXCLUDED.details ELSE added_validator.details END").
		Set("jailed = CASE WHEN EXCLUDED.jailed IS NOT NULL THEN EXCLUDED.jailed ELSE added_validator.jailed END").
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
			Column("id", "proposer_id", "height", "created_at", "deposit_time", "activation_time", "status", "type", "title", "description", "deposit", "metadata", "changes", "yes", "no", "no_with_veto", "abstain", "yes_vals", "no_vals", "no_with_veto_vals", "abstain_vals", "yes_addrs", "no_addrs", "no_with_veto_addrs", "abstain_addrs", "votes_count", "voting_power", "yes_voting_power", "no_voting_power", "no_with_veto_voting_power", "abstain_voting_power").
			Model(&add).
			On("CONFLICT (id) DO UPDATE")

		if proposals[i].Deposit.IsPositive() {
			query.Set("deposit = added_proposal.deposit + EXCLUDED.deposit")
		}

		if !proposals[i].EmptyStatus() {
			query.Set("status = EXCLUDED.status")
		}

		if proposals[i].ActivationTime != nil {
			query.Set("activation_time = EXCLUDED.activation_time")
		}

		if proposals[i].VotesCount > 0 {
			query.Set("votes_count = added_proposal.votes_count + EXCLUDED.votes_count")
		}

		if proposals[i].VotingPower.IsPositive() {
			query.Set("voting_power = EXCLUDED.voting_power")
		}

		if proposals[i].Yes > 0 {
			query.Set("yes = added_proposal.yes + EXCLUDED.yes")
		}
		if proposals[i].No > 0 {
			query.Set("no = added_proposal.no + EXCLUDED.no")
		}
		if proposals[i].NoWithVeto > 0 {
			query.Set("no_with_veto = added_proposal.no_with_veto + EXCLUDED.no_with_veto")
		}
		if proposals[i].Abstain > 0 {
			query.Set("abstain = added_proposal.abstain + EXCLUDED.abstain")
		}

		if proposals[i].YesValidators > 0 {
			query.Set("yes_vals = added_proposal.yes_vals + EXCLUDED.yes_vals")
		}
		if proposals[i].NoValidators > 0 {
			query.Set("no_vals = added_proposal.no_vals + EXCLUDED.no_vals")
		}
		if proposals[i].NoWithVetoValidators > 0 {
			query.Set("no_with_veto_vals = added_proposal.no_with_veto_vals + EXCLUDED.no_with_veto_vals")
		}
		if proposals[i].AbstainValidators > 0 {
			query.Set("abstain_vals = added_proposal.abstain_vals + EXCLUDED.abstain_vals")
		}

		if proposals[i].YesAddress > 0 {
			query.Set("yes_addrs = added_proposal.yes_addrs + EXCLUDED.yes_addrs")
		}
		if proposals[i].NoAddress > 0 {
			query.Set("no_addrs = added_proposal.no_addrs+ EXCLUDED.no_addrs")
		}
		if proposals[i].NoWithVetoAddress > 0 {
			query.Set("no_with_veto_addrs = added_proposal.no_with_veto_addrs + EXCLUDED.no_with_veto_addrs")
		}
		if proposals[i].AbstainAddress > 0 {
			query.Set("abstain_addrs = added_proposal.abstain_addrs + EXCLUDED.abstain_addrs")
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

func (tx Transaction) SaveVotes(ctx context.Context, votes ...*models.Vote) error {
	if len(votes) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&votes).Exec(ctx)
	return err
}

func (tx Transaction) UpdateSlashedDelegations(ctx context.Context, validatorId uint64, fraction decimal.Decimal) (balances []models.Balance, err error) {
	if validatorId == 0 || !fraction.IsPositive() {
		return nil, nil
	}

	fr, _ := fraction.Float64()
	_, err = tx.Tx().NewUpdate().
		Model((*models.Delegation)(nil)).
		Set("amount = amount * (1 - ?)", fr).
		Where("validator_id = ?", validatorId).
		Returning("address_id as id, 'utia' as currency, -(amount / (1 - ?) - amount) as delegated", fr).
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

func (tx Transaction) Validators(ctx context.Context) (validators []models.Validator, err error) {
	err = tx.Tx().NewSelect().Model(&validators).
		Column("id", "stake").
		Scan(ctx)
	return
}

func (tx Transaction) ProposalVotes(ctx context.Context, proposalId uint64, limit, offset int) (votes []models.Vote, err error) {
	query := tx.Tx().NewSelect().Model(&votes).
		Where("proposal_id = ?", proposalId)

	query = limitScope(query, limit)
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
