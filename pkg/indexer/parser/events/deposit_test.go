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

func Test_handleDeposit(t *testing.T) {
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
			name: "deposit test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.gov.v1.MsgDeposit",
					},
				}, {
					Height: 1745041,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "9990000000utia",
						"spender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 1745041,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "9990000000utia",
						"receiver": "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
					},
				}, {
					Height: 1745041,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "9990000000utia",
						"recipient": "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
						"sender":    "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 1745041,
					Type:   "proposal_deposit",
					Data: map[string]any{
						"amount":      "9990000000utia",
						"proposal_id": "2",
					},
				}, {
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"module": "governance",
						"sender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 1745041,
					Type:   "proposal_deposit",
					Data: map[string]any{
						"voting_period_start": "2",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDeposit,
				Height: 1745041,
			},
			idx: testsuite.Ptr(0),
			proposal: storage.Proposal{
				Id:             2,
				ActivationTime: &ts,
				Status:         types.ProposalStatusActive,
				Deposit:        decimal.RequireFromString("9990000000"),
			},
		}, {
			name: "deposit test 2",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.gov.v1.MsgDeposit",
					},
				}, {
					Height: 1745041,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "9990000000utia",
						"spender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 1745041,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "9990000000utia",
						"receiver": "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
					},
				}, {
					Height: 1745041,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "9990000000utia",
						"recipient": "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
						"sender":    "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				}, {
					Height: 1745041,
					Type:   "proposal_deposit",
					Data: map[string]any{
						"amount":      "9990000000utia",
						"proposal_id": "2",
					},
				}, {
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"module": "governance",
						"sender": "celestia1jkuw8rxxrsgn9pq009987kzelkp46cgcczuxp5",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDeposit,
				Height: 1745041,
			},
			idx: testsuite.Ptr(0),
			proposal: storage.Proposal{
				Id:      2,
				Status:  types.ProposalStatusInactive,
				Deposit: decimal.RequireFromString("9990000000"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ctx.Block = &storage.Block{
				Time: ts,
			}
			err := handleDeposit(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
			require.NotNil(t, tt.msg.Proposal)
			require.Equal(t, tt.proposal, *tt.msg.Proposal)
		})
	}
}
