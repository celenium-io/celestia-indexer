// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
		votes  []*storage.Vote
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
			votes: []*storage.Vote{
				{
					ProposalId: 5,
					Voter: &storage.Address{
						Address:    "celestia1lha3l9w5sca8gv98t27n5rezex3vafp8qa7h69",
						Height:     3762606,
						LastHeight: 3762606,
						Balance:    storage.EmptyBalance(),
						Hash:       []byte{0xfd, 0xfb, 0x1f, 0x95, 0xd4, 0x86, 0x3a, 0x74, 0x30, 0xa7, 0x5a, 0xbd, 0x3a, 0x0f, 0x22, 0xc9, 0xa2, 0xce, 0xa4, 0x27},
					},
					Option: types.VoteOptionYes,
					Weight: decimal.RequireFromString("1.000000000000000000"),
					Height: 3762606,
					Time:   ts,
				},
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
			votes: []*storage.Vote{
				{
					ProposalId: 3,
					Voter: &storage.Address{
						Address:    "celestia1davz40kat93t49ljrkmkl5uqhqq45e0tedgf8a",
						Height:     871324,
						LastHeight: 871324,
						Balance:    storage.EmptyBalance(),
						Hash:       []byte{0x6f, 0x58, 0x2a, 0xbe, 0xdd, 0x59, 0x62, 0xba, 0x97, 0xf2, 0x1d, 0xb7, 0x6f, 0xd3, 0x80, 0xb8, 0x01, 0x5a, 0x65, 0xeb},
					},
					Option: types.VoteOptionYes,
					Weight: decimal.RequireFromString("1.000000000000000000"),
					Height: 871324,
					Time:   ts,
				},
			},
		}, {
			name: "vote test 3",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 871324,
					Type:   "message",
					Data: map[string]any{
						"action":    "/cosmos.gov.v1.MsgVote",
						"sender":    "celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r",
						"module":    "gov",
						"msg_index": "0",
					},
				}, {
					Height: 871324,
					Type:   "proposal_vote",
					Data: map[string]any{
						"option":      "[{\"option\":1,\"weight\":\"1.000000000000000000\"}]",
						"proposal_id": "4",
						"voter":       "celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgVote,
				Height: 871324,
			},
			idx: testsuite.Ptr(0),
			votes: []*storage.Vote{
				{
					ProposalId: 4,
					Voter: &storage.Address{
						Address:    "celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r",
						Height:     871324,
						LastHeight: 871324,
						Balance:    storage.EmptyBalance(),
						Hash:       []byte{0x50, 0xa1, 0xec, 0xc6, 0x67, 0x0c, 0x9a, 0x72, 0x1f, 0x26, 0x7e, 0x08, 0xcd, 0x7b, 0x2b, 0xbb, 0x22, 0xfd, 0xe6, 0xc8},
					},
					Option: types.VoteOptionYes,
					Weight: decimal.RequireFromString("1.000000000000000000"),
					Height: 871324,
					Time:   ts,
				},
			},
		}, {
			name: "vote test 4",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 871324,
					Type:   "message",
					Data: map[string]any{
						"action":    "/cosmos.gov.v1.MsgVote",
						"sender":    "celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r",
						"module":    "gov",
						"msg_index": "0",
					},
				}, {
					Height: 871324,
					Type:   "proposal_vote",
					Data: map[string]any{
						"option":      "[{\"option\":1,\"weight\":\"0.500000000000000000\"},{\"option\":2,\"weight\":\"0.500000000000000000\"}]",
						"proposal_id": "4",
						"voter":       "celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgVote,
				Height: 871324,
			},
			idx: testsuite.Ptr(0),
			votes: []*storage.Vote{
				{
					ProposalId: 4,
					Voter: &storage.Address{
						Address:    "celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r",
						Height:     871324,
						LastHeight: 871324,
						Balance:    storage.EmptyBalance(),
						Hash:       []byte{0x50, 0xa1, 0xec, 0xc6, 0x67, 0x0c, 0x9a, 0x72, 0x1f, 0x26, 0x7e, 0x08, 0xcd, 0x7b, 0x2b, 0xbb, 0x22, 0xfd, 0xe6, 0xc8},
					},
					Option: types.VoteOptionYes,
					Weight: decimal.RequireFromString("0.500000000000000000"),
					Height: 871324,
					Time:   ts,
				}, {
					ProposalId: 4,
					Voter: &storage.Address{
						Address:    "celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r",
						Height:     871324,
						LastHeight: 871324,
						Balance:    storage.EmptyBalance(),
					},
					Option: types.VoteOptionAbstain,
					Weight: decimal.RequireFromString("0.500000000000000000"),
					Height: 871324,
					Time:   ts,
				},
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
			require.Len(t, tt.ctx.Votes, len(tt.votes))
			require.Equal(t, tt.votes, tt.ctx.Votes)
		})
	}
}
