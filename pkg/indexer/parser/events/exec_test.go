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

func Test_handleExec(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
		idx    int
	}{
		{
			name: "multiple delegations",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 844359,
					Type:   "message",
					Data: map[string]string{
						"action": "/cosmos.authz.v1beta1.MsgExec",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "101775utia",
						"authz_msg_index": "0",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "101775utia",
						"authz_msg_index": "0",
						"receiver":        "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
					},
				}, {
					Height: 844359,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "101775utia",
						"authz_msg_index": "0",
						"recipient":       "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "0",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "101775utia",
						"authz_msg_index": "0",
						"delegator":       "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "101774utia",
						"authz_msg_index": "0",
						"spender":         "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "101774utia",
						"authz_msg_index": "0",
						"receiver":        "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844359,
					Type:   "delegate",
					Data: map[string]string{
						"amount":          "101774utia",
						"authz_msg_index": "0",
						"new_shares":      "101774.000000000000000000",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "0",
						"module":          "staking",
						"sender":          "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "191045utia",
						"authz_msg_index": "1",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "191045utia",
						"authz_msg_index": "1",
						"receiver":        "celestia1gfvu3xpgze2jy20cy4lcfeq2qj3rww0a8cwuap",
					},
				}, {
					Height: 844359,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "191045utia",
						"authz_msg_index": "1",
						"recipient":       "celestia1gfvu3xpgze2jy20cy4lcfeq2qj3rww0a8cwuap",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "1",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "191045utia",
						"authz_msg_index": "1",
						"delegator":       "celestia1gfvu3xpgze2jy20cy4lcfeq2qj3rww0a8cwuap",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "191018utia",
						"authz_msg_index": "1",
						"spender":         "celestia1gfvu3xpgze2jy20cy4lcfeq2qj3rww0a8cwuap",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "191018utia",
						"authz_msg_index": "1",
						"receiver":        "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844359,
					Type:   "delegate",
					Data: map[string]string{
						"amount":          "191018utia",
						"authz_msg_index": "1",
						"new_shares":      "191018.000000000000000000",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "1",
						"module":          "staking",
						"sender":          "celestia1gfvu3xpgze2jy20cy4lcfeq2qj3rww0a8cwuap",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "106033utia",
						"authz_msg_index": "2",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "106033utia",
						"authz_msg_index": "2",
						"receiver":        "celestia1ghwz05j5s52nyvzau08eg9rqvkzgq72r92k56d",
					},
				}, {
					Height: 844359,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "106033utia",
						"authz_msg_index": "2",
						"recipient":       "celestia1ghwz05j5s52nyvzau08eg9rqvkzgq72r92k56d",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "2",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "106033utia",
						"authz_msg_index": "2",
						"delegator":       "celestia1ghwz05j5s52nyvzau08eg9rqvkzgq72r92k56d",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "106031utia",
						"authz_msg_index": "2",
						"spender":         "celestia1ghwz05j5s52nyvzau08eg9rqvkzgq72r92k56d",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "106031utia",
						"authz_msg_index": "2",
						"receiver":        "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844359,
					Type:   "delegate",
					Data: map[string]string{
						"amount":          "106031utia",
						"authz_msg_index": "2",
						"new_shares":      "106031.000000000000000000",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "2",
						"module":          "staking",
						"sender":          "celestia1ghwz05j5s52nyvzau08eg9rqvkzgq72r92k56d",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "114972utia",
						"authz_msg_index": "3",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "114972utia",
						"authz_msg_index": "3",
						"receiver":        "celestia1vx8nc79y47y7kuez8m9z8hxzjcu9sy8jyymfmn",
					},
				}, {
					Height: 844359,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "114972utia",
						"authz_msg_index": "3",
						"recipient":       "celestia1vx8nc79y47y7kuez8m9z8hxzjcu9sy8jyymfmn",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "3",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 844359,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "114972utia",
						"authz_msg_index": "3",
						"delegator":       "celestia1vx8nc79y47y7kuez8m9z8hxzjcu9sy8jyymfmn",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "114961utia",
						"authz_msg_index": "3",
						"spender":         "celestia1vx8nc79y47y7kuez8m9z8hxzjcu9sy8jyymfmn",
					},
				}, {
					Height: 844359,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "114961utia",
						"authz_msg_index": "3",
						"receiver":        "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 844359,
					Type:   "delegate",
					Data: map[string]string{
						"amount":          "114961utia",
						"authz_msg_index": "3",
						"new_shares":      "114961.000000000000000000",
						"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
					},
				}, {
					Height: 844359,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "3",
						"module":          "staking",
						"sender":          "celestia1vx8nc79y47y7kuez8m9z8hxzjcu9sy8jyymfmn",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgExec,
				Height: 844359,
				Data: map[string]any{
					"Grantee": "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					"Msgs": []any{
						map[string]any{}, map[string]any{}, map[string]any{}, map[string]any{},
					},
				},
				InternalMsgs: []string{
					"/cosmos.staking.v1beta1.MsgDelegate",
					"/cosmos.staking.v1beta1.MsgDelegate",
					"/cosmos.staking.v1beta1.MsgDelegate",
					"/cosmos.staking.v1beta1.MsgDelegate",
				},
			},
			idx: 0,
		}, {
			name: "multiple undelegations",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 595997,
					Type:   "message",
					Data: map[string]string{
						"action": "/cosmos.authz.v1beta1.MsgExec",
					},
				}, {
					Height: 595997,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "0stake",
						"authz_msg_index": "0",
						"delegator":       "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
						"validator":       "celestiavaloper1qq9dduljfua3uztkvgh6ytcahpzxs5qvkq97l3",
					},
				}, {
					Height: 595997,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "1000utia",
						"authz_msg_index": "0",
						"spender":         "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "1000utia",
						"authz_msg_index": "0",
						"receiver":        "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 595997,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "1000utia",
						"authz_msg_index": "0",
						"recipient":       "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"sender":          "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "0",
						"sender":          "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "unbond",
					Data: map[string]string{
						"amount":          "1000utia",
						"authz_msg_index": "0",
						"completion_time": "2023-12-18T18:38:51Z",
						"validator":       "celestiavaloper1qq9dduljfua3uztkvgh6ytcahpzxs5qvkq97l3",
					},
				}, {
					Height: 595997,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "0",
						"module":          "staking",
						"sender":          "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
					},
				}, {
					Height: 595997,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "0stake",
						"authz_msg_index": "1",
						"delegator":       "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
						"validator":       "celestiavaloper1qycj0ymu9fqvwgyw4xz93p3n4a83jjk7sm2wzh",
					},
				}, {
					Height: 595997,
					Type:   "unbond",
					Data: map[string]string{
						"amount":          "1000utia",
						"authz_msg_index": "1",
						"completion_time": "2023-12-18T18:38:51Z",
						"validator":       "celestiavaloper1qycj0ymu9fqvwgyw4xz93p3n4a83jjk7sm2wzh",
					},
				}, {
					Height: 595997,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "1",
						"module":          "staking",
						"sender":          "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
					},
				}, {
					Height: 595997,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "0stake",
						"authz_msg_index": "2",
						"delegator":       "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
						"validator":       "celestiavaloper1qyuwqj0cxe6hlzjru587nygwwmgh03ha9ve9ac",
					},
				}, {
					Height: 595997,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "1000utia",
						"authz_msg_index": "2",
						"spender":         "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "1000utia",
						"authz_msg_index": "2",
						"receiver":        "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 595997,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "1000utia",
						"authz_msg_index": "2",
						"recipient":       "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"sender":          "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "2",
						"sender":          "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 595997,
					Type:   "unbond",
					Data: map[string]string{
						"amount":          "1000utia",
						"authz_msg_index": "2",
						"completion_time": "2023-12-18T18:38:51Z",
						"validator":       "celestiavaloper1qyuwqj0cxe6hlzjru587nygwwmgh03ha9ve9ac",
					},
				}, {
					Height: 595997,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "2",
						"module":          "staking",
						"sender":          "celestia1um8q93lngf6hfvslqn2nph77f2xyeklppjlxlc",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgExec,
				Height: 595997,
				Data: map[string]any{
					"Grantee": "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					"Msgs": []any{
						map[string]any{}, map[string]any{}, map[string]any{},
					},
				},
				InternalMsgs: []string{
					"/cosmos.staking.v1beta1.MsgUndelegate",
					"/cosmos.staking.v1beta1.MsgUndelegate",
					"/cosmos.staking.v1beta1.MsgUndelegate",
				},
			},
			idx: 0,
		}, {
			name: "MsgWithdrawDelegatorReward",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 977944,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":  "57268utia",
						"spender": "celestia1gsvxuzts55h70c4338cmypzphv7l0exwc33jmp",
					},
				}, {
					Height: 977944,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":   "57268utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 977944,
					Type:   "transfer",
					Data: map[string]string{
						"amount":    "57268utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia1gsvxuzts55h70c4338cmypzphv7l0exwc33jmp",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]string{
						"sender": "celestia1gsvxuzts55h70c4338cmypzphv7l0exwc33jmp",
					},
				}, {
					Height: 977944,
					Type:   "tx",
					Data: map[string]string{
						"fee":       "57268utia",
						"fee_payer": "celestia1gsvxuzts55h70c4338cmypzphv7l0exwc33jmp",
					},
				}, {
					Height: 977944,
					Type:   "tx",
					Data: map[string]string{
						"acc_seq": "celestia1gsvxuzts55h70c4338cmypzphv7l0exwc33jmp/4",
					},
				}, {
					Height: 977944,
					Type:   "tx",
					Data: map[string]string{
						"signature": "KwNCgyd6IKCxDzqwh2s8uwS88HbXBNt+gBrx3MJPsI91xY9UKZXJhPIYifMdUS/Fdo0R6EL91VRWc7RXMGK0yA==",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]string{
						"action": "/cosmos.authz.v1beta1.MsgExec",
					},
				}, {
					Height: 977944,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "77utia",
						"authz_msg_index": "0",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "77utia",
						"authz_msg_index": "0",
						"receiver":        "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					},
				}, {
					Height: 977944,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "77utia",
						"authz_msg_index": "0",
						"recipient":       "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "0",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "77utia",
						"authz_msg_index": "0",
						"delegator":       "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
						"validator":       "celestiavaloper1pnzrk7yzx0nr9xrcjyswj7ram4qxlrfz28xvn6",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "0",
						"module":          "distribution",
						"sender":          "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					},
				}, {
					Height: 977944,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "127utia",
						"authz_msg_index": "1",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "127utia",
						"authz_msg_index": "1",
						"receiver":        "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					},
				}, {
					Height: 977944,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "127utia",
						"authz_msg_index": "1",
						"recipient":       "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "1",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 977944,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "127utia",
						"authz_msg_index": "1",
						"delegator":       "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
						"validator":       "celestiavaloper1nwu3ugynh8m6r7aphv0uxnca84t7gnruvyye9c",
					},
				}, {
					Height: 977944,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "1",
						"module":          "distribution",
						"sender":          "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgExec,
				Height: 977944,
				Data: map[string]any{
					"Grantee": "celestia1zq2atge5df93w0l6xhm87r8uspjva632aqz4fe",
					"Msgs": []any{
						map[string]any{}, map[string]any{},
					},
				},
				InternalMsgs: []string{
					"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
					"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
				},
			},
			idx: 7,
		}, {
			name: "unknown message",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 45631,
					Type:   "use_feegrant",
					Data: map[string]string{
						"grantee": "celestia1js8h76lxsl92qpqmsgd04u52aaqp82pr9n4p8f",
						"granter": "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "update_feegrant",
					Data: map[string]string{
						"grantee": "celestia1js8h76lxsl92qpqmsgd04u52aaqp82pr9n4p8f",
						"granter": "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":  "20443utia",
						"spender": "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":   "20443utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]string{
						"amount":    "20443utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"sender": "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"fee":       "20443utia",
						"fee_payer": "celestia1rcm7tth05klgkqpucdhm5hexnk49dfda5qnwts",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"acc_seq": "celestia1js8h76lxsl92qpqmsgd04u52aaqp82pr9n4p8f/0",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"signature": "wLj4V9p4OMCoP3OZLQU6RohalGFSHXQWXZ/7pMN1/FUx9n3YWh3eWwXb1PAYgkajLM3kGehn3mR770lufT7f+w==",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"action": "/cosmos.authz.v1beta1.MsgExec",
					},
				}, {
					Height: 45631,
					Type:   "cosmos.authz.v1beta1.EventGrant",
					Data: map[string]string{
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
				Data: map[string]any{
					"Grantee": "celestia10eykchznjdn8jdlwaj5v9wvlmdsp6kxx8ddhq6",
					"Msgs": []any{
						map[string]any{},
					},
				},
				InternalMsgs: []string{
					"/cosmos.authz.v1beta1.MsgGrant",
				},
			},
			idx: 9,
		}, {
			name: "signal version",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":  "210000utia",
						"spender": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":   "210000utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]string{
						"amount":    "210000utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"sender": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"fee":       "210000utia",
						"fee_payer": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"acc_seq": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw/0",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"signature": "pRnJae6JayIR/V4E1fwTpMC8myY3jldBM6YNWtgcuRIqLib2O6Tu06Ki3Yx3QEyjKiZSA1hH0jPRa3G/+/1Lcw==",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"action":    "/cosmos.authz.v1beta1.MsgExec",
						"module":    "authz",
						"msg_index": "0",
						"sender":    "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
					},
				}, {
					Height: 45631,
					Type:   "signal_version",
					Data: map[string]string{
						"action":            "/celestia.signal.v1.MsgSignalVersion",
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
			idx: 7,
		}, {
			name: "vote",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":  "2259utia",
						"spender": "celestia1cnkaqutgx24kyrge8mpp5qq0pw6mll93572anj",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":   "2259utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]string{
						"amount":    "2259utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia1cnkaqutgx24kyrge8mpp5qq0pw6mll93572anj",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"sender": "celestia1cnkaqutgx24kyrge8mpp5qq0pw6mll93572anj",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"fee":       "2259utia",
						"fee_payer": "celestia1cnkaqutgx24kyrge8mpp5qq0pw6mll93572anj",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"acc_seq": "celestia1cnkaqutgx24kyrge8mpp5qq0pw6mll93572anj/0",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"signature": "/FrbyKniUo9iKrZnIrb3czbur2Lewro/7l2aeDEHTMlG3QGdesQKHi2ILs54MOUWNTDsAdFMvD/xWMOQ4VNTPg==",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"action": "/cosmos.authz.v1beta1.MsgExec",
					},
				}, {
					Height: 45631,
					Type:   "proposal_vote",
					Data: map[string]string{
						"authz_msg_index": "0",
						"option":          "option:VOTE_OPTION_ABSTAIN weight:\"1.000000000000000000\"",
						"proposal_id":     "2",
						"voter":           "celestia138jl42zlxue4wpvnugcdqhxjmyd2vpt69gjdfk",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "0",
						"module":          "governance",
						"sender":          "celestia138jl42zlxue4wpvnugcdqhxjmyd2vpt69gjdfk",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgExec,
				Height: 45631,
				Data: map[string]any{
					"Grantee": "celestia1cnkaqutgx24kyrge8mpp5qq0pw6mll93572anj",
					"Msgs": []any{
						map[string]any{
							"Option":     2,
							"ProposalId": 2,
							"Voter":      "celestia138jl42zlxue4wpvnugcdqhxjmyd2vpt69gjdfk",
						},
					},
				},
				InternalMsgs: []string{
					"/cosmos.gov.v1.MsgVoteWeighted",
				},
			},
			idx: 7,
		}, {
			name: "withdraw commissions and delegator rewards",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":  "2000utia",
						"spender": "celestia1eqj8gju6sldleq76auawrl58nwqcl6dyh9xdry",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":   "2000utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]string{
						"amount":    "2000utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia1eqj8gju6sldleq76auawrl58nwqcl6dyh9xdry",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"sender": "celestia1eqj8gju6sldleq76auawrl58nwqcl6dyh9xdry",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"fee":       "2000utia",
						"fee_payer": "celestia1eqj8gju6sldleq76auawrl58nwqcl6dyh9xdry",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"acc_seq": "celestia1eqj8gju6sldleq76auawrl58nwqcl6dyh9xdry/8",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"signature": "I5CqepCwRdt5drYYFZGBaXJF4Idgr/CbRURa8zSOZD9WIfI61AlmHpfP1oiZQqwBVJQTUi1TT3NSbADUpZnqrA==",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"action":    "/cosmos.authz.v1beta1.MsgExec",
						"module":    "authz",
						"msg_index": "0",
						"sender":    "celestia1eqj8gju6sldleq76auawrl58nwqcl6dyh9xdry",
					},
				}, {
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "41419902utia",
						"authz_msg_index": "0",
						"msg_index":       "0",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "41419902utia",
						"authz_msg_index": "0",
						"msg_index":       "0",
						"receiver":        "celestia1fwyzhlyngfguxg0pkm34mqd0ycw8ua5ftnkqcm",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "41419902utia",
						"authz_msg_index": "0",
						"msg_index":       "0",
						"recipient":       "celestia1fwyzhlyngfguxg0pkm34mqd0ycw8ua5ftnkqcm",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "0",
						"msg_index":       "0",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "41419902utia",
						"authz_msg_index": "0",
						"delegator":       "celestia1gl0rg3g0pkcpr8umj2hvlhha06ecjd65p58r5j",
						"msg_index":       "0",
						"validator":       "celestiavaloper1gl0rg3g0pkcpr8umj2hvlhha06ecjd65yt96z5",
					},
				}, {
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "8323576892utia",
						"authz_msg_index": "1",
						"msg_index":       "0",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "8323576892utia",
						"authz_msg_index": "1",
						"msg_index":       "0",
						"receiver":        "celestia1fwyzhlyngfguxg0pkm34mqd0ycw8ua5ftnkqcm",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "8323576892utia",
						"authz_msg_index": "1",
						"msg_index":       "0",
						"recipient":       "celestia1fwyzhlyngfguxg0pkm34mqd0ycw8ua5ftnkqcm",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "1",
						"msg_index":       "0",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "withdraw_commission",
					Data: map[string]string{
						"amount":          "8323576892utia",
						"authz_msg_index": "1",
						"msg_index":       "0",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgExec,
				Height: 45631,
				Data: map[string]any{
					"Grantee": "celestia1eqj8gju6sldleq76auawrl58nwqcl6dyh9xdry",
					"Msgs": []any{
						map[string]any{
							"DelegatorAddress": "celestia1gl0rg3g0pkcpr8umj2hvlhha06ecjd65p58r5j",
							"ValidatorAddress": "celestiavaloper1gl0rg3g0pkcpr8umj2hvlhha06ecjd65yt96z5",
						},
						map[string]any{
							"ValidatorAddress": "celestiavaloper1gl0rg3g0pkcpr8umj2hvlhha06ecjd65yt96z5",
						},
					},
				},
				InternalMsgs: []string{
					"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
					"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
				},
			},
			idx: 7,
		}, {
			name: "withdraw commissions, delegator rewards and delegate",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":  "46797utia",
						"spender": "celestia1xvj465q8yu59uy78wtqs9fv6r9urlsa89q4s8f",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":   "46797utia",
						"receiver": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]string{
						"amount":    "46797utia",
						"recipient": "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s",
						"sender":    "celestia1xvj465q8yu59uy78wtqs9fv6r9urlsa89q4s8f",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"sender": "celestia1xvj465q8yu59uy78wtqs9fv6r9urlsa89q4s8f",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"fee":       "46797utia",
						"fee_payer": "celestia1xvj465q8yu59uy78wtqs9fv6r9urlsa89q4s8f",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"acc_seq": "celestia1xvj465q8yu59uy78wtqs9fv6r9urlsa89q4s8f/0",
					},
				}, {
					Height: 45631,
					Type:   "tx",
					Data: map[string]string{
						"signature": "I5CqepCwRdt5drYYFZGBaXJF4Idgr/CbRURa8zSOZD9WIfI61AlmHpfP1oiZQqwBVJQTUi1TT3NSbADUpZnqrA==",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"action": "/cosmos.authz.v1beta1.MsgExec",
					},
				}, {
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "164617964utia",
						"authz_msg_index": "0",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "164617964utia",
						"authz_msg_index": "0",
						"receiver":        "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "164617964utia",
						"authz_msg_index": "0",
						"recipient":       "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "0",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "164617964utia",
						"authz_msg_index": "0",
						"delegator":       "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
						"validator":       "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "0",
						"module":          "distribution",
						"sender":          "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
					},
				}, {
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "4174387254utia",
						"authz_msg_index": "1",
						"spender":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "4174387254utia",
						"authz_msg_index": "1",
						"receiver":        "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
					},
				}, {
					Height: 45631,
					Type:   "transfer",
					Data: map[string]string{
						"amount":          "4174387254utia",
						"authz_msg_index": "1",
						"recipient":       "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "1",
						"sender":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 45631,
					Type:   "withdraw_commission",
					Data: map[string]string{
						"amount":          "4174387254utia",
						"authz_msg_index": "1",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "1",
						"module":          "distribution",
						"sender":          "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
					},
				}, {
					Height: 45631,
					Type:   "withdraw_rewards",
					Data: map[string]string{
						"amount":          "0stake",
						"authz_msg_index": "2",
						"delegator":       "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
						"validator":       "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
					},
				}, {
					Height: 45631,
					Type:   "coin_spent",
					Data: map[string]string{
						"amount":          "4338363027utia",
						"authz_msg_index": "2",
						"spender":         "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
					},
				}, {
					Height: 45631,
					Type:   "coin_received",
					Data: map[string]string{
						"amount":          "4338363027utia",
						"authz_msg_index": "2",
						"receiver":        "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 45631,
					Type:   "delegate",
					Data: map[string]string{
						"amount":          "4338363027utia",
						"authz_msg_index": "2",
						"new_shares":      "4338363027.000000000000000000",
						"validator":       "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
					},
				}, {
					Height: 45631,
					Type:   "message",
					Data: map[string]string{
						"authz_msg_index": "2",
						"module":          "staking",
						"sender":          "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgExec,
				Height: 45631,
				Data: map[string]any{
					"Grantee": "celestia1eqj8gju6sldleq76auawrl58nwqcl6dyh9xdry",
					"Msgs": []any{
						map[string]any{
							"DelegatorAddress": "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
							"ValidatorAddress": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
						},
						map[string]any{
							"ValidatorAddress": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
						},
						map[string]any{
							"Amount": map[string]any{
								"Amount": "4338363027",
								"Denom":  "utia",
							},
							"DelegatorAddress": "celestia15urq2dtp9qce4fyc85m6upwm9xul3049d30meq",
							"ValidatorAddress": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
						},
					},
				},
				InternalMsgs: []string{
					"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
					"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
					"/cosmos.staking.v1beta1.MsgDelegate",
				},
			},
			idx: 7,
		},
	}
	for _, tt := range tests {
		tt.ctx.Block = &storage.Block{
			Time:   time.Now(),
			Height: tt.msg.Height,
		}
		t.Run(tt.name, func(t *testing.T) {
			c := NewCursor(tt.events)
			c.Skip(tt.idx)
			err := handleExec(tt.ctx, c, tt.msg)
			require.NoError(t, err)
		})
	}
}

func Test_handleExec_TryUpgrade(t *testing.T) {
	const granter = "celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8"

	ctx := context.NewContext()
	ctx.Block = &storage.Block{Time: time.Now(), Height: 45631}

	events := []storage.Event{
		{
			Height: 45631,
			Type:   "message",
			Data: map[string]string{
				"action":    "/cosmos.authz.v1beta1.MsgExec",
				"module":    "authz",
				"msg_index": "0",
				"sender":    "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
			},
		}, {
			Height: 45631,
			Type:   "signal_try_upgrade",
			Data: map[string]string{
				"action":          "/celestia.signal.v1.MsgTryUpgrade",
				"authz_msg_index": "0",
				"msg_index":       "0",
				"signer":          granter,
			},
		}, {
			Height: 45631,
			Type:   "message",
			Data: map[string]string{
				"action":    "/cosmos.bank.v1beta1.MsgSend",
				"msg_index": "1",
			},
		},
	}
	msg := &storage.Message{
		Id:     10,
		TxId:   5,
		Type:   types.MsgExec,
		Height: 45631,
		Time:   ctx.Block.Time,
		Data: map[string]any{
			"Grantee": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
			"Msgs": []any{
				map[string]any{"Signer": granter},
			},
		},
		InternalMsgs: []string{"/celestia.signal.v1.MsgTryUpgrade"},
	}

	c := NewCursor(events)
	require.NoError(t, handleExec(ctx, c, msg))

	require.NotNil(t, ctx.TryUpgrade)
	require.EqualValues(t, 45631, ctx.TryUpgrade.Height)
	require.EqualValues(t, 45631, ctx.TryUpgrade.EndHeight)
	require.EqualValues(t, 5, ctx.TryUpgrade.TxId)
	require.EqualValues(t, 10, ctx.TryUpgrade.MsgId)
	require.NotNil(t, ctx.TryUpgrade.Signer)
	require.Equal(t, granter, ctx.TryUpgrade.Signer.Address)

	// the granter is registered, so storage can resolve signer_id
	_, ok := ctx.Addresses.Get(ctx.TryUpgrade.Signer.String())
	require.True(t, ok)

	// the event of the next message is left for its handler
	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "1", next.Data["msg_index"])
}

// Inner messages have no rows of their own, so the forwarding points at the MsgExec.
func Test_handleExec_Forward(t *testing.T) {
	ctx := context.NewContext()
	ctx.Block = &storage.Block{Time: time.Now(), Height: 45631}

	events := []storage.Event{
		{
			Height: 45631,
			Type:   "message",
			Data: map[string]string{
				"action":    "/cosmos.authz.v1beta1.MsgExec",
				"module":    "authz",
				"msg_index": "0",
				"sender":    "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
			},
		}, {
			Height: 45631,
			Type:   types.EventTypeHyperlanewarpv1EventSendRemoteTransfer,
			Data: map[string]string{ //nolint:gosec
				"authz_msg_index":    "0",
				"msg_index":          "0",
				"destination_domain": "11155111",
				"recipient":          "\"0x000000000000000000000000d5e85e86fc692cedad6d6992f1f0ccf273e39913\"",
				"sender":             "\"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60\"",
				"token_id":           "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
			},
		}, {
			Height: 45631,
			Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
			Data: map[string]string{ //nolint:gosec
				"authz_msg_index": "0",
				"msg_index":       "0",
				"forward_addr":    "\"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60\"",
				"denom":           "\"utia\"",
				"amount":          "\"1000\"",
				"message_id":      "\"0xac8852bd411c0c88cdadfe9b2386b2bcd702f35479c25a4b2d2cc3fb49d095d4\"",
				"token_id":        "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
			},
		}, {
			Height: 45631,
			Type:   "message",
			Data: map[string]string{
				"action":    "/cosmos.bank.v1beta1.MsgSend",
				"msg_index": "1",
			},
		},
	}
	msg := &storage.Message{
		Id:     10,
		TxId:   5,
		Type:   types.MsgExec,
		Height: 45631,
		Time:   ctx.Block.Time,
		Data: map[string]any{
			"Grantee": "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw",
			"Msgs": []any{
				map[string]any{"ForwardAddr": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"},
			},
		},
		InternalMsgs: []string{"/celestia.forwarding.v1.MsgForward"},
	}

	c := NewCursor(events)
	require.NoError(t, handleExec(ctx, c, msg))

	require.Len(t, ctx.Forwardings, 1)
	fwd := ctx.Forwardings[0]
	require.EqualValues(t, 5, fwd.TxId)
	require.EqualValues(t, 10, fwd.MsgId)
	require.EqualValues(t, 11155111, fwd.DestDomain)
	require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", fwd.Address.Address)

	// the event of the next message is left for its handler
	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "1", next.Data["msg_index"])
}

func forwardExecEvents(authzIdx, domain, addr, amount string) []storage.Event {
	return []storage.Event{
		{
			Height: 45631,
			Type:   types.EventTypeHyperlanewarpv1EventSendRemoteTransfer,
			Data: map[string]string{ //nolint:gosec
				"authz_msg_index":    authzIdx,
				"msg_index":          "0",
				"destination_domain": domain,
				"recipient":          "\"0x000000000000000000000000d5e85e86fc692cedad6d6992f1f0ccf273e39913\"",
				"sender":             "\"" + addr + "\"",
				"token_id":           "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
			},
		}, {
			Height: 45631,
			Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
			Data: map[string]string{ //nolint:gosec
				"authz_msg_index": authzIdx,
				"msg_index":       "0",
				"forward_addr":    "\"" + addr + "\"",
				"denom":           "\"utia\"",
				"amount":          "\"" + amount + "\"",
				"message_id":      "\"0xac8852bd411c0c88cdadfe9b2386b2bcd702f35479c25a4b2d2cc3fb49d095d4\"",
				"token_id":        "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
			},
		},
	}
}

func Test_handleExec_ForwardBoundary(t *testing.T) {
	execEvent := storage.Event{
		Height: 45631,
		Type:   "message",
		Data: map[string]string{
			"action":    "/cosmos.authz.v1beta1.MsgExec",
			"msg_index": "0",
		},
	}
	newMsg := func(t time.Time, inner ...string) *storage.Message {
		msgs := make([]any, len(inner))
		for i := range msgs {
			msgs[i] = map[string]any{}
		}
		return &storage.Message{
			Id:           10,
			TxId:         5,
			Type:         types.MsgExec,
			Height:       45631,
			Time:         t,
			Data:         map[string]any{"Msgs": msgs},
			InternalMsgs: inner,
		}
	}

	t.Run("two forwards", func(t *testing.T) {
		ctx := context.NewContext()
		ctx.Block = &storage.Block{Time: time.Now(), Height: 45631}

		events := make([]storage.Event, 0, 5)
		events = append(events, execEvent)
		events = append(events, forwardExecEvents("0", "1", "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", "100")...)
		events = append(events, forwardExecEvents("1", "2", "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw", "200")...)

		msg := newMsg(ctx.Block.Time, "/celestia.forwarding.v1.MsgForward", "/celestia.forwarding.v1.MsgForward")
		c := NewCursor(events)
		require.NoError(t, handleExec(ctx, c, msg))

		require.Len(t, ctx.Forwardings, 2)
		require.EqualValues(t, 1, ctx.Forwardings[0].DestDomain)
		require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", ctx.Forwardings[0].Address.Address)
		require.Equal(t, "100", ctx.Forwardings[0].Amount.String())
		require.EqualValues(t, 2, ctx.Forwardings[1].DestDomain)
		require.Equal(t, "celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw", ctx.Forwardings[1].Address.Address)
		require.Equal(t, "200", ctx.Forwardings[1].Amount.String())

		_, ok := c.Peek()
		require.False(t, ok)
	})

	t.Run("forward then delegate", func(t *testing.T) {
		ctx := context.NewContext()
		ctx.Block = &storage.Block{Time: time.Now(), Height: 45631}

		events := make([]storage.Event, 0, 5)
		events = append(events, execEvent)
		events = append(events, forwardExecEvents("0", "1", "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", "100")...)
		events = append(events, storage.Event{
			Height: 45631,
			Type:   "delegate",
			Data: map[string]string{
				"amount":          "101774utia",
				"authz_msg_index": "1",
				"new_shares":      "101774.000000000000000000",
				"validator":       "celestiavaloper1j2jq259d3rrc24876gwxg0ksp0lhd8gy49k6st",
			},
		}, storage.Event{
			Height: 45631,
			Type:   "message",
			Data: map[string]string{
				"authz_msg_index": "1",
				"module":          "staking",
				"sender":          "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
			},
		})

		msg := newMsg(ctx.Block.Time, "/celestia.forwarding.v1.MsgForward", "/cosmos.staking.v1beta1.MsgDelegate")
		c := NewCursor(events)
		require.NoError(t, handleExec(ctx, c, msg))

		require.Len(t, ctx.Forwardings, 1)
		require.Equal(t, 1, ctx.Delegations.Len())
	})
}

func Test_handleExec_InnerBoundary(t *testing.T) {
	const (
		valA = "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x"
		valB = "celestiavaloper1gl0rg3g0pkcpr8umj2hvlhha06ecjd65yt96z5"
	)
	execEvent := storage.Event{
		Height: 45631,
		Type:   "message",
		Data: map[string]string{
			"action":    "/cosmos.authz.v1beta1.MsgExec",
			"msg_index": "0",
		},
	}
	nextTxMsg := storage.Event{
		Height: 45631,
		Type:   "message",
		Data: map[string]string{
			"action":    "/cosmos.bank.v1beta1.MsgSend",
			"module":    "bank",
			"msg_index": "1",
		},
	}
	newExec := func(ts time.Time, inner []string, data []any) *storage.Message {
		return &storage.Message{
			Id:           10,
			TxId:         5,
			Type:         types.MsgExec,
			Height:       45631,
			Time:         ts,
			Data:         map[string]any{"Msgs": data},
			InternalMsgs: inner,
		}
	}
	newCtx := func() *context.Context {
		ctx := context.NewContext()
		ctx.Block = &storage.Block{Time: time.Now(), Height: 45631}
		return ctx
	}

	// the first commission has no amount, so the old scan took the second validator's withdrawal
	t.Run("commission does not take next message withdrawal", func(t *testing.T) {
		ctx := newCtx()
		events := []storage.Event{
			execEvent,
			{
				Height: 45631,
				Type:   "withdraw_commission",
				Data: map[string]string{
					"amount":          "",
					"authz_msg_index": "0",
					"msg_index":       "0",
				},
			}, {
				Height: 45631,
				Type:   "withdraw_commission",
				Data: map[string]string{
					"amount":          "100utia",
					"authz_msg_index": "1",
					"msg_index":       "0",
				},
			},
			nextTxMsg,
		}
		msg := newExec(ctx.Block.Time,
			[]string{msgWithdrawValidatorCommission, msgWithdrawValidatorCommission},
			[]any{
				map[string]any{"ValidatorAddress": valA},
				map[string]any{"ValidatorAddress": valB},
			},
		)
		c := NewCursor(events)
		require.NoError(t, handleExec(ctx, c, msg))

		require.Equal(t, 1, ctx.Validators.Len())
		_, ok := ctx.Validators.Get(valA)
		require.False(t, ok)
		val, ok := ctx.Validators.Get(valB)
		require.True(t, ok)
		require.Equal(t, "-100", val.Commissions.String())

		next, ok := c.Peek()
		require.True(t, ok)
		require.Equal(t, "1", next.Data["msg_index"])
	})

	// since SDK 0.50 an unjail inside MsgExec emits no events
	t.Run("unjail without events", func(t *testing.T) {
		ctx := newCtx()
		events := []storage.Event{execEvent, nextTxMsg}
		msg := newExec(ctx.Block.Time,
			[]string{"/cosmos.slashing.v1beta1.MsgUnjail"},
			[]any{map[string]any{"ValidatorAddr": valA}},
		)
		c := NewCursor(events)
		require.NoError(t, handleExec(ctx, c, msg))

		val, ok := ctx.Validators.Get(valA)
		require.True(t, ok)
		require.NotNil(t, val.Jailed)
		require.False(t, *val.Jailed)

		next, ok := c.Peek()
		require.True(t, ok)
		require.Equal(t, "1", next.Data["msg_index"])
	})

	t.Run("unjail between delegations", func(t *testing.T) {
		ctx := newCtx()
		delegate := func(idx, amount, validator string) storage.Event {
			return storage.Event{
				Height: 45631,
				Type:   "delegate",
				Data: map[string]string{
					"amount":          amount + "utia",
					"authz_msg_index": idx,
					"msg_index":       "0",
					"delegator":       "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
					"new_shares":      amount + ".000000000000000000",
					"validator":       validator,
				},
			}
		}
		events := []storage.Event{execEvent, delegate("0", "100", valA), delegate("2", "200", valB), nextTxMsg}
		msg := newExec(ctx.Block.Time,
			[]string{
				"/cosmos.staking.v1beta1.MsgDelegate",
				"/cosmos.slashing.v1beta1.MsgUnjail",
				"/cosmos.staking.v1beta1.MsgDelegate",
			},
			[]any{
				map[string]any{},
				map[string]any{"ValidatorAddr": valA},
				map[string]any{},
			},
		)
		c := NewCursor(events)
		require.NoError(t, handleExec(ctx, c, msg))

		require.Equal(t, 2, ctx.Delegations.Len())
		val, ok := ctx.Validators.Get(valA)
		require.True(t, ok)
		require.False(t, *val.Jailed)

		next, ok := c.Peek()
		require.True(t, ok)
		require.Equal(t, "1", next.Data["msg_index"])
	})

	// a used-up grant emits EventRevoke without authz_msg_index before the inner message events
	t.Run("revoke before inner message", func(t *testing.T) {
		ctx := newCtx()
		delegate := func(idx, validator string) storage.Event {
			return storage.Event{
				Height: 45631,
				Type:   "delegate",
				Data: map[string]string{
					"amount":          "100utia",
					"authz_msg_index": idx,
					"msg_index":       "0",
					"delegator":       "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
					"new_shares":      "100.000000000000000000",
					"validator":       validator,
				},
			}
		}
		revoke := storage.Event{
			Height: 45631,
			Type:   types.EventTypeCosmosauthzv1beta1EventRevoke,
			Data: map[string]string{
				"granter":      "\"celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp\"",
				"grantee":      "\"celestia10vj4f36sd4nr27c9meta7elxt87t9ww9vw8euw\"",
				"msg_type_url": "\"/cosmos.staking.v1beta1.MsgDelegate\"",
				"msg_index":    "0",
			},
		}
		events := []storage.Event{execEvent, delegate("0", valA), revoke, delegate("1", valB), nextTxMsg}
		msg := newExec(ctx.Block.Time,
			[]string{"/cosmos.staking.v1beta1.MsgDelegate", "/cosmos.staking.v1beta1.MsgDelegate"},
			[]any{map[string]any{}, map[string]any{}},
		)
		c := NewCursor(events)
		require.NoError(t, handleExec(ctx, c, msg))
		require.Equal(t, 2, ctx.Delegations.Len())

		next, ok := c.Peek()
		require.True(t, ok)
		require.Equal(t, "1", next.Data["msg_index"])
	})

	t.Run("handler leaves its own trailing events", func(t *testing.T) {
		ctx := newCtx()
		events := []storage.Event{
			execEvent,
			{
				Height: 45631,
				Type:   "proposal_vote",
				Data: map[string]string{
					"authz_msg_index": "0",
					"msg_index":       "0",
					"option":          `[{"option":1,"weight":"1.000000000000000000"}]`,
					"proposal_id":     "3",
					"voter":           "celestia1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4c64grjp",
				},
			}, {
				Height: 45631,
				Type:   "coin_received",
				Data: map[string]string{
					"amount":          "1utia",
					"authz_msg_index": "0",
					"msg_index":       "0",
				},
			},
			nextTxMsg,
		}
		msg := newExec(ctx.Block.Time,
			[]string{"/cosmos.gov.v1.MsgVote"},
			[]any{map[string]any{}},
		)
		c := NewCursor(events)
		require.NoError(t, handleExec(ctx, c, msg))
		require.Equal(t, 1, ctx.Votes.Len())

		next, ok := c.Peek()
		require.True(t, ok)
		require.Equal(t, "1", next.Data["msg_index"])
	})
}
