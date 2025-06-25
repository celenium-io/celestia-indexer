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

func Test_handleHyperlaneProcessMessage(t *testing.T) {
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
						"action": "/hyperlane.core.v1.MsgProcessMessage",
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
					Type:   "hyperlane.warp.v1.EventReceiveRemoteTransfer",
					Data: map[string]any{
						"amount":        "6745utia",
						"sender":        "0x0000000000000000000000007b4bf9feccff207ef2cb7101ceb15b8516021acd",
						"token_id":      "0x726f757465725f61707000000000000000000000000000010000000000000000",
						"msg_index":     "0",
						"recipient":     "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
						"origin_domain": "56",
					},
				}, {
					Height: 1036866,
					Type:   "hyperlane.core.v1.EventProcess",
					Data: map[string]any{
						"origin":            "56",
						"sender":            "0x0000000000000000000000007b4bf9feccff207ef2cb7101ceb15b8516021acd",
						"message":           "0x0300046a65000000380000000000000000000000007b4bf9feccff207ef2cb7101ceb15b8516021acd6d696c6b726f757465725f6170700000000000000000000000000001000000000000000000000000000000000000000056a3cc5141289679223bc85169141c7d454143900000000000000000000000000000000000000000000000000000000002faf080",
						"msg_index":         "0",
						"recipient":         "0x726f757465725f61707000000000000000000000000000010000000000000000",
						"message_id":        "0xd3aa80e8a5f8082cba208900a4ad3aec15516f423a89ef7803a9da35e2e58369",
						"origin_mailbox_id": "0x68797065726c616e650000000000000000000000000000000000000000000000",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgProcessMessage,
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
				err := handleHyperlaneProcessMessage(tt.ctx, tt.events, tt.msg[i], tt.idx)
				require.NoError(t, err)
				require.NotNil(t, tt.msg[i].HLTransfer)
				require.NotNil(t, tt.msg[i].HLTransfer.Address)
			}
		})
	}
}
