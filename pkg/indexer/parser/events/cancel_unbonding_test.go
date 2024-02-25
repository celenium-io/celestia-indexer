// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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

func Test_handleCancelUnbonding(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
		idx    *int
	}{
		{
			name: "test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 844287,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation",
					},
				}, {
					Height: 844287,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "1314utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844287,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "1314utia",
						"receiver": "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9",
					},
				}, {
					Height: 844287,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "1314utia",
						"recipient": "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844287,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844287,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "1314utia",
						"delegator": "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9",
						"validator": "celestiavaloper1qe8uuf5x69c526h4nzxwv4ltftr73v7q5qhs58",
					},
				}, {
					Height: 844287,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "45000000utia",
						"spender": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 844287,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "45000000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844287,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "45000000utia",
						"recipient": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
						"sender":    "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 844287,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 844287,
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
				Type:   types.MsgDelegate,
				Height: 844287,
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleCancelUnbonding(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
		})
	}
}
