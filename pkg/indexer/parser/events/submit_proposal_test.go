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

func Test_handleSubmitProposal(t *testing.T) {
	ts := time.Now()
	tests := []struct {
		name     string
		ctx      *context.Context
		events   []storage.Event
		msg      *storage.Message
		idx      *int
		proposal storage.Proposal
	}{
		{
			name: "submit_proposal test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 3648325,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.gov.v1beta1.MsgSubmitProposal",
					},
				}, {
					Height: 3648325,
					Type:   "submit_proposal",
					Data: map[string]any{
						"proposal_id":       "5",
						"proposal_messages": ",/cosmos.gov.v1.MsgExecLegacyContent",
					},
				}, {
					Height: 3648325,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "10000000000utia",
						"spender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 3648325,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "10000000000utia",
						"receiver": "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
					},
				}, {
					Height: 3648325,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "10000000000utia",
						"recipient": "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
						"sender":    "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 3648325,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 3648325,
					Type:   "proposal_deposit",
					Data: map[string]any{
						"amount":      "10000000000utia",
						"proposal_id": "5",
					},
				}, {
					Height: 3648325,
					Type:   "message",
					Data: map[string]any{
						"module": "governance",
						"sender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 3648325,
					Type:   "submit_proposal",
					Data: map[string]any{
						"voting_period_start": "5",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgSubmitProposal,
				Height: 3648325,
				Proposal: &storage.Proposal{
					Status: types.ProposalStatusInactive,
				},
			},
			idx: testsuite.Ptr(0),
			proposal: storage.Proposal{
				Id:             5,
				ActivationTime: &ts,
				Status:         types.ProposalStatusActive,
				Deposit:        decimal.RequireFromString("10000000000"),
			},
		}, {
			name: "submit_proposal test 2",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 3648325,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.gov.v1beta1.MsgSubmitProposal",
					},
				}, {
					Height: 3648325,
					Type:   "submit_proposal",
					Data: map[string]any{
						"proposal_id":       "5",
						"proposal_messages": ",/cosmos.gov.v1.MsgExecLegacyContent",
					},
				}, {
					Height: 3648325,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "10000000000utia",
						"spender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 3648325,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "10000000000utia",
						"receiver": "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
					},
				}, {
					Height: 3648325,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "10000000000utia",
						"recipient": "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
						"sender":    "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 3648325,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 3648325,
					Type:   "proposal_deposit",
					Data: map[string]any{
						"amount":      "10000000000utia",
						"proposal_id": "5",
					},
				}, {
					Height: 3648325,
					Type:   "message",
					Data: map[string]any{
						"module": "governance",
						"sender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgSubmitProposal,
				Height: 3648325,
				Proposal: &storage.Proposal{
					Status: types.ProposalStatusInactive,
				},
			},
			idx: testsuite.Ptr(0),
			proposal: storage.Proposal{
				Id:      5,
				Status:  types.ProposalStatusInactive,
				Deposit: decimal.RequireFromString("10000000000"),
			},
		}, {
			name: "submit_proposal test 3",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 58507,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.gov.v1.MsgSubmitProposal",
					},
				}, {
					Height: 58507,
					Type:   "submit_proposal",
					Data: map[string]any{
						"proposal_id":       "2",
						"proposal_messages": ",/cosmos.gov.v1beta1.MsgSubmitProposal",
					},
				}, {
					Height: 58507,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "",
						"spender": "celestia17adsjkuecgjheugrdrwdqv9uh3qkrfmj9xzawx",
					},
				}, {
					Height: 58507,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "",
						"receiver": "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
					},
				}, {
					Height: 58507,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "",
						"recipient": "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
						"sender":    "celestia17adsjkuecgjheugrdrwdqv9uh3qkrfmj9xzawx",
					},
				}, {
					Height: 58507,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia17adsjkuecgjheugrdrwdqv9uh3qkrfmj9xzawx",
					},
				}, {
					Height: 58507,
					Type:   "proposal_deposit",
					Data: map[string]any{
						"amount":      "",
						"proposal_id": "2",
					},
				}, {
					Height: 58507,
					Type:   "message",
					Data: map[string]any{
						"module": "governance",
						"sender": "celestia17adsjkuecgjheugrdrwdqv9uh3qkrfmj9xzawx",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgSubmitProposal,
				Height: 58507,
				Proposal: &storage.Proposal{
					Status: types.ProposalStatusInactive,
				},
			},
			idx: testsuite.Ptr(0),
			proposal: storage.Proposal{
				Id:      2,
				Status:  types.ProposalStatusInactive,
				Deposit: decimal.Zero,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ctx.Block = &storage.Block{
				Time: ts,
			}
			err := handleSubmitProposal(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
			require.NotNil(t, tt.msg.Proposal)
			require.Equal(t, tt.proposal, *tt.msg.Proposal)
		})
	}
}
