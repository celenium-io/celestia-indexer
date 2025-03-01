// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type IVote interface {
	sdk.Table[*Vote]
}

type Vote struct {
	bun.BaseModel `bun:"vote" comment:"Table with proposal's votes"`

	Id         uint64           `bun:"id,pk,autoincrement"     comment:"Unique internal identity"`
	Height     pkgTypes.Level   `bun:"height,notnull"          comment:"The number (height) of this block"`
	Time       time.Time        `bun:"time,pk,notnull"         comment:"The time of block"`
	Option     types.VoteOption `bun:"option,type:vote_option" comment:"Selected vote option"`
	Weight     decimal.Decimal  `bun:"weight,type:numeric"     comment:"Vote's weight"`
	VoterId    uint64           `bun:"voter_id"                comment:"Voter internal identity"`
	ProposalId uint64           `bun:"proposal_id"             comment:"Proposal id"`

	Voter    *Address  `bun:"rel:belongs-to,join:voter_id=id"`
	Proposal *Proposal `bun:"rel:belongs-to,join:proposal_id=id"`
}

// TableName -
func (Vote) TableName() string {
	return "vote"
}
