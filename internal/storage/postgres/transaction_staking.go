// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	pg "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

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
		Column("id", "delegator", "address", "cons_address", "moniker", "website", "identity", "contacts", "details", "rate", "max_rate",
			"max_change_rate", "min_self_delegation", "stake", "jailed", "commissions", "rewards", "height", "version",
			"messages_count", "creation_time", "fibre_host", "fibre_host_height", "power", "bond_updates_count").
		On("CONFLICT ON CONSTRAINT address_validator DO UPDATE").
		Set("rate = CASE WHEN EXCLUDED.rate > 0 THEN EXCLUDED.rate ELSE added_validator.rate END").
		Set("max_rate = CASE WHEN EXCLUDED.max_rate > 0 THEN EXCLUDED.max_rate ELSE added_validator.max_rate END").
		Set("min_self_delegation = CASE WHEN EXCLUDED.min_self_delegation > 0 THEN EXCLUDED.min_self_delegation ELSE added_validator.min_self_delegation END").
		Set("stake = added_validator.stake + EXCLUDED.stake").
		Set("power = COALESCE(EXCLUDED.power, added_validator.power)").
		Set("commissions = added_validator.commissions + EXCLUDED.commissions").
		Set("rewards = added_validator.rewards + EXCLUDED.rewards").
		Set("messages_count = added_validator.messages_count + EXCLUDED.messages_count").
		Set("bond_updates_count = added_validator.bond_updates_count + EXCLUDED.bond_updates_count").
		Set("moniker = CASE WHEN EXCLUDED.moniker != '[do-not-modify]' THEN EXCLUDED.moniker ELSE added_validator.moniker END").
		Set("website = CASE WHEN EXCLUDED.website != '[do-not-modify]' THEN EXCLUDED.website ELSE added_validator.website END").
		Set("identity = CASE WHEN EXCLUDED.identity != '[do-not-modify]' THEN EXCLUDED.identity ELSE added_validator.identity END").
		Set("contacts = CASE WHEN EXCLUDED.contacts != '[do-not-modify]' THEN EXCLUDED.contacts ELSE added_validator.contacts END").
		Set("details = CASE WHEN EXCLUDED.details != '[do-not-modify]' THEN EXCLUDED.details ELSE added_validator.details END").
		Set("jailed = CASE WHEN EXCLUDED.jailed IS NOT NULL THEN EXCLUDED.jailed ELSE added_validator.jailed END").
		Set("version = CASE WHEN EXCLUDED.version > 0 THEN EXCLUDED.version ELSE added_validator.version END").
		Set("fibre_host = COALESCE(EXCLUDED.fibre_host, added_validator.fibre_host)").
		Set("fibre_host_height = COALESCE(EXCLUDED.fibre_host_height, added_validator.fibre_host_height)").
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

func (tx Transaction) SaveStakingLogs(ctx context.Context, logs ...models.StakingLog) error {
	return pg.SaveBulkWithCopy(ctx, tx, logs, copyThreshold)
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

func (tx Transaction) UpdateSlashedDelegations(ctx context.Context, validatorId uint64, burned storageTypes.Numeric) (balances []models.Balance, err error) {
	if validatorId == 0 || !burned.IsPositive() {
		return nil, nil
	}

	totalQuery := tx.Tx().NewSelect().
		Model((*models.Delegation)(nil)).
		ColumnExpr("sum(amount) as amount").
		Where("validator_id = ?", validatorId)

	burnedParts := tx.Tx().NewSelect().
		Table("total", "delegation").
		ColumnExpr("case when total.amount > 0 then (delegation.amount * ? / total.amount) else 0 end as amount", burned.String()).
		ColumnExpr("delegation.address_id as address_id").
		Where("validator_id = ?", validatorId)

	_, err = tx.Tx().NewUpdate().
		With("total", totalQuery).
		With("burned", burnedParts).
		Model((*models.Delegation)(nil)).
		TableExpr("burned").
		Set("amount = delegation.amount - burned.amount").
		Where("validator_id = ?", validatorId).
		Where("burned.address_id = delegation.address_id").
		Returning("delegation.address_id as id, 'utia' as currency, -burned.amount as delegated").
		Exec(ctx, &balances)
	return
}

func (tx Transaction) RollbackValidators(ctx context.Context, height types.Level) (validators []models.Validator, err error) {
	_, err = tx.Tx().NewDelete().Model(&validators).Where("height = ?", height).Returning("id").Exec(ctx)
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

func (tx Transaction) RollbackBondUpdates(ctx context.Context, height types.Level) (updates []models.ValidatorBondUpdate, err error) {
	_, err = tx.Tx().NewDelete().Model(&updates).
		Where("height = ?", height).
		Returning("*").
		Exec(ctx)
	return
}

func (tx Transaction) DeleteDelegationsByValidator(ctx context.Context, ids ...uint64) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := tx.Tx().NewDelete().
		Model((*models.Delegation)(nil)).
		Where("validator_id IN ?", bun.Tuple(ids)).
		Exec(ctx)
	return err
}

func (tx Transaction) CancelUnbondings(ctx context.Context, cancellations ...*models.Undelegation) error {
	if len(cancellations) == 0 {
		return nil
	}

	// one row per unbonding entry: cancellations of the same entry are summed by the decode
	// context, and the update below joins _data, so several matching rows would apply once
	data := make([]cancelUnbonding, len(cancellations))
	for i := range cancellations {
		data[i] = cancelUnbonding{
			Height:      cancellations[i].CreationHeight,
			ValidatorId: cancellations[i].ValidatorId,
			AddressId:   cancellations[i].AddressId,
			Amount:      cancellations[i].Amount,
		}
	}

	// the delete runs even though nothing selects from it: a data-modifying CTE
	// is always executed, and its snapshot does not see the update's own changes
	deleted := tx.Tx().NewDelete().
		Model((*models.Undelegation)(nil)).
		TableExpr("_data").
		Where("undelegation.height = _data.height").
		Where("undelegation.validator_id = _data.validator_id").
		Where("undelegation.address_id = _data.address_id").
		Where("undelegation.amount <= _data.amount")

	_, err := tx.Tx().NewUpdate().
		With("_data", tx.Tx().NewValues(&data)).
		With("_deleted", deleted).
		Model((*models.Undelegation)(nil)).
		TableExpr("_data").
		Set("amount = undelegation.amount - _data.amount").
		Where("undelegation.height = _data.height").
		Where("undelegation.validator_id = _data.validator_id").
		Where("undelegation.address_id = _data.address_id").
		Where("undelegation.amount > _data.amount").
		Exec(ctx)
	return err
}

type cancelUnbonding struct {
	bun.BaseModel `bun:"_data"`

	Height      types.Level          `bun:"height"`
	ValidatorId uint64               `bun:"validator_id"`
	AddressId   uint64               `bun:"address_id"`
	Amount      storageTypes.Numeric `bun:"amount,type:numeric"`
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
		Set("jailed = COALESCE(_data.jailed, validator.jailed)").
		Set("commissions = validator.commissions + _data.commissions").
		Set("rewards = validator.rewards + _data.rewards").
		Set("power = COALESCE(_data.power, validator.power)").
		Set("bond_updates_count = validator.bond_updates_count + _data.bond_updates_count").
		Where("validator.id = _data.id").
		Exec(ctx)
	return err
}
