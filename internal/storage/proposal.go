// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IProposal interface {
	ListWithFilters(ctx context.Context, filters ListProposalFilters) (proposals []Proposal, err error)
	ById(ctx context.Context, id uint64) (Proposal, error)
}

type Proposal struct {
	bun.BaseModel `bun:"proposal" comment:"Table with proposals"`

	Id             uint64               `bun:"id,pk"                       comment:"Unique identity"`
	Height         pkgTypes.Level       `bun:"height"                      comment:"The number (height) of this block"`
	ProposerId     uint64               `bun:"proposer_id"                 comment:"Proposer address identity"`
	CreatedAt      time.Time            `bun:"created_at,notnull"          comment:"Creation time"`
	DepositTime    time.Time            `bun:"deposit_time"                comment:"Time to end of deposit period"`
	ActivationTime *time.Time           `bun:"activation_time"             comment:"Voting start time"`
	Status         types.ProposalStatus `bun:"status,type:proposal_status" comment:"Proposal status"`
	Type           types.ProposalType   `bun:"type,type:proposal_type"     comment:"Proposal type"`
	Title          string               `bun:"title"                       comment:"Title"`
	Description    string               `bun:"description"                 comment:"Proposal description"`
	Deposit        decimal.Decimal      `bun:"deposit,type:numeric"        comment:"Deposit"`
	Metadata       string               `bun:"metadata"                    comment:"Metadata"`
	Changes        []byte               `bun:"changes,type:bytea"          comment:"JSON object with proposal changes"`

	VotesCount int64 `bun:"votes_count"  comment:"Total votes count"`
	Yes        int64 `bun:"yes"          comment:"Count of yes votes"`
	No         int64 `bun:"no"           comment:"Count of no votes"`
	NoWithVeto int64 `bun:"no_with_veto" comment:"Count of no votes with veto"`
	Abstain    int64 `bun:"abstain"      comment:"Count of abstain votes"`

	YesValidators        int64 `bun:"yes_vals"          comment:"Count of yes votes by validators"`
	NoValidators         int64 `bun:"no_vals"           comment:"Count of no votes by validators"`
	NoWithVetoValidators int64 `bun:"no_with_veto_vals" comment:"Count of no votes with veto by validators"`
	AbstainValidators    int64 `bun:"abstain_vals"      comment:"Count of abstain votes by validators"`

	YesAddress        int64 `bun:"yes_addrs"          comment:"Count of yes votes by addresses"`
	NoAddress         int64 `bun:"no_addrs"           comment:"Count of no votes by addresses"`
	NoWithVetoAddress int64 `bun:"no_with_veto_addrs" comment:"Count of no votes with veto by addresses"`
	AbstainAddress    int64 `bun:"abstain_addrs"      comment:"Count of abstain votes by addresses"`

	VotingPower           decimal.Decimal `bun:"voting_power,type:numeric"              comment:"Total voting power"`
	YesVotingPower        decimal.Decimal `bun:"yes_voting_power,type:numeric"          comment:"Yes voting power"`
	NoVotingPower         decimal.Decimal `bun:"no_voting_power,type:numeric"           comment:"No voting power"`
	NoWithVetoVotingPower decimal.Decimal `bun:"no_with_veto_voting_power,type:numeric" comment:"No with veto voting power"`
	AbstainVotingPower    decimal.Decimal `bun:"abstain_voting_power,type:numeric"      comment:"Abstain voting power"`

	Proposer *Address `bun:"rel:belongs-to,join:proposer_id=id"`
}

// TableName -
func (Proposal) TableName() string {
	return "proposal"
}

func (p Proposal) EmptyStatus() bool {
	return p.Status == "" || p.Status == types.ProposalStatusInactive
}

func (p Proposal) Finished() bool {
	return p.Status == types.ProposalStatusApplied || p.Status == types.ProposalStatusRejected
}

type ListProposalFilters struct {
	Limit      int
	Offset     int
	ProposerId uint64
	Sort       sdk.SortOrder
	Status     []types.ProposalStatus
	Type       []types.ProposalType
}
