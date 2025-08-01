// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBalance interface {
	storage.Table[*Balance]
}

type Balance struct {
	bun.BaseModel `bun:"balance" comment:"Table with account balances."`

	Id        uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Currency  string          `bun:"currency,pk,notnull"         comment:"Balance currency"`
	Spendable decimal.Decimal `bun:"spendable,type:numeric"      comment:"Spendable balance"`
	Delegated decimal.Decimal `bun:"delegated,type:numeric"      comment:"Delegated balance"`
	Unbonding decimal.Decimal `bun:"unbonding,type:numeric"      comment:"Unbonding balance"`
}

func (Balance) TableName() string {
	return "balance"
}

func (b Balance) IsEmpty() bool {
	return b.Currency == "" && b.Spendable.IsZero()
}

func EmptyBalance() Balance {
	return Balance{
		Currency:  currency.DefaultCurrency,
		Spendable: decimal.Zero,
		Delegated: decimal.Zero,
		Unbonding: decimal.Zero,
	}
}
