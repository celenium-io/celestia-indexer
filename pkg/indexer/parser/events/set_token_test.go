// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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

func Test_handleSetToken(t *testing.T) {
	tests := []struct {
		name      string
		ctx       *context.Context
		events    []storage.Event
		msg       []*storage.Message
		wantEmpty bool
		idx       *int
	}{
		{
			name: "test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1036866,
					Type:   "message",
					Data: map[string]any{
						"action": "/hyperlane.warp.v1.MsgSetToken",
					},
				}, {
					Height: 1036866,
					Type:   "hyperlane.warp.v1.EventSetToken",
					Data: map[string]any{
						"owner":              "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
						"ism_id":             "",
						"token_id":           "0x726f757465725f61707000000000000000000000000000010000000000000000",
						"msg_index":          "0",
						"new_owner":          "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
						"renounce_ownership": "false",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgSetToken,
					Height: 1036866,
				},
			},
			idx:       testsuite.Ptr(0),
			wantEmpty: false,
		}, {
			name: "test 2: empty new owner",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1036866,
					Type:   "message",
					Data: map[string]any{
						"action": "/hyperlane.warp.v1.MsgSetToken",
					},
				}, {
					Height: 1036866,
					Type:   "hyperlane.warp.v1.EventSetToken",
					Data: map[string]any{
						"owner":              "celestia1lg0e9n4pt29lpq2k4ptue4ckw09dx0aujlpe4j",
						"ism_id":             "0x726f757465725f69736d00000000000000000000000000000000000000000001",
						"token_id":           "0x726f757465725f61707000000000000000000000000000020000000000000001",
						"msg_index":          "0",
						"new_owner":          "",
						"renounce_ownership": "false",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgSetToken,
					Height: 1036866,
				},
			},
			idx:       testsuite.Ptr(0),
			wantEmpty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ctx.Block = &storage.Block{
				Height: 1036866,
				Time:   time.Now(),
			}
			for i := range tt.msg {
				err := handleSetToken(tt.ctx, tt.events, tt.msg[i], tt.idx)
				require.NoError(t, err)
				if !tt.wantEmpty {
					require.NotNil(t, tt.msg[i].HLToken)
				} else {
					require.Nil(t, tt.msg[i].HLToken)
				}
			}
		})
	}
}
