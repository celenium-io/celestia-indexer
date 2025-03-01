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
	sdk.Table[*Proposal]

	ByProposer(ctx context.Context, id uint64, limit, offset int) ([]Proposal, error)
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

	Yes        int64 `bun:"yes"          comment:"Count of yes votes"`
	No         int64 `bun:"no"           comment:"Count of no votes"`
	NoWithVeto int64 `bun:"no_with_veto" comment:"Count of no votes with veto"`
	Abstain    int64 `bun:"abstain"      comment:"Count of abstain votes"`

	Proposer *Address `bun:"rel:belongs-to,join:proposer_id=id"`
}

// TableName -
func (Proposal) TableName() string {
	return "proposal"
}

func (p Proposal) EmptyStatus() bool {
	return p.Status == "" || p.Status == types.ProposalStatusInactive
}
