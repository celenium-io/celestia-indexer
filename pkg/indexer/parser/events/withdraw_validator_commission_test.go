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

func Test_handleWithdrawValidatorCommission(t *testing.T) {
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
					Height: 848613,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
					},
				}, {
					Height: 848613,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "3003622utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 848613,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "3003622utia",
						"receiver": "celestia1s0lankh33kprer2l22nank5rvsuh9ksa0e65uz",
					},
				}, {
					Height: 848613,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "3003622utia",
						"recipient": "celestia1s0lankh33kprer2l22nank5rvsuh9ksa0e65uz",
						"sender":    "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 848613,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 848613,
					Type:   "withdraw_commission",
					Data: map[string]any{
						"amount": "3003622utia",
					},
				}, {
					Height: 848613,
					Type:   "message",
					Data: map[string]any{
						"module": "distribution",
						"sender": "celestiavaloper1s0lankh33kprer2l22nank5rvsuh9ksa2xcd2y",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgDelegate,
				Height: 848613,
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleWithdrawValidatorCommission(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
		})
	}
}
