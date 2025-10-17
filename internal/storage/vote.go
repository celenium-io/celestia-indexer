// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
type IVote interface {
	sdk.Table[*Vote]

	ByProposalId(ctx context.Context, proposalId uint64, fltrs VoteFilters) ([]Vote, error)
	ByVoterId(ctx context.Context, voterId uint64, fltrs VoteFilters) ([]Vote, error)
	ByValidatorId(ctx context.Context, validatorId uint64, fltrs VoteFilters) ([]Vote, error)
}

type Vote struct {
	bun.BaseModel `bun:"vote" comment:"Table with proposal's votes"`

	Id          uint64           `bun:"id,pk,autoincrement"     comment:"Unique internal identity"`
	Height      pkgTypes.Level   `bun:"height,notnull"          comment:"The number (height) of this block"`
	Time        time.Time        `bun:"time,pk,notnull"         comment:"The time of block"`
	Option      types.VoteOption `bun:"option,type:vote_option" comment:"Selected vote option"`
	Weight      decimal.Decimal  `bun:"weight,type:numeric"     comment:"Vote's weight"`
	VoterId     uint64           `bun:"voter_id"                comment:"Voter internal identity"`
	ProposalId  uint64           `bun:"proposal_id"             comment:"Proposal id"`
	ValidatorId *uint64          `bun:"validator_id"            comment:"Validator id"`

	Voter     *Address   `bun:"rel:belongs-to,join:voter_id=id"`
	Proposal  *Proposal  `bun:"rel:belongs-to,join:proposal_id=id"`
	Validator *Validator `bun:"rel:belongs-to,join:validator_id=id"`
}

// TableName -
func (Vote) TableName() string {
	return "vote"
}

type VoteFilters struct {
	Limit       int
	Offset      int
	Option      []types.VoteOption
	VoterType   types.VoterType
	AddressId   *uint64
	ValidatorId *uint64
}
type VotesCount struct {
	VotesCount int64
	Yes        int64
	No         int64
	NoWithVeto int64
	Abstain    int64

	YesValidators        int64
	NoValidators         int64
	NoWithVetoValidators int64
	AbstainValidators    int64

	YesAddress        int64
	NoAddress         int64
	NoWithVetoAddress int64
	AbstainAddress    int64
}

func (vc *VotesCount) Update(count int64, vote Vote) {
	vc.VotesCount += count

	switch vote.Option {
	case types.VoteOptionAbstain:
		vc.Abstain += count
		if vote.ValidatorId != nil {
			vc.AbstainValidators += count
		} else {
			vc.AbstainAddress += count
		}

	case types.VoteOptionNo:
		vc.No += count
		if vote.ValidatorId != nil {
			vc.NoValidators += count
		} else {
			vc.NoAddress += count
		}

	case types.VoteOptionNoWithVeto:
		vc.NoWithVeto += count
		if vote.ValidatorId != nil {
			vc.NoWithVetoValidators += count
		} else {
			vc.NoWithVetoAddress += count
		}

	case types.VoteOptionYes:
		vc.Yes += count
		if vote.ValidatorId != nil {
			vc.YesValidators += count
		} else {
			vc.YesAddress += count
		}
	}
}
