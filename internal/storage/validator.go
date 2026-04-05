// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/math"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ValidatorMetrics struct {
	Id                    uint64        `bun:"id"`
	Moniker               string        `bun:"moniker"`
	MaxRate               types.Numeric `bun:"max_rate"`
	MaxChangeRate         types.Numeric `bun:"max_change_rate"`
	Stake                 types.Numeric `bun:"stake"`
	CreationTime          time.Time     `bun:"creation_time"`
	SelfDelegationAmount  types.Numeric `bun:"self_delegation_amount"`
	AppliedProposalsCount uint64        `bun:"applied_proposals_count"`
	VotesCount            uint64        `bun:"votes_count"`
	BlockMissedCount      uint64        `bun:"block_missed_count"`

	VotesMetric          types.Numeric `bun:"votes_metric"`
	CommissionMetric     types.Numeric `bun:"commission_metric"`
	OperationTimeMetric  types.Numeric `bun:"operation_time_metric"`
	SelfDelegationMetric types.Numeric `bun:"self_delegation_metric"`
	BlockMissedMetric    types.Numeric `bun:"block_missed_metric"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IValidator interface {
	storage.Table[*Validator]

	ByAddress(ctx context.Context, address string) (Validator, error)
	TotalVotingPower(ctx context.Context, maxVals int) (types.Numeric, error)
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

	Rate              types.Numeric `bun:"rate,type:numeric"                comment:"Commission rate charged to delegators, as a fraction"                   json:"-"`
	MaxRate           types.Numeric `bun:"max_rate,type:numeric"            comment:"Maximum commission rate which validator can ever charge, as a fraction" json:"-"`
	MaxChangeRate     types.Numeric `bun:"max_change_rate,type:numeric"     comment:"Maximum daily increase of the validator commission, as a fraction"      json:"-"`
	MinSelfDelegation types.Numeric `bun:"min_self_delegation,type:numeric" comment:""                                                                       json:"-"`

	Stake       types.Numeric  `bun:"stake,type:numeric"       comment:"Validator's stake"                 json:"-"`
	Rewards     types.Numeric  `bun:"rewards,type:numeric"     comment:"Validator's rewards"               json:"-"`
	Commissions types.Numeric  `bun:"commissions,type:numeric" comment:"Commissions"                       json:"-"`
	Height      pkgTypes.Level `bun:"height"                   comment:"Height when validator was created" json:"-"`
	Version     uint64         `bun:"version,default:0"        comment:"Signal version"                    json:"-"`

	Jailed *bool `bun:"jailed" comment:"True if validator was punished" json:"-"`

	MessagesCount uint64 `bun:"messages_count" comment:"Count of validator messages" json:"-"`

	CreationTime time.Time `bun:"creation_time" comment:"Creation time"`
}

func (Validator) TableName() string {
	return "validator"
}

func (v Validator) VotingPower() types.Numeric {
	return math.SharesNumeric(v.Stake)
}

const DoNotModify = "[do-not-modify]"

func EmptyValidator() Validator {
	return Validator{
		Rate:              types.NewNumeric(decimal.Zero),
		MaxRate:           types.NewNumeric(decimal.Zero),
		MaxChangeRate:     types.NewNumeric(decimal.Zero),
		MinSelfDelegation: types.NewNumeric(decimal.Zero),
		Rewards:           types.NewNumeric(decimal.Zero),
		Commissions:       types.NewNumeric(decimal.Zero),
		Stake:             types.NewNumeric(decimal.Zero),
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
