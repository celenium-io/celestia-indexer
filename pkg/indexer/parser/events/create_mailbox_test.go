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

func Test_handleCreateMailbox(t *testing.T) {
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
						"action": "/hyperlane.core.v1.MsgCreateMailbox",
					},
				}, {
					Height: 1036866,
					Type:   "hyperlane.core.v1.EventCreateMailbox",
					Data: map[string]any{
						"owner":         "\"celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj\"",
						"mailbox_id":    "\"0x68797065726c616e650000000000000000000000000000000000000000000000\"",
						"default_ism":   "\"0x726f757465725f69736d00000000000000000000000000010000000000000000\"",
						"default_hook":  "null",
						"local_domain":  "1835625579",
						"required_hook": "null",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgCreateMailbox,
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
				err := handleCreateMailbox(tt.ctx, tt.events, tt.msg[i], tt.idx)
				require.NoError(t, err)
				require.NotNil(t, tt.msg[i].HLMailbox)
			}
		})
	}
}
