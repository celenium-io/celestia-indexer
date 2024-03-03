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

func Test_handleRedelegate(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
		idx    *int
	}{
		{
			name: "redelegate test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 841682,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgBeginRedelegate",
					},
				}, {
					Height: 841682,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "384192utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 841682,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":  "384192utia",
						"spender": "celestia1kj4m7nhhcr5f5jpcdaec5ymhepc3x22yvtjh7j",
					},
				}, {
					Height: 841682,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "384192utia",
						"recipient": "celestia1kj4m7nhhcr5f5jpcdaec5ymhepc3x22yvtjh7j",
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
						"amount":    "384192utia",
						"delegator": "celestia1kj4m7nhhcr5f5jpcdaec5ymhepc3x22yvtjh7j",
						"validator": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 841682,
					Type:   "redelegate",
					Data: map[string]any{
						"amount":                "20000000utia",
						"completion_time":       "2024-03-15T00:04:38Z",
						"destination_validator": "celestiavaloper107lwx458gy345ag2afx9a7e2kkl7x49y3433gj",
						"source_validator":      "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 841682,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1kj4m7nhhcr5f5jpcdaec5ymhepc3x22yvtjh7j",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 841682,
				Data: map[string]any{
					"Amount": map[string]string{
						"Amount": "20000000",
						"Denom":  "utia",
					},
					"DelegatorAddress":    "celestia1kj4m7nhhcr5f5jpcdaec5ymhepc3x22yvtjh7j",
					"ValidatorDstAddress": "celestiavaloper107lwx458gy345ag2afx9a7e2kkl7x49y3433gj",
					"ValidatorSrcAddress": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "redelegate test 2",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 241,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgBeginRedelegate",
					},
				}, {
					Height: 241,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0stake",
						"delegator": "celestia1ze2ye5u5k3qdlexvt2e0nn0508p04094r9atu2",
						"validator": "celestiavaloper1e2p4u5vqwgum7pm9vhp0yjvl58gvhfc6yfatw4",
					},
				}, {
					Height: 241,
					Type:   "redelegate",
					Data: map[string]any{
						"amount":                "1000utia",
						"completion_time":       "2023-11-21T14:50:46Z",
						"destination_validator": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
						"source_validator":      "celestiavaloper1e2p4u5vqwgum7pm9vhp0yjvl58gvhfc6yfatw4",
					},
				}, {
					Height: 241,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1ze2ye5u5k3qdlexvt2e0nn0508p04094r9atu2",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 241,
				Data: map[string]any{
					"Amount": map[string]string{
						"Amount": "1000",
						"Denom":  "utia",
					},
					"DelegatorAddress":    "celestia1ze2ye5u5k3qdlexvt2e0nn0508p04094r9atu2",
					"ValidatorDstAddress": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
					"ValidatorSrcAddress": "celestiavaloper1e2p4u5vqwgum7pm9vhp0yjvl58gvhfc6yfatw4",
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "redelegate test 3",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 315,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgBeginRedelegate",
					},
				}, {
					Height: 315,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "44283utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "44283utia",
						"receiver": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
					},
				}, {
					Height: 315,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "44283utia",
						"recipient": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "44283utia",
						"delegator": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
						"validator": "celestiavaloper1uwmf03ke52vld2sa9khs0nslpgzwsm5xs5e4pn",
					},
				}, {
					Height: 315,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "3953utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "3953utia",
						"receiver": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
					},
				}, {
					Height: 315,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "3953utia",
						"recipient": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "3953utia",
						"delegator": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
						"validator": "celestiavaloper19f0w9svr905fhefusyx4z8sf83j6et0g57nch8",
					},
				}, {
					Height: 315,
					Type:   "redelegate",
					Data: map[string]any{
						"amount":                "125000000utia",
						"completion_time":       "2023-11-21T15:05:39Z",
						"destination_validator": "celestiavaloper19f0w9svr905fhefusyx4z8sf83j6et0g57nch8",
						"source_validator":      "celestiavaloper1uwmf03ke52vld2sa9khs0nslpgzwsm5xs5e4pn",
					},
				}, {
					Height: 315,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 315,
				Data: map[string]any{
					"Amount": map[string]string{
						"Amount": "125000000",
						"Denom":  "utia",
					},
					"DelegatorAddress":    "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
					"ValidatorDstAddress": "celestiavaloper19f0w9svr905fhefusyx4z8sf83j6et0g57nch8",
					"ValidatorSrcAddress": "celestiavaloper1uwmf03ke52vld2sa9khs0nslpgzwsm5xs5e4pn",
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "redelegate test 4",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 584,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgBeginRedelegate",
					},
				}, {
					Height: 584,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0stake",
						"delegator": "celestia18ccg4scfdee9qvcnxxx4wutjt69sct2jshau9f",
						"validator": "celestiavaloper19mm3s0y676453w3ja58d376ysd84wlf6hq8ae3",
					},
				}, {
					Height: 584,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "250000000utia",
						"spender": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 584,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "250000000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 584,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "250000000utia",
						"recipient": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
						"sender":    "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 584,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 584,
					Type:   "redelegate",
					Data: map[string]any{
						"amount":                "250000000utia",
						"completion_time":       "0001-01-01T00:00:00Z",
						"destination_validator": "celestiavaloper1yecxnyegvgm5dwsx0r3jsgr74ju6mlxdwkxx8g",
						"source_validator":      "celestiavaloper19mm3s0y676453w3ja58d376ysd84wlf6hq8ae3",
					},
				}, {
					Height: 584,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia18ccg4scfdee9qvcnxxx4wutjt69sct2jshau9f",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 315,
				Data: map[string]any{
					"Amount": map[string]string{
						"Amount": "250000000",
						"Denom":  "utia",
					},
					"DelegatorAddress":    "celestia18ccg4scfdee9qvcnxxx4wutjt69sct2jshau9f",
					"ValidatorDstAddress": "celestiavaloper1yecxnyegvgm5dwsx0r3jsgr74ju6mlxdwkxx8g",
					"ValidatorSrcAddress": "celestiavaloper19mm3s0y676453w3ja58d376ysd84wlf6hq8ae3",
				},
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleRedelegate(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
		})
	}
}
