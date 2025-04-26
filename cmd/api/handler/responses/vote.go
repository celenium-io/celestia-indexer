// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/shopspring/decimal"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
)

type Vote struct {
	Id         uint64           `example:"321"                       format:"int64"     json:"id"           swaggertype:"integer"`
	Height     pkgTypes.Level   `example:"100"                       format:"int64"     json:"height"       swaggertype:"integer"`
	Time       time.Time        `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"deposit_time" swaggertype:"string"`
	Option     types.VoteOption `example:"yes"                       format:"string"    json:"status"       swaggertype:"string"`
	Weight     decimal.Decimal  `example:"12345678"                  format:"int64"     json:"weight"       swaggertype:"integer"`
	VoterId    uint64           `example:"1"                         format:"int64"     json:"voter_id"     swaggertype:"integer"`
	ProposalId uint64           `example:"2"                         format:"int64"     json:"proposal_id"  swaggertype:"integer"`

	Voter     *ShortAddress   `json:"proposer,omitempty"`
	Validator *ShortValidator `json:"validator,omitempty"`
	Proposal  Proposal        `json:"-"`
}

func NewVote(vote storage.Vote) Vote {
	result := Vote{
		Id:         vote.Id,
		Height:     vote.Height,
		Time:       vote.Time,
		Option:     vote.Option,
		Weight:     vote.Weight,
		VoterId:    vote.VoterId,
		ProposalId: vote.ProposalId,
	}

	if vote.Voter != nil {
		result.Voter = NewShortAddress(vote.Voter)
	}

	if vote.Validator != nil {
		result.Validator = NewShortValidator(*vote.Validator)
	}

	return result
}
