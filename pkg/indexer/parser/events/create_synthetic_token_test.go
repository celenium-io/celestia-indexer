// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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

func Test_handleCreateSyntheticToken(t *testing.T) {
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
						"action": "/hyperlane.warp.v1.MsgCreateSyntheticToken",
					},
				}, {
					Height: 1036866,
					Type:   "hyperlane.warp.v1.EventCreateSyntheticToken",
					Data: map[string]any{
						"owner":          "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
						"token_id":       "0x726f757465725f61707000000000000000000000000000010000000000000000",
						"msg_index":      "0",
						"origin_denom":   "utia",
						"origin_mailbox": "0x68797065726c616e650000000000000000000000000000000000000000000000",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgCreateSyntheticToken,
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
				err := handleCreateSyntheticToken(tt.ctx, tt.events, tt.msg[i], tt.idx)
				require.NoError(t, err)
				require.NotNil(t, tt.msg[i].HLToken)
			}
		})
	}
}
