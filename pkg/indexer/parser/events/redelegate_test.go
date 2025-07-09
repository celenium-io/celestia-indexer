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

func Test_handleRedelegate(t *testing.T) {
	ts := time.Now()
	tests := []struct {
		name         string
		ctx          *context.Context
		events       []storage.Event
		msg          *storage.Message
		idx          *int
		redelegation *storage.Redelegation
	}{
		{
			name: "redelegate test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 841682,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgBeginRedelegate",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "384192utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":  "384192utia",
						"spender": "celestia1kj4m7nhhcr5f5jpcdaec5ymhepc3x22yvtjh7j",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "384192utia",
						"recipient": "celestia1kj4m7nhhcr5f5jpcdaec5ymhepc3x22yvtjh7j",
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
						"amount":    "384192utia",
						"delegator": "celestia1kj4m7nhhcr5f5jpcdaec5ymhepc3x22yvtjh7j",
						"validator": "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "redelegate",
					Data: map[string]any{
						"amount":                "20000000utia",
						"completion_time":       "2024-03-15T00:04:38Z",
						"destination_validator": "celestiavaloper107lwx458gy345ag2afx9a7e2kkl7x49y3433gj",
						"source_validator":      "celestiavaloper1uqj5ul7jtpskk9ste9mfv6jvh0y3w34vtpz3gw",
					},
				}, {
					Height: 841682,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1kj4m7nhhcr5f5jpcdaec5ymhepc3x22yvtjh7j",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgBeginRedelegate,
				Time:   ts,
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
			redelegation: &storage.Redelegation{
				Time:           ts,
				Height:         841682,
				Amount:         decimal.RequireFromString("20000000"),
				CompletionTime: time.Date(2024, 3, 15, 0, 4, 38, 0, time.UTC),
				Address: &storage.Address{
					Address:    "celestia1kj4m7nhhcr5f5jpcdaec5ymhepc3x22yvtjh7j",
					Height:     841682,
					LastHeight: 841682,
					Hash:       []byte{0xb4, 0xab, 0xbf, 0x4e, 0xf7, 0xc0, 0xe8, 0x9a, 0x48, 0x38, 0x6f, 0x73, 0x8a, 0x13, 0x77, 0xc8, 0x71, 0x13, 0x29, 0x44},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.Zero,
					},
				},
				Source: &storage.Validator{
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
					Stake:             decimal.RequireFromString("-20000000"),
				},
				Destination: &storage.Validator{
					Address:           "celestiavaloper107lwx458gy345ag2afx9a7e2kkl7x49y3433gj",
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
					Stake:             decimal.RequireFromString("20000000"),
				},
			},
		}, {
			name: "redelegate test 2",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 241,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgBeginRedelegate",
					},
				}, {
					Height: 241,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0stake",
						"delegator": "celestia1ze2ye5u5k3qdlexvt2e0nn0508p04094r9atu2",
						"validator": "celestiavaloper1e2p4u5vqwgum7pm9vhp0yjvl58gvhfc6yfatw4",
					},
				}, {
					Height: 241,
					Time:   ts,
					Type:   "redelegate",
					Data: map[string]any{
						"amount":                "1000utia",
						"completion_time":       "2023-11-21T14:50:46Z",
						"destination_validator": "celestiavaloper15urq2dtp9qce4fyc85m6upwm9xul3049gwdz0x",
						"source_validator":      "celestiavaloper1e2p4u5vqwgum7pm9vhp0yjvl58gvhfc6yfatw4",
					},
				}, {
					Height: 241,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1ze2ye5u5k3qdlexvt2e0nn0508p04094r9atu2",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgBeginRedelegate,
				Height: 241,
				Time:   ts,
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
			redelegation: &storage.Redelegation{
				Time:           ts,
				Height:         241,
				Amount:         decimal.RequireFromString("1000"),
				CompletionTime: time.Date(2023, 11, 21, 14, 50, 46, 0, time.UTC),
				Address: &storage.Address{
					Address:    "celestia1ze2ye5u5k3qdlexvt2e0nn0508p04094r9atu2",
					Height:     241,
					LastHeight: 241,
					Hash:       []byte{0x16, 0x54, 0x4c, 0xd3, 0x94, 0xb4, 0x40, 0xdf, 0xe4, 0xcc, 0x5a, 0xb2, 0xf9, 0xcd, 0xf4, 0x79, 0xc2, 0xfa, 0xbc, 0xb5},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.Zero,
					},
				},
				Source: &storage.Validator{
					Address:           "celestiavaloper1e2p4u5vqwgum7pm9vhp0yjvl58gvhfc6yfatw4",
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
					Stake:             decimal.RequireFromString("-1000"),
				},
				Destination: &storage.Validator{
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
					Stake:             decimal.RequireFromString("1000"),
				},
			},
		}, {
			name: "redelegate test 3",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 315,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgBeginRedelegate",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "44283utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "44283utia",
						"receiver": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "44283utia",
						"recipient": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "44283utia",
						"delegator": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
						"validator": "celestiavaloper1uwmf03ke52vld2sa9khs0nslpgzwsm5xs5e4pn",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "3953utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "3953utia",
						"receiver": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "3953utia",
						"recipient": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "3953utia",
						"delegator": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
						"validator": "celestiavaloper19f0w9svr905fhefusyx4z8sf83j6et0g57nch8",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "redelegate",
					Data: map[string]any{
						"amount":                "125000000utia",
						"completion_time":       "2023-11-21T15:05:39Z",
						"destination_validator": "celestiavaloper19f0w9svr905fhefusyx4z8sf83j6et0g57nch8",
						"source_validator":      "celestiavaloper1uwmf03ke52vld2sa9khs0nslpgzwsm5xs5e4pn",
					},
				}, {
					Height: 315,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgBeginRedelegate,
				Height: 315,
				Time:   ts,
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
			redelegation: &storage.Redelegation{
				Time:           ts,
				Height:         315,
				Amount:         decimal.RequireFromString("125000000"),
				CompletionTime: time.Date(2023, 11, 21, 15, 05, 39, 0, time.UTC),
				Address: &storage.Address{
					Address:    "celestia1mcrnf2n2vmmut5j7upluehstxh4yenztxhmhd0",
					Height:     315,
					LastHeight: 315,
					Hash:       []byte{0xde, 0x07, 0x34, 0xaa, 0x6a, 0x66, 0xf7, 0xc5, 0xd2, 0x5e, 0xe0, 0x7f, 0xcc, 0xde, 0x0b, 0x35, 0xea, 0x4c, 0xcc, 0x4b},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.Zero,
					},
				},
				Source: &storage.Validator{
					Address:           "celestiavaloper1uwmf03ke52vld2sa9khs0nslpgzwsm5xs5e4pn",
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
					Stake:             decimal.RequireFromString("-125000000"),
				},
				Destination: &storage.Validator{
					Address:           "celestiavaloper19f0w9svr905fhefusyx4z8sf83j6et0g57nch8",
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
					Stake:             decimal.RequireFromString("125000000"),
				},
			},
		}, {
			name: "redelegate test 4",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 584,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.staking.v1beta1.MsgBeginRedelegate",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0stake",
						"delegator": "celestia18ccg4scfdee9qvcnxxx4wutjt69sct2jshau9f",
						"validator": "celestiavaloper19mm3s0y676453w3ja58d376ysd84wlf6hq8ae3",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "250000000utia",
						"spender": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "250000000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "250000000utia",
						"recipient": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
						"sender":    "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "redelegate",
					Data: map[string]any{
						"amount":                "250000000utia",
						"completion_time":       "0001-01-01T00:00:00Z",
						"destination_validator": "celestiavaloper1yecxnyegvgm5dwsx0r3jsgr74ju6mlxdwkxx8g",
						"source_validator":      "celestiavaloper19mm3s0y676453w3ja58d376ysd84wlf6hq8ae3",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia18ccg4scfdee9qvcnxxx4wutjt69sct2jshau9f",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgBeginRedelegate,
				Height: 315,
				Time:   ts,
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
			redelegation: &storage.Redelegation{
				Time:           ts,
				Height:         315,
				Amount:         decimal.RequireFromString("250000000"),
				CompletionTime: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				Address: &storage.Address{
					Address:    "celestia18ccg4scfdee9qvcnxxx4wutjt69sct2jshau9f",
					Height:     315,
					LastHeight: 315,
					Hash:       []byte{0x3e, 0x30, 0x8a, 0xc3, 0x09, 0x6e, 0x72, 0x50, 0x33, 0x13, 0x31, 0x8d, 0x57, 0x71, 0x72, 0x5e, 0x8b, 0x0c, 0x2d, 0x52},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.Zero,
					},
				},
				Source: &storage.Validator{
					Address:           "celestiavaloper19mm3s0y676453w3ja58d376ysd84wlf6hq8ae3",
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
					Stake:             decimal.RequireFromString("-250000000"),
				},
				Destination: &storage.Validator{
					Address:           "celestiavaloper1yecxnyegvgm5dwsx0r3jsgr74ju6mlxdwkxx8g",
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
					Stake:             decimal.RequireFromString("250000000"),
				},
			},
		}, {
			name: "redelegate test 5",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 584,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action":    "/cosmos.staking.v1beta1.MsgBeginRedelegate",
						"module":    "staking",
						"msg_index": "0",
						"sender":    "celestia18ccg4scfdee9qvcnxxx4wutjt69sct2jshau9f",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0stake",
						"delegator": "celestia18ccg4scfdee9qvcnxxx4wutjt69sct2jshau9f",
						"validator": "celestiavaloper19mm3s0y676453w3ja58d376ysd84wlf6hq8ae3",
						"msg_index": "0",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":    "250000000utia",
						"spender":   "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"msg_index": "0",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":    "250000000utia",
						"receiver":  "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
						"msg_index": "0",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "250000000utia",
						"recipient": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
						"sender":    "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"msg_index": "0",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"sender":    "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"msg_index": "0",
					},
				}, {
					Height: 584,
					Time:   ts,
					Type:   "redelegate",
					Data: map[string]any{
						"amount":                "250000000utia",
						"completion_time":       "0001-01-01T00:00:00Z",
						"destination_validator": "celestiavaloper1yecxnyegvgm5dwsx0r3jsgr74ju6mlxdwkxx8g",
						"source_validator":      "celestiavaloper19mm3s0y676453w3ja58d376ysd84wlf6hq8ae3",
						"msg_index":             "0",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgBeginRedelegate,
				Height: 315,
				Time:   ts,
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
			redelegation: &storage.Redelegation{
				Time:           ts,
				Height:         315,
				Amount:         decimal.RequireFromString("250000000"),
				CompletionTime: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				Address: &storage.Address{
					Address:    "celestia18ccg4scfdee9qvcnxxx4wutjt69sct2jshau9f",
					Height:     315,
					LastHeight: 315,
					Hash:       []byte{0x3e, 0x30, 0x8a, 0xc3, 0x09, 0x6e, 0x72, 0x50, 0x33, 0x13, 0x31, 0x8d, 0x57, 0x71, 0x72, 0x5e, 0x8b, 0x0c, 0x2d, 0x52},
					Balance: storage.Balance{
						Currency:  currency.Utia,
						Delegated: decimal.Zero,
					},
				},
				Source: &storage.Validator{
					Address:           "celestiavaloper19mm3s0y676453w3ja58d376ysd84wlf6hq8ae3",
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
					Stake:             decimal.RequireFromString("-250000000"),
				},
				Destination: &storage.Validator{
					Address:           "celestiavaloper1yecxnyegvgm5dwsx0r3jsgr74ju6mlxdwkxx8g",
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
					Stake:             decimal.RequireFromString("250000000"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleRedelegate(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
			require.Len(t, tt.ctx.Redelegations, 1)
			require.Equal(t, *tt.redelegation, tt.ctx.Redelegations[0])
		})
	}
}
