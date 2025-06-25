// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test_handleVote(t *testing.T) {
	ts := time.Now()
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
		idx    *int
		vote   storage.Vote
	}{
		{
			name: "vote test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 3762606,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.gov.v1beta1.MsgVote",
					},
				}, {
					Height: 3762606,
					Type:   "proposal_vote",
					Data: map[string]any{
						"option":      "option:VOTE_OPTION_YES weight:\"1.000000000000000000\"",
						"proposal_id": "5",
						"voter":       "celestia1lha3l9w5sca8gv98t27n5rezex3vafp8qa7h69",
					},
				}, {
					Height: 3762606,
					Type:   "message",
					Data: map[string]any{
						"module": "governance",
						"sender": "celestia1lha3l9w5sca8gv98t27n5rezex3vafp8qa7h69",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgVote,
				Height: 3762606,
			},
			idx: testsuite.Ptr(0),
			vote: storage.Vote{
				ProposalId: 5,
				Voter: &storage.Address{
					Address: "celestia1lha3l9w5sca8gv98t27n5rezex3vafp8qa7h69",
				},
				Option: types.VoteOptionYes,
				Weight: decimal.RequireFromString("1.000000000000000000"),
				Height: 3762606,
				Time:   ts,
			},
		}, {
			name: "vote test 2",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 871324,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.gov.v1.MsgVoteWeighted",
					},
				}, {
					Height: 871324,
					Type:   "proposal_vote",
					Data: map[string]any{
						"option":      "option:VOTE_OPTION_YES weight:\"1.000000000000000000\"",
						"proposal_id": "3",
						"voter":       "celestia1davz40kat93t49ljrkmkl5uqhqq45e0tedgf8a",
					},
				}, {
					Height: 871324,
					Type:   "message",
					Data: map[string]any{
						"module": "governance",
						"sender": "celestia1davz40kat93t49ljrkmkl5uqhqq45e0tedgf8a",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgVoteWeighted,
				Height: 871324,
			},
			idx: testsuite.Ptr(0),
			vote: storage.Vote{
				ProposalId: 3,
				Voter: &storage.Address{
					Address: "celestia1davz40kat93t49ljrkmkl5uqhqq45e0tedgf8a",
				},
				Option: types.VoteOptionYes,
				Weight: decimal.RequireFromString("1.000000000000000000"),
				Height: 871324,
				Time:   ts,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ctx.Block = &storage.Block{
				Time:   ts,
				Height: tt.msg.Height,
			}
			err := handleVote(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
			require.Len(t, tt.ctx.Votes, 1)
			require.Equal(t, tt.vote, *tt.ctx.Votes[0])
		})
	}
}
