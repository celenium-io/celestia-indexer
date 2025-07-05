// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
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

func Test_handleDelegate(t *testing.T) {
	ts := time.Now()
	tests := []struct {
		name       string
		ctx        *context.Context
		events     []storage.Event
		msg        *storage.Message
		idx        *int
		delegation *storage.Delegation
	}{
		{
			name: "top up",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 841682,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgDelegate",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "6745utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":  "6745utia",
						"spender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "6745utia",
						"recipient": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "6745utia",
						"delegator": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
						"validator": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "5690000utia",
						"spender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "5690000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "delegate",
					Data: map[string]any{
						"amount":     "5690000utia",
						"new_shares": "5690000.000000000000000000",
						"validator":  "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 841682,
				Time:   ts,
				Data: map[string]any{
					"Amount": map[string]string{
						"Amount": "5690000",
						"Denom":  "utia",
					},
					"DelegatorAddress": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					"ValidatorAddress": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
				},
			},
			idx: testsuite.Ptr(0),
			delegation: &storage.Delegation{
				Amount: decimal.RequireFromString("5690000"),
				Address: &storage.Address{
					Height:     841682,
					LastHeight: 841682,
					Address:    "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					Hash:       []byte{0xe7, 0xeb, 0x3b, 0x22, 0x85, 0x79, 0xa0, 0x93, 0xe0, 0x33, 0xd3, 0xc4, 0xf4, 0x6e, 0x46, 0xdb, 0xb8, 0x69, 0x97, 0xa9},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.RequireFromString("5690000"),
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
					Stake:             decimal.RequireFromString("5690000"),
				},
			},
		}, {
			name: "first time delegation",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 841682,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgDelegate",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "5690000utia",
						"spender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "5690000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "delegate",
					Data: map[string]any{
						"amount":     "5690000utia",
						"new_shares": "5690000.000000000000000000",
						"validator":  "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 841682,
				Time:   ts,
				Data: map[string]any{
					"Amount": map[string]string{
						"Amount": "5690000",
						"Denom":  "utia",
					},
					"DelegatorAddress": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					"ValidatorAddress": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
				},
			},
			idx: testsuite.Ptr(0),
			delegation: &storage.Delegation{
				Amount: decimal.RequireFromString("5690000"),
				Address: &storage.Address{
					Height:     841682,
					LastHeight: 841682,
					Address:    "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					Hash:       []byte{0xe7, 0xeb, 0x3b, 0x22, 0x85, 0x79, 0xa0, 0x93, 0xe0, 0x33, 0xd3, 0xc4, 0xf4, 0x6e, 0x46, 0xdb, 0xb8, 0x69, 0x97, 0xa9},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.RequireFromString("5690000"),
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
					Stake:             decimal.RequireFromString("5690000"),
				},
			},
		}, {
			name: "delegate with zero stake",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 105,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "24205utia",
						"spender": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "24205utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "24205utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "tx",
					Data: map[string]any{
						"fee":       "24205utia",
						"fee_payer": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "tx",
					Data: map[string]any{
						"acc_seq": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7/1",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "tx",
					Data: map[string]any{
						"signature": "UWAiCoF5Kgyp2o1R/ud25/azKfk5OEkp3ynOJZ6V79o3kRNEtn8jxvYRs8UGuex2Gy4eP7abdkuMMn8SpWpOmg==",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgDelegate",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0stake",
						"delegator": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
						"validator": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "100000000utia",
						"spender": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "100000000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "delegate",
					Data: map[string]any{
						"amount":     "100000000utia",
						"new_shares": "100000000.000000000000000000",
						"validator":  "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 105,
				Time:   ts,
				Data: map[string]any{
					"Amount": map[string]string{
						"Amount": "100000000",
						"Denom":  "utia",
					},
					"DelegatorAddress": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					"ValidatorAddress": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
				},
			},
			idx: testsuite.Ptr(7),
			delegation: &storage.Delegation{
				Amount: decimal.RequireFromString("100000000"),
				Address: &storage.Address{
					Height:     105,
					LastHeight: 105,
					Address:    "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					Hash:       []byte{0x2a, 0x30, 0x5a, 0xbd, 0x8d, 0x87, 0xca, 0xb0, 0xd2, 0xcc, 0xb0, 0xb3, 0x8f, 0x6c, 0xf3, 0x97, 0xf8, 0x1f, 0x74, 0x8f},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.RequireFromString("100000000"),
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
					Stake:             decimal.RequireFromString("100000000"),
				},
			},
		}, {
			name: "delegate v4",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 105,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{

						"action":    "/cosmos.staking.v1beta1.MsgDelegate",
						"module":    "staking",
						"msg_index": "0",
						"sender":    "celestia1768fz9r2c4w86kvckxfhpvha054rkc7rzwl83x",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":    "45000000utia",
						"msg_index": "0",
						"spender":   "celestia1768fz9r2c4w86kvckxfhpvha054rkc7rzwl83x",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":    "45000000utia",
						"msg_index": "0",
						"receiver":  "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 105,
					Time:   ts,
					Type:   "delegate",
					Data: map[string]any{
						"amount":     "45000000utia",
						"delegator":  "celestia1768fz9r2c4w86kvckxfhpvha054rkc7rzwl83x",
						"msg_index":  "0",
						"new_shares": "45000000.000000000000000000",
						"validator":  "celestiavaloper1jt9w26mpxxjsk63mvd4m2ynj0af09cslh5d096",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 105,
				Time:   ts,
				Data: map[string]any{
					"Amount": map[string]string{
						"Amount": "45000000",
						"Denom":  "utia",
					},
					"DelegatorAddress": "celestia1768fz9r2c4w86kvckxfhpvha054rkc7rzwl83x",
					"ValidatorAddress": "celestiavaloper1jt9w26mpxxjsk63mvd4m2ynj0af09cslh5d096",
				},
			},
			idx: testsuite.Ptr(0),
			delegation: &storage.Delegation{
				Amount: decimal.RequireFromString("45000000"),
				Address: &storage.Address{
					Height:     105,
					LastHeight: 105,
					Address:    "celestia1768fz9r2c4w86kvckxfhpvha054rkc7rzwl83x",
					Hash:       []byte{0xf6, 0x8e, 0x91, 0x14, 0x6a, 0xc5, 0x5c, 0x7d, 0x59, 0x98, 0xb1, 0x93, 0x70, 0xb2, 0xfd, 0x7d, 0x2a, 0x3b, 0x63, 0xc3},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.RequireFromString("45000000"),
					},
				},
				Validator: &storage.Validator{
					Address:           "celestiavaloper1jt9w26mpxxjsk63mvd4m2ynj0af09cslh5d096",
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
			err := handleDelegate(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
			require.EqualValues(t, 1, tt.ctx.Delegations.Len())
			err = tt.ctx.Delegations.Range(func(_ string, value *storage.Delegation) (error, bool) {
				require.Equal(t, tt.delegation, value)
				return nil, true
			})
			require.NoError(t, err)
		})
	}
}
