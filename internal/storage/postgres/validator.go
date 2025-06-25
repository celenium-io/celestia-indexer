// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/shopspring/decimal"
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

func (v *Validator) TotalVotingPower(ctx context.Context) (decimal.Decimal, error) {
	q := v.DB().NewSelect().
		Model((*storage.Validator)(nil)).
		Column("stake").
		Where("jailed = false").
		Order("stake desc").
		Limit(100)

	var power decimal.Decimal
	err := v.DB().NewSelect().
		With("q", q).
		Table("q").
		ColumnExpr("sum(floor(stake / 1000000))").
		Scan(ctx, &power)
	return power, err
}

func (v *Validator) ListByPower(ctx context.Context, fltrs storage.ValidatorFilters) (validators []storage.Validator, err error) {
	query := v.DB().NewSelect().Model(&validators).
		OrderExpr("(not jailed)::int * stake desc")

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

	err = query.Scan(ctx)
	return
}

func (v *Validator) JailedCount(ctx context.Context) (int, error) {
	return v.DB().NewSelect().
		Model((*storage.Validator)(nil)).
		Where("jailed = true").
		Count(ctx)
}
