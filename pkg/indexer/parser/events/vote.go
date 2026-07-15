// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func handleVote(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	action := decoder.StringFromMap(event.Data, "action")
	isValid := action == "/cosmos.gov.v1beta1.MsgVote" || action == "/cosmos.gov.v1.MsgVote" || action == "/cosmos.gov.v1.MsgVoteWeighted" || action == "/cosmos.gov.v1beta1.MsgVoteWeighted"
	if !isValid {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processVote(ctx, c, msg)
}

func processVote(ctx *context.Context, c *Cursor, _ *storage.Message) error {
	event, ok := c.Peek()
	if !ok {
		return errors.New("not enough events for vote")
	}
	if event.Type != types.EventTypeProposalVote {
		return errors.Errorf("vote unexpected event type: %s", event.Type)
	}

	proposalId, err := decoder.Uint64FromMap(event.Data, "proposal_id")
	if err != nil {
		return errors.Errorf("vote can't receive proposal id: %##v", event.Data)
	}
	voter := decoder.StringFromMap(event.Data, "voter")
	option := decoder.StringFromMap(event.Data, "option")

	if err := parseOption(ctx, proposalId, voter, option, c); err != nil {
		return errors.Wrap(err, "parse option")
	}

	ctx.AddProposal(&storage.Proposal{
		Id: proposalId,
	})
	return nil
}

type optionType struct {
	Option int             `json:"option"`
	Weight decimal.Decimal `json:"weight"`
}

func parseOption(ctx *context.Context, proposalId uint64, voter, option string, c *Cursor) error {
	var opts []optionType
	if err := json.Unmarshal([]byte(option), &opts); err == nil {
		if len(opts) == 0 {
			return errors.New("empty vote options array")
		}

		for i := range opts {
			vote := storage.Vote{
				ProposalId: proposalId,
				Time:       ctx.Block.Time,
				Height:     ctx.Block.Height,
				Voter: &storage.Address{
					Height:     ctx.Block.Height,
					LastHeight: ctx.Block.Height,
					Address:    voter,
					Balances:   []storage.Balance{storage.EmptyBalance()},
				},
			}
			if err := ctx.AddAddress(vote.Voter); err != nil {
				return err
			}

			switch opts[i].Option {
			case int(cosmosGovTypesV1.OptionAbstain):
				vote.Option = types.VoteOptionAbstain
			case int(cosmosGovTypesV1.OptionNo):
				vote.Option = types.VoteOptionNo
			case int(cosmosGovTypesV1.OptionNoWithVeto):
				vote.Option = types.VoteOptionNoWithVeto
			case int(cosmosGovTypesV1.OptionYes):
				vote.Option = types.VoteOptionYes
			}
			vote.Weight = types.NewNumeric(opts[i].Weight)

			ctx.AddVote(&vote)
		}
		c.Skip(1)
		return nil
	}

	vote := storage.Vote{
		ProposalId: proposalId,
		Time:       ctx.Block.Time,
		Height:     ctx.Block.Height,
		Voter: &storage.Address{
			Height:     ctx.Block.Height,
			LastHeight: ctx.Block.Height,
			Address:    voter,
			Balances:   []storage.Balance{storage.EmptyBalance()},
		},
	}

	if err := ctx.AddAddress(vote.Voter); err != nil {
		return err
	}

	optionParts := strings.Split(option, " ")
	for i := range optionParts {
		values := strings.Split(optionParts[i], ":")
		if len(values) != 2 {
			continue
		}
		switch values[0] {
		case "option":
			switch values[1] {
			case "VOTE_OPTION_YES":
				vote.Option = types.VoteOptionYes
			case "VOTE_OPTION_NO":
				vote.Option = types.VoteOptionNo
			case "VOTE_OPTION_NO_WITH_VETO":
				vote.Option = types.VoteOptionNoWithVeto
			case "VOTE_OPTION_ABSTAIN":
				vote.Option = types.VoteOptionAbstain
			}
		case "weight":
			value, err := strconv.Unquote(values[1])
			if err != nil {
				return errors.Errorf("unquote weight in vote option: %s", values[1])
			}
			w, err := types.NumericFromString(value)
			if err != nil {
				return errors.Wrap(err, "parse vote weight")
			}
			vote.Weight = w
		}
	}

	ctx.AddVote(&vote)
	c.Skip(2)
	return nil
}
