// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/math"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ValidatorMetrics struct {
	Id                    uint64          `bun:"id"`
	Moniker               string          `bun:"moniker"`
	MaxRate               decimal.Decimal `bun:"max_rate"`
	MaxChangeRate         decimal.Decimal `bun:"max_change_rate"`
	Stake                 decimal.Decimal `bun:"stake"`
	CreationTime          time.Time       `bun:"creation_time"`
	SelfDelegationAmount  decimal.Decimal `bun:"self_delegation_amount"`
	AppliedProposalsCount uint64          `bun:"applied_proposals_count"`
	VotesCount            uint64          `bun:"votes_count"`
	BlockMissedCount      uint64          `bun:"block_missed_count"`

	VotesMetric          decimal.Decimal `bun:"votes_metric"`
	CommissionMetric     decimal.Decimal `bun:"commission_metric"`
	OperationTimeMetric  decimal.Decimal `bun:"operation_time_metric"`
	SelfDelegationMetric decimal.Decimal `bun:"self_delegation_metric"`
	BlockMissedMetric    decimal.Decimal `bun:"block_missed_metric"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IValidator interface {
	storage.Table[*Validator]

	ByAddress(ctx context.Context, address string) (Validator, error)
	TotalVotingPower(ctx context.Context, maxVals int) (decimal.Decimal, error)
	ListByPower(ctx context.Context, fltrs ValidatorFilters) ([]Validator, error)
	JailedCount(ctx context.Context) (int, error)
	Messages(ctx context.Context, id uint64, fltrs ValidatorMessagesFilters) ([]MsgValidator, error)
	Metrics(ctx context.Context, id uint64) (ValidatorMetrics, error)
	TopNMetrics(ctx context.Context, n int) (ValidatorMetrics, error)
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
	Version     uint64          `bun:"version,default:0"        comment:"Signal version"                    json:"-"`

	Jailed *bool `bun:"jailed" comment:"True if validator was punished" json:"-"`

	MessagesCount uint64 `bun:"messages_count" comment:"Count of validator messages" json:"-"`

	CreationTime time.Time `bun:"creation_time" comment:"Creation time"`
}

func (Validator) TableName() string {
	return "validator"
}

func (v Validator) VotingPower() decimal.Decimal {
	return math.VotingPower(v.Stake)
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
	Limit   int
	Offset  int
	Jailed  *bool
	Version *int
}

type ValidatorMessagesFilters struct {
	Limit  int
	Offset int
	Sort   storage.SortOrder
	From   *time.Time
	To     *time.Time
}
