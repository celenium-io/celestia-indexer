// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/json"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
)

type Proposal struct {
	Id              uint64         `example:"321"                       format:"int64"     json:"id"                        swaggertype:"integer"`
	Height          pkgTypes.Level `example:"100"                       format:"int64"     json:"height"                    swaggertype:"integer"`
	CreatedAt       time.Time      `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"created_at"                swaggertype:"string"`
	DepositTime     time.Time      `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"deposit_time"              swaggertype:"string"`
	ActivationtTime time.Time      `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"activation_time,omitempty" swaggertype:"string"`
	Status          string         `example:"active"                    format:"string"    json:"status"                    swaggertype:"string"`
	Type            string         `example:"param_changed"             format:"string"    json:"type"                      swaggertype:"string"`
	Title           string         `example:"Proposal title"            format:"string"    json:"title"                     swaggertype:"string"`
	Description     string         `example:"Some proposal description" format:"string"    json:"description"               swaggertype:"string"`
	Deposit         string         `example:"1000000000"                format:"string"    json:"deposit"                   swaggertype:"string"`
	Metadata        string         `example:"metadata"                  format:"string"    json:"metadata,omitempty"        swaggertype:"string"`

	VotesCount int64 `example:"12354" json:"votes_count"  swaggertype:"integer"`
	Yes        int64 `example:"1234"  json:"yes"          swaggertype:"integer"`
	No         int64 `example:"1234"  json:"no"           swaggertype:"integer"`
	NoWithVeto int64 `example:"1234"  json:"no_with_veto" swaggertype:"integer"`
	Abstain    int64 `example:"1234"  json:"abstain"      swaggertype:"integer"`

	YesVals        int64 `example:"1234" json:"yes_vals"          swaggertype:"integer"`
	NoVals         int64 `example:"1234" json:"no_vals"           swaggertype:"integer"`
	NoWithVetoVals int64 `example:"1234" json:"no_with_veto_vals" swaggertype:"integer"`
	AbstainVals    int64 `example:"1234" json:"abstain_vals"      swaggertype:"integer"`

	YesAddr        int64 `example:"1234" json:"yes_addrs"          swaggertype:"integer"`
	NoAddr         int64 `example:"1234" json:"no_addrs"           swaggertype:"integer"`
	NoWithVetoAddr int64 `example:"1234" json:"no_with_veto_addrs" swaggertype:"integer"`
	AbstainAddr    int64 `example:"1234" json:"abstain_addrs"      swaggertype:"integer"`

	VotingPower           string `example:"1000000000" format:"string" json:"voting_power"              swaggertype:"string"`
	YesVotingPower        string `example:"1000000000" format:"string" json:"yes_voting_power"          swaggertype:"string"`
	NoVotingPower         string `example:"1000000000" format:"string" json:"no_voting_power"           swaggertype:"string"`
	NoWithVetoVotingPower string `example:"1000000000" format:"string" json:"no_with_veto_voting_power" swaggertype:"string"`
	AbstainVotingPower    string `example:"1000000000" format:"string" json:"abstain_voting_power"      swaggertype:"string"`

	Changes  json.RawMessage `json:"changes,omitempty"`
	Proposer *ShortAddress   `json:"proposer,omitempty"`
}

func NewProposal(proposal storage.Proposal) Proposal {
	result := Proposal{
		Id:                    proposal.Id,
		Height:                proposal.Height,
		Proposer:              NewShortAddress(proposal.Proposer),
		CreatedAt:             proposal.CreatedAt,
		DepositTime:           proposal.DepositTime,
		Status:                proposal.Status.String(),
		Type:                  proposal.Type.String(),
		Title:                 proposal.Title,
		Description:           proposal.Description,
		Metadata:              proposal.Metadata,
		Deposit:               proposal.Deposit.String(),
		VotesCount:            proposal.VotesCount,
		Yes:                   proposal.Yes,
		No:                    proposal.No,
		NoWithVeto:            proposal.NoWithVeto,
		Abstain:               proposal.Abstain,
		YesVals:               proposal.YesValidators,
		NoVals:                proposal.NoValidators,
		NoWithVetoVals:        proposal.NoWithVetoValidators,
		AbstainVals:           proposal.AbstainValidators,
		YesAddr:               proposal.YesAddress,
		NoAddr:                proposal.NoAddress,
		NoWithVetoAddr:        proposal.NoWithVetoAddress,
		AbstainAddr:           proposal.AbstainAddress,
		VotingPower:           proposal.VotingPower.String(),
		YesVotingPower:        proposal.YesVotingPower.String(),
		NoVotingPower:         proposal.NoVotingPower.String(),
		NoWithVetoVotingPower: proposal.NoWithVetoVotingPower.String(),
		AbstainVotingPower:    proposal.AbstainVotingPower.String(),
		Changes:               proposal.Changes,
	}
	if proposal.ActivationTime != nil {
		result.ActivationtTime = *proposal.ActivationTime
	}

	return result
}
