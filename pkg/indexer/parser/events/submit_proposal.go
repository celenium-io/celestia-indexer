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

func handleSubmitProposal(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if msg.Proposal == nil {
		return nil
	}
	event, _ := c.Peek()
	action := decoder.StringFromMap(event.Data, "action")
	isValid := action == "/cosmos.gov.v1beta1.MsgSubmitProposal" || action == "/cosmos.gov.v1.MsgSubmitProposal"
	if !isValid {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processSubmitProposal(ctx, c, msg)
}

func processSubmitProposal(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	event, ok := c.Peek()
	if !ok {
		return errors.New("not enough events for submit proposal")
	}
	if event.Type != types.EventTypeSubmitProposal {
		return errors.Errorf("submit proposal unexpected event type: %s", event.Type)
	}

	proposalId, err := decoder.Uint64FromMap(event.Data, "proposal_id")
	if err != nil {
		return errors.Errorf("submit proposal can't receive proposal id: %##v", event.Data)
	}
	msg.Proposal.Id = proposalId
	c.Skip(5)
	event, ok = c.Peek()
	if !ok {
		return errors.New("not enough events for submit proposal after parsing proposal id")
	}

	if event.Type != types.EventTypeProposalDeposit {
		return errors.Errorf("submit proposal unexpected event type: %s", event.Type)
	}
	deposit, err := decoder.NumericAmountFromMap(event.Data, "amount")
	if err != nil {
		return errors.Wrap(err, "parse proposal deposit amount")
	}
	msg.Proposal.Deposit = deposit
	if _, isV1 := event.Data["msg_index"]; isV1 {
		c.Skip(1)
	} else {
		c.Skip(2)
	}

	if event, ok := c.Peek(); ok {
		if event.Type == types.EventTypeSubmitProposal {
			msg.Proposal.Status = types.ProposalStatusActive
			msg.Proposal.ActivationTime = &ctx.Block.Time
		}
	}
	ctx.AddProposal(msg.Proposal)
	return nil
}
