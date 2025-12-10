package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func TestSendHandler(t *testing.T) {
	tests := []struct {
		name       string
		ctx        *context.Context
		events     []storage.Event
		msg        *storage.Message
		idx        *int
		requireIdx int
	}{
		{
			name: "send test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.bank.v1beta1.MsgSend",
					},
				}, {
					Height: 1745041,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "6000000000utia",
						"spender": "celestia1j64nlm43umrrv62krql2ek48dfdjkh38d6g4y0",
					},
				}, {
					Height: 1745041,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "6000000000utia",
						"receiver": "celestia1yj3z5paheg78h2hcfyc3fc68h3r0jl3zc8n5y0",
					},
				}, {
					Height: 1745041,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "6000000000utia",
						"recipient": "celestia1yj3z5paheg78h2hcfyc3fc68h3r0jl3zc8n5y0",
						"sender":    "celestia1j64nlm43umrrv62krql2ek48dfdjkh38d6g4y0",
					},
				}, {
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1j64nlm43umrrv62krql2ek48dfdjkh38d6g4y0",
					},
				}, {
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"module": "bank",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgSend,
				Height: 1745041,
			},
			idx:        testsuite.Ptr(0),
			requireIdx: 6,
		}, {
			name: "send test 2",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"action":    "/cosmos.bank.v1beta1.MsgSend",
						"msg_index": "0",
					},
				}, {
					Height: 1745041,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":    "6000000000utia",
						"spender":   "celestia1j64nlm43umrrv62krql2ek48dfdjkh38d6g4y0",
						"msg_index": "0",
					},
				}, {
					Height: 1745041,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":    "6000000000utia",
						"receiver":  "celestia1yj3z5paheg78h2hcfyc3fc68h3r0jl3zc8n5y0",
						"msg_index": "0",
					},
				}, {
					Height: 1745041,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "6000000000utia",
						"recipient": "celestia1yj3z5paheg78h2hcfyc3fc68h3r0jl3zc8n5y0",
						"sender":    "celestia1j64nlm43umrrv62krql2ek48dfdjkh38d6g4y0",
						"msg_index": "0",
					},
				}, {
					Height: 1745041,
					Type:   "message",
					Data: map[string]any{
						"sender":    "celestia1j64nlm43umrrv62krql2ek48dfdjkh38d6g4y0",
						"msg_index": "0",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgSend,
				Height: 1745041,
			},
			idx:        testsuite.Ptr(0),
			requireIdx: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleSend(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
			require.Equal(t, tt.requireIdx, *tt.idx)
		})
	}
}
