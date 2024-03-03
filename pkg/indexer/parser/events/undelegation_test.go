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

func Test_handleUndelegate(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
		idx    *int
	}{
		{
			name: "undelegate",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 844186,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgUndelegate",
					},
				}, {
					Height: 844186,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "3105utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844186,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "3105utia",
						"receiver": "celestia144tkegx6vw9pr6tx0gg68zmy57vcjpapxgwn4q",
					},
				}, {
					Height: 844186,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "3105utia",
						"recipient": "celestia144tkegx6vw9pr6tx0gg68zmy57vcjpapxgwn4q",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844186,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844186,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "3105utia",
						"delegator": "celestia144tkegx6vw9pr6tx0gg68zmy57vcjpapxgwn4q",
						"validator": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 844186,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "144000000utia",
						"spender": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844186,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "144000000utia",
						"receiver": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 844186,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "144000000utia",
						"recipient": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"sender":    "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844186,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844186,
					Type:   "unbond",
					Data: map[string]any{
						"amount":          "144000000utia",
						"completion_time": "2024-03-15T00:25:17Z",
						"validator":       "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 844186,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia144tkegx6vw9pr6tx0gg68zmy57vcjpapxgwn4q",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 844186,
				Data: map[string]any{
					"Amount": map[string]string{
						"Amount": "144000000",
						"Denom":  "utia",
					},
					"DelegatorAddress": "celestia144tkegx6vw9pr6tx0gg68zmy57vcjpapxgwn4q",
					"ValidatorAddress": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "undelegate zero",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 75,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgUndelegate",
					},
				}, {
					Height: 75,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0stake",
						"delegator": "celestia15e9hkqujx0c8m464w6a35glc8vtqdrxnshcmwq",
						"validator": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
					},
				}, {
					Height: 75,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "30000utia",
						"spender": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 75,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "30000utia",
						"receiver": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 75,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "30000utia",
						"recipient": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"sender":    "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 75,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 75,
					Type:   "unbond",
					Data: map[string]any{
						"amount":          "30000utia",
						"completion_time": "2023-11-21T14:16:41Z",
						"validator":       "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
					},
				}, {
					Height: 75,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia15e9hkqujx0c8m464w6a35glc8vtqdrxnshcmwq",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 75,
				Data: map[string]any{
					"Amount": map[string]string{
						"Amount": "30000",
						"Denom":  "utia",
					},
					"DelegatorAddress": "celestia15e9hkqujx0c8m464w6a35glc8vtqdrxnshcmwq",
					"ValidatorAddress": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
				},
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleUndelegate(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
		})
	}
}
