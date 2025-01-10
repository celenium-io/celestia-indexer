// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IValidator interface {
	storage.Table[*Validator]

	ByAddress(ctx context.Context, address string) (Validator, error)
	TotalVotingPower(ctx context.Context) (decimal.Decimal, error)
	ListByPower(ctx context.Context, fltrs ValidatorFilters) ([]Validator, error)
	JailedCount(ctx context.Context) (int, error)
}

type Validator struct {
	bun.BaseModel `bun:"validator" comment:"Table with celestia validators."`

	Id          uint64 `bun:"id,pk,notnull,autoincrement"                comment:"Unique internal identity"`
	Delegator   string `bun:"delegator,type:text"                        comment:"Delegator address"        json:"-"`
	Address     string `bun:"address,unique:address_validator,type:text" comment:"Validator address"        json:"-"`
	ConsAddress string `bun:"cons_address"                               comment:"Consensus address"        json:"-"`

	Moniker  string `bun:"moniker,type:text"  comment:"Human-readable name for the validator" json:"-"`
	Website  string `bun:"website,type:text"  comment:"Website link"                          json:"-"`
	Identity string `bun:"identity,type:text" comment:"Optional identity signature"           json:"-"`
	Contacts string `bun:"contacts,type:text" comment:"Contacts"                              json:"-"`
	Details  string `bun:"details,type:text"  comment:"Detailed information about validator"  json:"-"`

	Rate              decimal.Decimal `bun:"rate,type:numeric"                comment:"Commission rate charged to delegators, as a fraction"                   json:"-"`
	MaxRate           decimal.Decimal `bun:"max_rate,type:numeric"            comment:"Maximum commission rate which validator can ever charge, as a fraction" json:"-"`
	MaxChangeRate     decimal.Decimal `bun:"max_change_rate,type:numeric"     comment:"Maximum daily increase of the validator commission, as a fraction"      json:"-"`
	MinSelfDelegation decimal.Decimal `bun:"min_self_delegation,type:numeric" comment:""                                                                       json:"-"`

	Stake       decimal.Decimal `bun:"stake,type:numeric"       comment:"Validator's stake"                 json:"-"`
	Rewards     decimal.Decimal `bun:"rewards,type:numeric"     comment:"Validator's rewards"               json:"-"`
	Commissions decimal.Decimal `bun:"commissions,type:numeric" comment:"Commissions"                       json:"-"`
	Height      pkgTypes.Level  `bun:"height"                   comment:"Height when validator was created" json:"-"`

	Jailed *bool `bun:"jailed" comment:"True if validator was punished" json:"-"`
}

func (Validator) TableName() string {
	return "validator"
}

const DoNotModify = "[do-not-modify]"

func EmptyValidator() Validator {
	return Validator{
		Rate:              decimal.Zero,
		MaxRate:           decimal.Zero,
		MaxChangeRate:     decimal.Zero,
		MinSelfDelegation: decimal.Zero,
		Rewards:           decimal.Zero,
		Commissions:       decimal.Zero,
		Stake:             decimal.Zero,
		Contacts:          DoNotModify,
		Details:           DoNotModify,
		Identity:          DoNotModify,
		Moniker:           DoNotModify,
		Website:           DoNotModify,
	}
}

type ValidatorFilters struct {
	Limit  int
	Offset int
	Jailed *bool
}
