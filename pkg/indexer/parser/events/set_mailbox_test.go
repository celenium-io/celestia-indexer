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

func Test_handleSetMailbox(t *testing.T) {
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
						"action": "/hyperlane.core.v1.MsgSetMailbox",
					},
				}, {
					Height: 1036866,
					Type:   "hyperlane.core.v1.EventSetMailbox",
					Data: map[string]any{
						"owner":              "celestia1ul4nkg590xsf8cpn60z0gmjxmwuxn9afzar42t",
						"msg_index":          "0",
						"new_owner":          "",
						"mailbox_id":         "0x68797065726c616e650000000000000000000000000000000000000000000000",
						"default_ism":        "null",
						"default_hook":       "0x726f757465725f706f73745f6469737061746368000000040000000000000000",
						"renounce_ownership": "false",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgSetMailbox,
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
				err := handleSetMailbox(tt.ctx, tt.events, tt.msg[i], tt.idx)
				require.NoError(t, err)
				require.NotNil(t, tt.msg[i].HLMailbox)
			}
		})
	}
}
