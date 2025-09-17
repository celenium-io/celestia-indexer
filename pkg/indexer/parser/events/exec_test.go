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

func Test_handleExec(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
		idx    *int
	}{
		{
			name: "multiple delegations",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 844359,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.authz.v1beta1.MsgExec",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "101775utia",
						"authz_msg_index": "0",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "101775utia",
						"authz_msg_index": "0",
						"receiver":        "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
					},
				}, {
					Height: 844359,
					Type:   "transfer",
					Data: map[string]any{
						"amount":          "101775utia",
						"authz_msg_index": "0",
						"recipient":       "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "0",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":          "101775utia",
						"authz_msg_index": "0",
						"delegator":       "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "101774utia",
						"authz_msg_index": "0",
						"spender":         "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "101774utia",
						"authz_msg_index": "0",
						"receiver":        "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844359,
					Type:   "delegate",
					Data: map[string]any{
						"amount":          "101774utia",
						"authz_msg_index": "0",
						"new_shares":      "101774.000000000000000000",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "0",
						"module":          "staking",
						"sender":          "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "191045utia",
						"authz_msg_index": "1",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "191045utia",
						"authz_msg_index": "1",
						"receiver":        "celestia1gfvu3xpgze2jy20cy4lcfeq2qj3rww0a8cwuap",
					},
				}, {
					Height: 844359,
					Type:   "transfer",
					Data: map[string]any{
						"amount":          "191045utia",
						"authz_msg_index": "1",
						"recipient":       "celestia1gfvu3xpgze2jy20cy4lcfeq2qj3rww0a8cwuap",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "1",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":          "191045utia",
						"authz_msg_index": "1",
						"delegator":       "celestia1gfvu3xpgze2jy20cy4lcfeq2qj3rww0a8cwuap",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "191018utia",
						"authz_msg_index": "1",
						"spender":         "celestia1gfvu3xpgze2jy20cy4lcfeq2qj3rww0a8cwuap",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "191018utia",
						"authz_msg_index": "1",
						"receiver":        "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844359,
					Type:   "delegate",
					Data: map[string]any{
						"amount":          "191018utia",
						"authz_msg_index": "1",
						"new_shares":      "191018.000000000000000000",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "1",
						"module":          "staking",
						"sender":          "celestia1gfvu3xpgze2jy20cy4lcfeq2qj3rww0a8cwuap",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "106033utia",
						"authz_msg_index": "2",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "106033utia",
						"authz_msg_index": "2",
						"receiver":        "celestia1ghwz05j5s52nyvzau08eg9rqvkzgq72r92k56d",
					},
				}, {
					Height: 844359,
					Type:   "transfer",
					Data: map[string]any{
						"amount":          "106033utia",
						"authz_msg_index": "2",
						"recipient":       "celestia1ghwz05j5s52nyvzau08eg9rqvkzgq72r92k56d",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "2",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":          "106033utia",
						"authz_msg_index": "2",
						"delegator":       "celestia1ghwz05j5s52nyvzau08eg9rqvkzgq72r92k56d",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "106031utia",
						"authz_msg_index": "2",
						"spender":         "celestia1ghwz05j5s52nyvzau08eg9rqvkzgq72r92k56d",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "106031utia",
						"authz_msg_index": "2",
						"receiver":        "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844359,
					Type:   "delegate",
					Data: map[string]any{
						"amount":          "106031utia",
						"authz_msg_index": "2",
						"new_shares":      "106031.000000000000000000",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "2",
						"module":          "staking",
						"sender":          "celestia1ghwz05j5s52nyvzau08eg9rqvkzgq72r92k56d",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "114972utia",
						"authz_msg_index": "3",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "114972utia",
						"authz_msg_index": "3",
						"receiver":        "celestia1vx8nc79y47y7kuez8m9z8hxzjcu9sy8jyymfmn",
					},
				}, {
					Height: 844359,
					Type:   "transfer",
					Data: map[string]any{
						"amount":          "114972utia",
						"authz_msg_index": "3",
						"recipient":       "celestia1vx8nc79y47y7kuez8m9z8hxzjcu9sy8jyymfmn",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "3",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":          "114972utia",
						"authz_msg_index": "3",
						"delegator":       "celestia1vx8nc79y47y7kuez8m9z8hxzjcu9sy8jyymfmn",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "114961utia",
						"authz_msg_index": "3",
						"spender":         "celestia1vx8nc79y47y7kuez8m9z8hxzjcu9sy8jyymfmn",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "114961utia",
						"authz_msg_index": "3",
						"receiver":        "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844359,
					Type:   "delegate",
					Data: map[string]any{
						"amount":          "114961utia",
						"authz_msg_index": "3",
						"new_shares":      "114961.000000000000000000",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "3",
						"module":          "staking",
						"sender":          "celestia1vx8nc79y47y7kuez8m9z8hxzjcu9sy8jyymfmn",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgExec,
				Height: 844359,
				InternalMsgs: []string{
					"/cosmos.staking.v1beta1.MsgDelegate",
					"/cosmos.staking.v1beta1.MsgDelegate",
					"/cosmos.staking.v1beta1.MsgDelegate",
					"/cosmos.staking.v1beta1.MsgDelegate",
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "multiple undelegations",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 595997,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.authz.v1beta1.MsgExec",
					},
				}, {
					Height: 595997,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":          "0stake",
						"authz_msg_index": "0",
						"delegator":       "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
						"validator":       "celestiavaloper1qq9dduljfua3uztkvgh6ytcahpzxs5qvkq97l3",
					},
				}, {
					Height: 595997,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "1000utia",
						"authz_msg_index": "0",
						"spender":         "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "1000utia",
						"authz_msg_index": "0",
						"receiver":        "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 595997,
					Type:   "transfer",
					Data: map[string]any{
						"amount":          "1000utia",
						"authz_msg_index": "0",
						"recipient":       "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"sender":          "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "0",
						"sender":          "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "unbond",
					Data: map[string]any{
						"amount":          "1000utia",
						"authz_msg_index": "0",
						"completion_time": "2023-12-18T18:38:51Z",
						"validator":       "celestiavaloper1qq9dduljfua3uztkvgh6ytcahpzxs5qvkq97l3",
					},
				}, {
					Height: 595997,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "0",
						"module":          "staking",
						"sender":          "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
					},
				}, {
					Height: 595997,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":          "0stake",
						"authz_msg_index": "1",
						"delegator":       "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
						"validator":       "celestiavaloper1qycj0ymu9fqvwgyw4xz93p3n4a83jjk7sm2wzh",
					},
				}, {
					Height: 595997,
					Type:   "unbond",
					Data: map[string]any{
						"amount":          "1000utia",
						"authz_msg_index": "1",
						"completion_time": "2023-12-18T18:38:51Z",
						"validator":       "celestiavaloper1qycj0ymu9fqvwgyw4xz93p3n4a83jjk7sm2wzh",
					},
				}, {
					Height: 595997,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "1",
						"module":          "staking",
						"sender":          "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
					},
				}, {
					Height: 595997,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":          "0stake",
						"authz_msg_index": "2",
						"delegator":       "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
						"validator":       "celestiavaloper1qyuwqj0cxe6hlzjru587nygwwmgh03ha9ve9ac",
					},
				}, {
					Height: 595997,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "1000utia",
						"authz_msg_index": "2",
						"spender":         "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "1000utia",
						"authz_msg_index": "2",
						"receiver":        "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 595997,
					Type:   "transfer",
					Data: map[string]any{
						"amount":          "1000utia",
						"authz_msg_index": "2",
						"recipient":       "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"sender":          "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "2",
						"sender":          "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "unbond",
					Data: map[string]any{
						"amount":          "1000utia",
						"authz_msg_index": "2",
						"completion_time": "2023-12-18T18:38:51Z",
						"validator":       "celestiavaloper1qyuwqj0cxe6hlzjru587nygwwmgh03ha9ve9ac",
					},
				}, {
					Height: 595997,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "2",
						"module":          "staking",
						"sender":          "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgExec,
				Height: 595997,
				InternalMsgs: []string{
					"/cosmos.staking.v1beta1.MsgUndelegate",
					"/cosmos.staking.v1beta1.MsgUndelegate",
					"/cosmos.staking.v1beta1.MsgUndelegate",
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "MsgWithdrawDelegatorReward",
			events: []storage.Event{
				{
					Height: 977944,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "57268utia",
						"spender": "celestia1gsvxuzts55h70c4338cmypzphv7l0exwc33jmp",
					},
				}, {
					Height: 977944,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "57268utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 977944,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "57268utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia1gsvxuzts55h70c4338cmypzphv7l0exwc33jmp",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1gsvxuzts55h70c4338cmypzphv7l0exwc33jmp",
					},
				}, {
					Height: 977944,
					Type:   "tx",
					Data: map[string]any{
						"fee":       "57268utia",
						"fee_payer": "celestia1gsvxuzts55h70c4338cmypzphv7l0exwc33jmp",
					},
				}, {
					Height: 977944,
					Type:   "tx",
					Data: map[string]any{
						"acc_seq": "celestia1gsvxuzts55h70c4338cmypzphv7l0exwc33jmp/4",
					},
				}, {
					Height: 977944,
					Type:   "tx",
					Data: map[string]any{
						"signature": "KwNCgyd6IKCxDzqwh2s8uwS88HbXBNt+gBrx3MJPsI91xY9UKZXJhPIYifMdUS/Fdo0R6EL91VRWc7RXMGK0yA==",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.authz.v1beta1.MsgExec",
					},
				}, {
					Height: 977944,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "77utia",
						"authz_msg_index": "0",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "77utia",
						"authz_msg_index": "0",
						"receiver":        "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					},
				}, {
					Height: 977944,
					Type:   "transfer",
					Data: map[string]any{
						"amount":          "77utia",
						"authz_msg_index": "0",
						"recipient":       "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "0",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":          "77utia",
						"authz_msg_index": "0",
						"delegator":       "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
						"validator":       "celestiavaloper1pnzrk7yzx0nr9xrcjyswj7ram4qxlrfz28xvn6",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "0",
						"module":          "distribution",
						"sender":          "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					},
				}, {
					Height: 977944,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":          "127utia",
						"authz_msg_index": "1",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":          "127utia",
						"authz_msg_index": "1",
						"receiver":        "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					},
				}, {
					Height: 977944,
					Type:   "transfer",
					Data: map[string]any{
						"amount":          "127utia",
						"authz_msg_index": "1",
						"recipient":       "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "1",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":          "127utia",
						"authz_msg_index": "1",
						"delegator":       "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
						"validator":       "celestiavaloper1nwu3ugynh8m6r7aphv0uxnca84t7gnruvyye9c",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]any{
						"authz_msg_index": "1",
						"module":          "distribution",
						"sender":          "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgExec,
				Height: 977944,
				InternalMsgs: []string{
					"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
					"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
				},
			},
			idx: testsuite.Ptr(7),
		}, {
			name: "unknown message",
			events: []storage.Event{
				{
					Height: 45631,
					Type:   "use_feegrant",
					Data: map[string]any{
						"grantee": "celestia1js8h76lxsl92qpqmsgd04u52aaqp82pr9n4p8f",
						"granter": "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "update_feegrant",
					Data: map[string]any{
						"grantee": "celestia1js8h76lxsl92qpqmsgd04u52aaqp82pr9n4p8f",
						"granter": "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "20443utia",
						"spender": "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "20443utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "20443utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]any{
						"fee":       "20443utia",
						"fee_payer": "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]any{
						"acc_seq": "celestia1js8h76lxsl92qpqmsgd04u52aaqp82pr9n4p8f/0",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]any{
						"signature": "wLj4V9p4OMCoP3OZLQU6RohalGFSHXQWXZ/7pMN1/FUx9n3YWh3eWwXb1PAYgkajLM3kGehn3mR770lufT7f+w==",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.authz.v1beta1.MsgExec",
					},
				}, {
					Height: 45631,
					Type:   "cosmos.authz.v1beta1.EventGrant",
					Data: map[string]any{
						"authz_msg_index": "0",
						"grantee":         "celestia10eykchznjdn8jdlwaj5v9wvlmdsp6kxx8ddhq6",
						"granter":         "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
						"msg_type_url":    "/cosmos.gov.v1beta1.MsgVote",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgGrant,
				Height: 45631,
				InternalMsgs: []string{
					"/cosmos.authz.v1beta1.MsgGrant",
				},
			},
			idx: testsuite.Ptr(9),
		}, {
			name: "signal version",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "210000utia",
						"spender": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "210000utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "210000utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]any{
						"fee":       "210000utia",
						"fee_payer": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]any{
						"acc_seq": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw/0",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]any{
						"signature": "pRnJae6JayIR/V4E1fwTpMC8myY3jldBM6YNWtgcuRIqLib2O6Tu06Ki3Yx3QEyjKiZSA1hH0jPRa3G/+/1Lcw==",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]any{
						"action":    "/cosmos.authz.v1beta1.MsgExec",
						"module":    "authz",
						"msg_index": "0",
						"sender":    "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					},
				}, {
					Height: 45631,
					Type:   "signal_version",
					Data: map[string]any{
						"action":            "/celestia.signal.v1.Msg/SignalVersion",
						"authz_msg_index":   "0",
						"msg_index":         "0",
						"validator_address": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgExec,
				Height: 45631,
				Data: map[string]any{
					"Grantee": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					"Msgs": []any{
						map[string]any{
							"ValidatorAddress": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
							"Version":          5,
						},
					},
				},
				InternalMsgs: []string{
					"/celestia.signal.v1.Msg/SignalVersion",
				},
			},
			idx: testsuite.Ptr(7),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleExec(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
		})
	}
}
