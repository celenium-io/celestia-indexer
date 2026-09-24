// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Validator -
type Validator struct {
	*postgres.Table[*storage.Validator]
}

// NewValidator -
func NewValidator(db *database.Bun) *Validator {
	return &Validator{
		Table: postgres.NewTable[*storage.Validator](db),
	}
}

func (v *Validator) ByAddress(ctx context.Context, address string) (validator storage.Validator, err error) {
	err = v.DB().NewSelect().Model(&validator).
		Where("address = ?", address).
		Scan(ctx)
	return
}

func (v *Validator) TotalVotingPower(ctx context.Context) (storageTypes.Numeric, error) {
	var power storageTypes.Numeric
	err := v.DB().NewSelect().
		Model((*storage.Validator)(nil)).
		ColumnExpr("coalesce(sum(power), 0)").
		Scan(ctx, &power)
	return power, err
}

func (v *Validator) ListByPower(ctx context.Context, fltrs storage.ValidatorFilters) (validators []storage.Validator, err error) {
	query := v.DB().NewSelect().Model(&validators).
		OrderExpr("(jailed IS NOT TRUE) DESC, power DESC NULLS LAST, id DESC")

	query = limitScope(query, fltrs.Limit)
	if fltrs.Offset > 0 {
		query = query.Offset(fltrs.Offset)
	}
	if fltrs.Jailed != nil {
		if *fltrs.Jailed {
			query = query.Where("jailed = true")
		} else {
			query = query.Where("jailed = false")
		}
	}
	if fltrs.Version != nil {
		query = query.Where("version = ?", *fltrs.Version)
	}
	switch fltrs.Status {
	case storageTypes.ValidatorStatusActive:
		query = query.Where("power > 0 AND jailed IS NOT TRUE")
	case storageTypes.ValidatorStatusJailed:
		query = query.Where("jailed IS TRUE")
	case storageTypes.ValidatorStatusNotActive:
		query = query.Where("COALESCE(power, 0) = 0 AND jailed IS NOT TRUE")
	}

	err = query.Scan(ctx)
	return
}

func (v *Validator) Messages(ctx context.Context, id uint64, fltrs storage.ValidatorMessagesFilters) ([]storage.MsgValidator, error) {
	subQuery := v.DB().NewSelect().
		Model((*storage.MsgValidator)(nil)).
		Offset(fltrs.Offset).
		Where("validator_id = ?", id)

	subQuery = limitScope(subQuery, fltrs.Limit)
	subQuery = sortScope(subQuery, "time", fltrs.Sort)

	if fltrs.To != nil {
		subQuery = subQuery.Where("time < ?", *fltrs.To)
	}
	if fltrs.From != nil {
		subQuery = subQuery.Where("time >= ?", *fltrs.From)
	}

	var response []storage.MsgValidator
	err := v.DB().NewSelect().
		TableExpr("(?) as msgs", subQuery).
		ColumnExpr("msgs.*").
		ColumnExpr("message.position as msg__position, message.type as msg__type, message.size as msg__size, message.data as msg__data").
		Join("left join message on message.id = msg_id").
		Scan(ctx, &response)

	return response, err
}

func (v *Validator) Metrics(ctx context.Context, id uint64) (metrics storage.ValidatorMetrics, err error) {
	err = v.DB().NewSelect().
		Table(storage.ViewValidatorMetrics).
		Where("id = ?", id).
		Scan(ctx, &metrics)
	return
}

func (v *Validator) TopNMetrics(ctx context.Context, n int) (metrics storage.ValidatorMetrics, err error) {
	subQuery := v.DB().NewSelect().
		Table(storage.ViewValidatorMetrics).
		Order("stake desc").
		Limit(n)

	err = v.DB().NewSelect().
		With("metrics", subQuery).
		Table("metrics").
		ColumnExpr("avg(commission_metric) as commission_metric").
		ColumnExpr("avg(votes_metric) as votes_metric").
		ColumnExpr("avg(operation_time_metric) as operation_time_metric").
		ColumnExpr("avg(self_delegation_metric) as self_delegation_metric").
		ColumnExpr("avg(block_missed_metric) as block_missed_metric").
		Scan(ctx, &metrics)
	return
}

func (v *Validator) CountByStatus(ctx context.Context) (response storage.CountByStatus, err error) {
	err = v.DB().NewSelect().
		Model((*storage.Validator)(nil)).
		ColumnExpr("count(*) AS total").
		ColumnExpr("count(*) FILTER (WHERE jailed IS TRUE) AS jailed").
		ColumnExpr("count(*) FILTER (WHERE jailed IS NOT TRUE AND COALESCE(power, 0) = 0) AS not_active").
		ColumnExpr("count(*) FILTER (WHERE jailed IS NOT TRUE AND power > 0) AS active").
		Scan(ctx, &response)
	return
}
