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

func Test_handleUndelegate(t *testing.T) {
	ts := time.Now()
	tests := []struct {
		name         string
		ctx          *context.Context
		events       []storage.Event
		msg          *storage.Message
		idx          *int
		undelegation *storage.Undelegation
	}{
		{
			name: "undelegate",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 844186,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgUndelegate",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "3105utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "3105utia",
						"receiver": "celestia144tkegx6vw9pr6tx0gg68zmy57vcjpapxgwn4q",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "3105utia",
						"recipient": "celestia144tkegx6vw9pr6tx0gg68zmy57vcjpapxgwn4q",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "3105utia",
						"delegator": "celestia144tkegx6vw9pr6tx0gg68zmy57vcjpapxgwn4q",
						"validator": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "144000000utia",
						"spender": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "144000000utia",
						"receiver": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "144000000utia",
						"recipient": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"sender":    "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "unbond",
					Data: map[string]any{
						"amount":          "144000000utia",
						"completion_time": "2024-03-15T00:25:17Z",
						"validator":       "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 844186,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia144tkegx6vw9pr6tx0gg68zmy57vcjpapxgwn4q",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUndelegate,
				Height: 844186,
				Time:   ts,
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
			undelegation: &storage.Undelegation{
				Height:         844186,
				Time:           ts,
				Amount:         decimal.RequireFromString("144000000"),
				CompletionTime: time.Date(2024, 3, 15, 0, 25, 17, 0, time.UTC),
				Address: &storage.Address{
					Height:     844186,
					LastHeight: 844186,
					Address:    "celestia144tkegx6vw9pr6tx0gg68zmy57vcjpapxgwn4q",
					Hash:       []byte{0xad, 0x57, 0x6c, 0xa0, 0xda, 0x63, 0x8a, 0x11, 0xe9, 0x66, 0x7a, 0x11, 0xa3, 0x8b, 0x64, 0xa7, 0x99, 0x89, 0x07, 0xa1},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.RequireFromString("-144000000"),
						Spendable: decimal.Zero,
						Unbonding: decimal.RequireFromString("144000000"),
					},
				},
				Validator: &storage.Validator{
					Address:           "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
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
					Stake:             decimal.RequireFromString("-144000000"),
				},
			},
		}, {
			name: "undelegate zero",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 75,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgUndelegate",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0stake",
						"delegator": "celestia15e9hkqujx0c8m464w6a35glc8vtqdrxnshcmwq",
						"validator": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "30000utia",
						"spender": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "30000utia",
						"receiver": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "30000utia",
						"recipient": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"sender":    "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "unbond",
					Data: map[string]any{
						"amount":          "30000utia",
						"completion_time": "2023-11-21T14:16:41Z",
						"validator":       "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia15e9hkqujx0c8m464w6a35glc8vtqdrxnshcmwq",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUndelegate,
				Height: 75,
				Time:   ts,
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
			undelegation: &storage.Undelegation{
				Height:         75,
				Time:           ts,
				Amount:         decimal.RequireFromString("30000"),
				CompletionTime: time.Date(2023, 11, 21, 14, 16, 41, 0, time.UTC),
				Address: &storage.Address{
					Height:     75,
					LastHeight: 75,
					Address:    "celestia15e9hkqujx0c8m464w6a35glc8vtqdrxnshcmwq",
					Hash:       []byte{0xa6, 0x4b, 0x7b, 0x03, 0x92, 0x33, 0xf0, 0x7d, 0xd7, 0x55, 0x76, 0xbb, 0x1a, 0x23, 0xf8, 0x3b, 0x16, 0x06, 0x8c, 0xd3},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.RequireFromString("-30000"),
						Spendable: decimal.Zero,
						Unbonding: decimal.RequireFromString("30000"),
					},
				},
				Validator: &storage.Validator{
					Address:           "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
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
					Stake:             decimal.RequireFromString("-30000"),
				},
			},
		}, {
			name: "undelegate v4",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 75,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action":    "/cosmos.staking.v1beta1.MsgUndelegate",
						"module":    "staking",
						"msg_index": "0",
						"sender":    "celestia1xs55snr6lxsalaqrwc63cxlmgn437zzv2gew35",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0utia",
						"delegator": "celestia1xs55snr6lxsalaqrwc63cxlmgn437zzv2gew35",
						"msg_index": "0",
						"validator": "celestiavaloper109nzhf6fvqvfan3tayzc8cywcsk6a5q45lmk5s",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":    "1000000utia",
						"msg_index": "0",
						"spender":   "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":    "1000000utia",
						"msg_index": "0",
						"receiver":  "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "1000000utia",
						"msg_index": "0",
						"recipient": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"sender":    "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"msg_index": "0",
						"sender":    "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 75,
					Time:   ts,
					Type:   "unbond",
					Data: map[string]any{
						"amount":          "1000000utia",
						"completion_time": "2025-07-23T11:56:30Z",
						"delegator":       "celestia1xs55snr6lxsalaqrwc63cxlmgn437zzv2gew35",
						"msg_index":       "0",
						"validator":       "celestiavaloper109nzhf6fvqvfan3tayzc8cywcsk6a5q45lmk5s",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUndelegate,
				Height: 75,
				Time:   ts,
				Data: map[string]any{
					"Amount": map[string]any{
						"Amount": "1000000",
						"Denom":  "utia",
					},
					"DelegatorAddress": "celestia1xs55snr6lxsalaqrwc63cxlmgn437zzv2gew35",
					"ValidatorAddress": "celestiavaloper109nzhf6fvqvfan3tayzc8cywcsk6a5q45lmk5s",
				},
			},
			idx: testsuite.Ptr(0),
			undelegation: &storage.Undelegation{
				Height:         75,
				Time:           ts,
				Amount:         decimal.RequireFromString("1000000"),
				CompletionTime: time.Date(2025, 7, 23, 11, 56, 30, 0, time.UTC),
				Address: &storage.Address{
					Height:     75,
					LastHeight: 75,
					Address:    "celestia1xs55snr6lxsalaqrwc63cxlmgn437zzv2gew35",
					Hash:       []byte{0x34, 0x29, 0x48, 0x4c, 0x7a, 0xf9, 0xa1, 0xdf, 0xf4, 0x03, 0x76, 0x35, 0x1c, 0x1b, 0xfb, 0x44, 0xeb, 0x1f, 0x08, 0x4c},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.RequireFromString("-1000000"),
						Spendable: decimal.Zero,
						Unbonding: decimal.RequireFromString("1000000"),
					},
				},
				Validator: &storage.Validator{
					Address:           "celestiavaloper109nzhf6fvqvfan3tayzc8cywcsk6a5q45lmk5s",
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
					Stake:             decimal.RequireFromString("-1000000"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ctx.Block = &storage.Block{
				Time: ts,
			}
			err := handleUndelegate(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
			require.Len(t, tt.ctx.Undelegations, 1)
			require.Equal(t, *tt.undelegation, tt.ctx.Undelegations[0])
		})
	}
}
