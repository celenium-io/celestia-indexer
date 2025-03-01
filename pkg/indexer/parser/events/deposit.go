// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleDeposit(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/cosmos.gov.v1.MsgDeposit" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processDeposit(ctx, events, msg, idx)
}

func processDeposit(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	*idx += 4
	if events[*idx].Type != types.EventTypeProposalDeposit {
		return errors.Errorf("proposal deposit unexpected event type: %s", events[*idx].Type)
	}

	proposalId, err := decoder.Uint64FromMap(events[*idx].Data, "proposal_id")
	if err != nil {
		return errors.Errorf("submit proposal can't receive proposal id: %##v", events[*idx].Data)
	}
	amount := decoder.AmountFromMap(events[*idx].Data, "amount")
	msg.Proposal = &storage.Proposal{
		Id:      proposalId,
		Deposit: amount,
		Status:  types.ProposalStatusInactive,
	}

	*idx += 2
	if len(events) > *idx {
		if events[*idx].Type == types.EventTypeProposalDeposit {
			msg.Proposal.Status = types.ProposalStatusActive
			msg.Proposal.ActivationTime = &ctx.Block.Time
		}
		*idx += 1
	}
	ctx.AddProposal(msg.Proposal)

	return nil
}
