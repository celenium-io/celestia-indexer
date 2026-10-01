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

		votes := make([]*storage.Vote, len(opts))
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
			votes[i] = &vote
		}
		ctx.AddVotes(votes...)
		c.Skip(1)
		return nil
	}

	var votes []*storage.Vote
	for _, field := range strings.Fields(option) {
		key, value, ok := strings.Cut(field, ":")
		if !ok {
			continue
		}
		switch key {
		case "option":
			voterAddress := &storage.Address{
				Height:     ctx.Block.Height,
				LastHeight: ctx.Block.Height,
				Address:    voter,
				Balances:   []storage.Balance{storage.EmptyBalance()},
			}
			votes = append(votes, &storage.Vote{
				ProposalId: proposalId,
				Time:       ctx.Block.Time,
				Height:     ctx.Block.Height,
				Voter:      voterAddress,
			})

			if err := ctx.AddAddress(voterAddress); err != nil {
				return err
			}

			switch value {
			case "VOTE_OPTION_YES":
				votes[len(votes)-1].Option = types.VoteOptionYes
			case "VOTE_OPTION_NO":
				votes[len(votes)-1].Option = types.VoteOptionNo
			case "VOTE_OPTION_NO_WITH_VETO":
				votes[len(votes)-1].Option = types.VoteOptionNoWithVeto
			case "VOTE_OPTION_ABSTAIN":
				votes[len(votes)-1].Option = types.VoteOptionAbstain
			}
		case "weight":
			if len(votes) == 0 {
				return errors.Errorf("weight before option: %s", option)
			}
			unquoted, err := strconv.Unquote(value)
			if err != nil {
				return errors.Wrap(err, "unquote weight")
			}
			weight, err := types.NumericFromString(unquoted)
			if err != nil {
				return errors.Wrap(err, "parse weight")
			}
			votes[len(votes)-1].Weight = weight
		}
	}
	ctx.AddVotes(votes...)
	c.Skip(2)
	return nil
}
