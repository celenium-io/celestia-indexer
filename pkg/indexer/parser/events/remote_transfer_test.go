// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func Test_handleHyperlaneRemoteTransfer(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    []*storage.Message
		idx    *int
	}{
		{
			name: "test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1036866,
					Type:   "message",
					Data: map[string]any{
						"action": "/hyperlane.warp.v1.MsgRemoteTransfer",
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
					Height: 1036866,
					Type:   "hyperlane.warp.v1.EventSendRemoteTransfer",
					Data: map[string]any{
						"amount":             "6745utia",
						"sender":             "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
						"token_id":           "0x726f757465725f61707000000000000000000000000000010000000000000000",
						"msg_index":          "0",
						"recipient":          "0x000000000000000000000000e0f9f661f106d6da1974fdc12a904e936834b3f8",
						"destination_domain": "56",
					},
				}, {
					Height: 1036866,
					Type:   "hyperlane.core.v1.EventDispatch",
					Data: map[string]any{
						"sender":            "0x726f757465725f61707000000000000000000000000000010000000000000000",
						"message":           "0x03000000e86d696c6b726f757465725f61707000000000000000000000000000010000000000000000000000380000000000000000000000007b4bf9feccff207ef2cb7101ceb15b8516021acd000000000000000000000000e0f9f661f106d6da1974fdc12a904e936834b3f8000000000000000000000000000000000000000000000000000000000c939ac0",
						"msg_index":         "0",
						"recipient":         "0x0000000000000000000000007b4bf9feccff207ef2cb7101ceb15b8516021acd",
						"destination":       "56",
						"origin_mailbox_id": "0x68797065726c616e650000000000000000000000000000000000000000000000",
					},
				}, {
					Height: 1036866,
					Type:   "hyperlane.core.post_dispatch.v1.EventInsertedIntoTree",
					Data: map[string]any{
						"index":               "232",
						"msg_index":           "0",
						"message_id":          "0xdcdb3f985ecd20c313c58c0f6b2a0d7ea980349134ee4813f6bd53cfe5bf0a1e",
						"merkle_tree_hook_id": "0x726f757465725f706f73745f6469737061746368000000030000000000000001",
					},
				}, {
					Height: 841682,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "4379232utia",
						"spender": "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
					},
				}, {
					Height: 841682,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":  "4379232utia",
						"spender": "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
					},
				}, {
					Height: 841682,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "4379232utia",
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
					Height: 1036866,
					Type:   "hyperlane.core.post_dispatch.v1.EventGasPayment",
					Data: map[string]any{
						"igp_id":      "0x726f757465725f706f73745f6469737061746368000000040000000000000000",
						"payment":     "4379232utia",
						"msg_index":   "0",
						"gas_amount":  "64000",
						"message_id":  "0xdcdb3f985ecd20c313c58c0f6b2a0d7ea980349134ee4813f6bd53cfe5bf0a1e",
						"destination": "56",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgRemoteTransfer,
					Height: 1036866,
				},
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ctx.Block = &storage.Block{
				Height: 1036866,
				Time:   time.Now(),
			}
			for i := range tt.msg {
				err := handleHyperlaneRemoteTransfer(tt.ctx, tt.events, tt.msg[i], tt.idx)
				require.NoError(t, err)
				require.NotNil(t, tt.msg[i].HLTransfer)
			}
		})
	}
}
