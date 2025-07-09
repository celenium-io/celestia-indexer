// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func Test_handleWithdrawRewards(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msgs   []*storage.Message
		idx    *int
	}{
		{
			name: "test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 848613,
					Type:   "message",
					Data: map[string]any{
						"action":    "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
						"sender":    "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p",
						"module":    "distribution",
						"msg_index": "0",
					},
				}, {
					Height: 848613,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0utia",
						"validator": "celestiavaloper1u5pshtqpexjmuudrvq6q335qym2zggzhyp5ee8",
						"delegator": "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p",
						"msg_index": "0",
					},
				},
			},
			msgs: []*storage.Message{
				{
					Type:   types.MsgWithdrawDelegatorReward,
					Height: 848613,
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "test 2",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 848613,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
					},
				}, {
					Height: 848613,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0utia",
						"validator": "celestiavaloper1u5pshtqpexjmuudrvq6q335qym2zggzhyp5ee8",
						"delegator": "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p",
					},
				}, {
					Height: 848613,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p",
						"module": "distribution",
					},
				},
			},
			msgs: []*storage.Message{
				{
					Type:   types.MsgWithdrawDelegatorReward,
					Height: 848613,
				},
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.msgs {
				err := handleWithdrawDelegatorRewards(tt.ctx, tt.events, tt.msgs[i], tt.idx)
				require.NoError(t, err)
			}
		})
	}
}
