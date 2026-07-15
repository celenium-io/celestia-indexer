// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleDeposit(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/cosmos.gov.v1.MsgDeposit" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processDeposit(ctx, c, msg)
}

func processDeposit(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	c.Skip(4)
	event, ok := c.Peek()
	if !ok {
		return errors.New("proposal deposit unexpected events length")
	}
	if event.Type != types.EventTypeProposalDeposit {
		return errors.Errorf("proposal deposit unexpected event type: %s", event.Type)
	}

	proposalId, err := decoder.Uint64FromMap(event.Data, "proposal_id")
	if err != nil {
		return errors.Errorf("submit proposal can't receive proposal id: %##v", event.Data)
	}
	depositAmount, err := decoder.NumericAmountFromMap(event.Data, "amount")
	if err != nil {
		return errors.Wrap(err, "parse deposit amount")
	}
	msg.Proposal = &storage.Proposal{
		Id:      proposalId,
		Deposit: depositAmount,
		Status:  types.ProposalStatusInactive,
	}

	c.Next()
	for {
		event, ok := c.Peek()
		if !ok {
			break
		}
		if event.Type == types.EventTypeProposalDeposit {
			votingPeriodStart, err := decoder.Uint64FromMap(event.Data, "voting_period_start")
			if err != nil {
				return errors.Errorf("submit proposal can't receive voting_period_start: %##v", event.Data)
			}
			if votingPeriodStart == proposalId {
				msg.Proposal.Status = types.ProposalStatusActive
				msg.Proposal.ActivationTime = &ctx.Block.Time
			}
			break
		}

		if action := decoder.StringFromMap(event.Data, "action"); action != "" {
			break
		}
		c.Next()
	}
	ctx.AddProposal(msg.Proposal)

	return nil
}
