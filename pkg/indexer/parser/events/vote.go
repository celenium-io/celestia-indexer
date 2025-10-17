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

func handleVote(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	action := decoder.StringFromMap(events[*idx].Data, "action")
	isValid := action == "/cosmos.gov.v1beta1.MsgVote" || action == "/cosmos.gov.v1.MsgVote" || action == "/cosmos.gov.v1.MsgVoteWeighted" || action == "/cosmos.gov.v1beta1.MsgVoteWeighted"
	if !isValid {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processVote(ctx, events, msg, idx)
}

func processVote(ctx *context.Context, events []storage.Event, _ *storage.Message, idx *int) error {
	if events[*idx].Type != types.EventTypeProposalVote {
		return errors.Errorf("vote unexpected event type: %s", events[*idx].Type)
	}

	proposalId, err := decoder.Uint64FromMap(events[*idx].Data, "proposal_id")
	if err != nil {
		return errors.Errorf("vote can't receive proposal id: %##v", events[*idx].Data)
	}
	voter := decoder.StringFromMap(events[*idx].Data, "voter")
	option := decoder.StringFromMap(events[*idx].Data, "option")

	proposal := storage.Proposal{
		Id: proposalId,
	}

	if err := parseOption(ctx, proposalId, voter, option, &proposal, idx); err != nil {
		return errors.Wrap(err, "parse option")
	}

	ctx.AddProposal(&proposal)
	return nil
}

type optionType struct {
	Option int             `json:"option"`
	Weight decimal.Decimal `json:"weight"`
}

func parseOption(ctx *context.Context, proposalId uint64, voter, option string, proposal *storage.Proposal, idx *int) error {
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
					Balance:    storage.EmptyBalance(),
				},
			}
			if err := ctx.AddAddress(vote.Voter); err != nil {
				return err
			}

			switch opts[i].Option {
			case int(cosmosGovTypesV1.OptionAbstain):
				vote.Option = types.VoteOptionAbstain
				proposal.Abstain += 1
			case int(cosmosGovTypesV1.OptionNo):
				vote.Option = types.VoteOptionNo
				proposal.No += 1
			case int(cosmosGovTypesV1.OptionNoWithVeto):
				vote.Option = types.VoteOptionNoWithVeto
				proposal.NoWithVeto += 1
			case int(cosmosGovTypesV1.OptionYes):
				vote.Option = types.VoteOptionYes
				proposal.Yes += 1
			}
			vote.Weight = opts[i].Weight

			ctx.AddVote(&vote)
		}
		*idx += 1
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
			Balance:    storage.EmptyBalance(),
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
				proposal.Yes += 1
			case "VOTE_OPTION_NO":
				vote.Option = types.VoteOptionNo
				proposal.No += 1
			case "VOTE_OPTION_NO_WITH_VETO":
				vote.Option = types.VoteOptionNoWithVeto
				proposal.NoWithVeto += 1
			case "VOTE_OPTION_ABSTAIN":
				vote.Option = types.VoteOptionAbstain
				proposal.Abstain += 1
			}
		case "weight":
			value, err := strconv.Unquote(values[1])
			if err != nil {
				return errors.Errorf("unquote weight in vote option: %s", values[1])
			}
			vote.Weight = decimal.RequireFromString(value)
		}
	}

	ctx.AddVote(&vote)
	*idx += 2
	return nil
}
