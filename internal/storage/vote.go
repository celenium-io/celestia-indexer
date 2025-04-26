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
	ValidatorId uint64           `bun:"validator_id"            comment:"Validator id"`

	Voter     *Address   `bun:"rel:belongs-to,join:voter_id=id"`
	Proposal  *Proposal  `bun:"rel:belongs-to,join:proposal_id=id"`
	Validator *Validator `bun:"rel:belongs-to,join:validator_id=id"`
}

// TableName -
func (Vote) TableName() string {
	return "vote"
}

type VoteFilters struct {
	Limit     int
	Offset    int
	Option    types.VoteOption
	VoterType types.VoterType
}
