// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
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

func Test_handleDelegate(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
		idx    *int
	}{
		{
			name: "top up",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 841682,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgDelegate",
					},
				}, {
					Height: 841682,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "6745utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 841682,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":  "6745utia",
						"spender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				}, {
					Height: 841682,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "6745utia",
						"recipient": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 841682,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 841682,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "6745utia",
						"delegator": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
						"validator": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 841682,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "5690000utia",
						"spender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				}, {
					Height: 841682,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "5690000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 841682,
					Type:   "delegate",
					Data: map[string]any{
						"amount":     "5690000utia",
						"new_shares": "5690000.000000000000000000",
						"validator":  "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 841682,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 841682,
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
		}, {
			name: "first time delegation",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 841682,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgDelegate",
					},
				}, {
					Height: 841682,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "5690000utia",
						"spender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				}, {
					Height: 841682,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "5690000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 841682,
					Type:   "delegate",
					Data: map[string]any{
						"amount":     "5690000utia",
						"new_shares": "5690000.000000000000000000",
						"validator":  "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 841682,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 841682,
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
		}, {
			name: "delegate with zero stake",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 105,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "24205utia",
						"spender": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				}, {
					Height: 105,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "24205utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 105,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "24205utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				}, {
					Height: 105,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				}, {
					Height: 105,
					Type:   "tx",
					Data: map[string]any{
						"fee":       "24205utia",
						"fee_payer": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				}, {
					Height: 105,
					Type:   "tx",
					Data: map[string]any{
						"acc_seq": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7/1",
					},
				}, {
					Height: 105,
					Type:   "tx",
					Data: map[string]any{
						"signature": "UWAiCoF5Kgyp2o1R/ud25/azKfk5OEkp3ynOJZ6V79o3kRNEtn8jxvYRs8UGuex2Gy4eP7abdkuMMn8SpWpOmg==",
					},
				}, {
					Height: 105,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgDelegate",
					},
				}, {
					Height: 105,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0stake",
						"delegator": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
						"validator": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 105,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "100000000utia",
						"spender": "celestia19gc940vdsl9tp5kvkzec7m8njlup7ay0frlfu7",
					},
				}, {
					Height: 105,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "100000000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 105,
					Type:   "delegate",
					Data: map[string]any{
						"amount":     "100000000utia",
						"new_shares": "100000000.000000000000000000",
						"validator":  "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 105,
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleDelegate(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
		})
	}
}
