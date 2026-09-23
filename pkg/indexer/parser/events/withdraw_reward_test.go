// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func Test_handleWithdrawRewards(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msgs   []*storage.Message
		idx    int
	}{
		{
			name: "test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 848613,
					Type:   "message",
					Data: map[string]string{
						"action":    "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
						"sender":    "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p",
						"module":    "distribution",
						"msg_index": "0",
					},
				}, {
					Height: 848613,
					Type:   "withdraw_rewards",
					Data: map[string]string{
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
			idx: 0,
		}, {
			name: "test 2",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 848613,
					Type:   "message",
					Data: map[string]string{
						"action": "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
					},
				}, {
					Height: 848613,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":    "0utia",
						"validator": "celestiavaloper1u5pshtqpexjmuudrvq6q335qym2zggzhyp5ee8",
						"delegator": "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p",
					},
				}, {
					Height: 848613,
					Type:   "message",
					Data: map[string]string{
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
			idx: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.msgs {
				c := NewCursor(tt.events)
				c.Skip(tt.idx)
				err := handleWithdrawDelegatorRewards(tt.ctx, c, tt.msgs[i])
				require.NoError(t, err)
			}
		})
	}
}

// The event carries the full withdrawn amount, which leaves the validator's reward pool,
// so the context must hold it as a negative delta. Regression: the guard used to be
// `rewards.Delegator != ""` while the decoder rejects an empty delegator, so the whole
// function returned before doing anything.
func Test_parseWithdrawRewards(t *testing.T) {
	const (
		validator = "celestiavaloper1u5pshtqpexjmuudrvq6q335qym2zggzhyp5ee8"
		delegator = "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p"
	)

	msg := &storage.Message{
		Type:   types.MsgWithdrawDelegatorReward,
		Height: 848613,
		Time:   time.Now(),
	}

	t.Run("withdrawn amount is subtracted from rewards", func(t *testing.T) {
		ctx := context.NewContext()

		err := parseWithdrawRewards(ctx, msg, map[string]string{
			"amount":    "1314utia",
			"validator": validator,
			"delegator": delegator,
		})
		require.NoError(t, err)

		val, ok := ctx.Validators.Get(validator)
		require.True(t, ok, "validator must reach the context")
		require.Equal(t, "-1314", val.Rewards.String())

		require.Len(t, ctx.StakingLogs, 1)
		require.Equal(t, types.StakingLogTypeRewards, ctx.StakingLogs[0].Type)
		require.Equal(t, "-1314", ctx.StakingLogs[0].Change.String())
		require.Equal(t, msg.Height, ctx.StakingLogs[0].Height)

		_, ok = ctx.Addresses.Get(delegator)
		require.True(t, ok, "reward receiver must reach the context")
	})

	t.Run("two withdrawals in one block accumulate", func(t *testing.T) {
		ctx := context.NewContext()

		for _, amount := range []string{"1314utia", "686utia"} {
			require.NoError(t, parseWithdrawRewards(ctx, msg, map[string]string{
				"amount":    amount,
				"validator": validator,
				"delegator": delegator,
			}))
		}

		val, ok := ctx.Validators.Get(validator)
		require.True(t, ok)
		require.Equal(t, "-2000", val.Rewards.String())
		require.Len(t, ctx.StakingLogs, 2)
	})

	// CIP-30 claims rewards on every delegation change; the emitted amount is zero there
	t.Run("zero amount writes no log", func(t *testing.T) {
		ctx := context.NewContext()

		err := parseWithdrawRewards(ctx, msg, map[string]string{
			"amount":    "0utia",
			"validator": validator,
			"delegator": delegator,
		})
		require.NoError(t, err)
		require.Empty(t, ctx.StakingLogs)
	})

	t.Run("event without a delegator is an error", func(t *testing.T) {
		ctx := context.NewContext()

		err := parseWithdrawRewards(ctx, msg, map[string]string{
			"amount":    "1314utia",
			"validator": validator,
		})
		require.Error(t, err)
	})
}
