// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test_handleCancelUnbonding(t *testing.T) {
	ts := time.Now()
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
		idx    *int
		cancel *storage.Undelegation
	}{
		{
			name: "test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 844287,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation",
					},
				}, {
					Height: 844287,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "1314utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844287,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "1314utia",
						"receiver": "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9",
					},
				}, {
					Height: 844287,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "1314utia",
						"recipient": "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844287,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844287,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "1314utia",
						"delegator": "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9",
						"validator": "celestiavaloper1qe8uuf5x69c526h4nzxwv4ltftr73v7q5qhs58",
					},
				}, {
					Height: 844287,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "45000000utia",
						"spender": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 844287,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "45000000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844287,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "45000000utia",
						"recipient": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
						"sender":    "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 844287,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 844287,
					Time:   ts,
					Type:   "cancel_unbonding_delegation",
					Data: map[string]any{
						"amount":          "45000000utia",
						"creation_height": "842069",
						"delegator":       "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9",
						"validator":       "celestiavaloper1qe8uuf5x69c526h4nzxwv4ltftr73v7q5qhs58",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgCancelUnbondingDelegation,
				Height: 844287,
				Time:   ts,
				Data: types.PackedBytes{

					"Amount": map[string]any{
						"Amount": "45000000",
						"Denom":  "utia",
					},
					"CreationHeight":   842069,
					"DelegatorAddress": "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9",
					"ValidatorAddress": "celestiavaloper1qe8uuf5x69c526h4nzxwv4ltftr73v7q5qhs58",
				},
			},
			idx: testsuite.Ptr(0),
			cancel: &storage.Undelegation{
				Height: 844287,
				Time:   ts,
				Amount: decimal.RequireFromString("45000000"),
				Address: &storage.Address{
					Address:    "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9",
					Height:     844287,
					LastHeight: 844287,
					Hash:       []byte{0xfd, 0x86, 0xd3, 0xeb, 0x83, 0x1e, 0xef, 0xb8, 0x31, 0x02, 0x49, 0x11, 0x28, 0x23, 0x8c, 0x8e, 0x10, 0xc7, 0xbe, 0x3a},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.RequireFromString("45000000"),
						Unbonding: decimal.RequireFromString("-45000000"),
					},
				},
				Validator: &storage.Validator{
					Address:           "celestiavaloper1qe8uuf5x69c526h4nzxwv4ltftr73v7q5qhs58",
					Moniker:           storage.DoNotModify,
					Website:           storage.DoNotModify,
					Identity:          storage.DoNotModify,
					Contacts:          storage.DoNotModify,
					Details:           storage.DoNotModify,
					Rate:              decimal.Zero,
					MaxRate:           decimal.Zero,
					MaxChangeRate:     decimal.Zero,
					MinSelfDelegation: decimal.Zero,
					Rewards:           decimal.Zero,
					Commissions:       decimal.Zero,
					Stake:             decimal.RequireFromString("45000000"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleCancelUnbonding(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
			require.Len(t, tt.ctx.CancelUnbonding, 1)
			require.Equal(t, *tt.cancel, tt.ctx.CancelUnbonding[0])
		})
	}
}
