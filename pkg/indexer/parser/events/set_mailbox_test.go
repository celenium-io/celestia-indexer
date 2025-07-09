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

func Test_handleSetMailbox(t *testing.T) {
	ts := time.Now()
	tests := []struct {
		name    string
		ctx     *context.Context
		events  []storage.Event
		msg     *storage.Message
		idx     *int
		mailbox *storage.HLMailbox
	}{
		{
			name: "test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1036866,
					Time:   ts,
					Type:   "message",
					Data: map[string]any{
						"action": "/hyperlane.core.v1.MsgSetMailbox",
					},
				}, {
					Height: 1036866,
					Time:   ts,
					Type:   "hyperlane.core.v1.EventSetMailbox",
					Data: map[string]any{
						"default_hook":       "\"0x726f757465725f706f73745f6469737061746368000000040000000000000001\"",
						"default_ism":        "null",
						"mailbox_id":         "\"0x68797065726c616e650000000000000000000000000000000000000000000000\"",
						"msg_index":          "0",
						"new_owner":          "\"\"",
						"owner":              "\"celestia1zvdlcmplx4gdh4hajwlsegnn2xzzfy470gjw4c\"",
						"renounce_ownership": "false",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgSetMailbox,
				Height: 1036866,
				Time:   ts,
			},
			idx: testsuite.Ptr(0),
			mailbox: &storage.HLMailbox{
				Height:      1036866,
				Time:        ts,
				Mailbox:     []byte{0x68, 0x79, 0x70, 0x65, 0x72, 0x6c, 0x61, 0x6e, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				DefaultHook: []byte{0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
				Owner: &storage.Address{
					Address: "celestia1zvdlcmplx4gdh4hajwlsegnn2xzzfy470gjw4c",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ctx.Block = &storage.Block{
				Height: 1036866,
				Time:   ts,
			}
			err := handleSetMailbox(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
			require.NotNil(t, tt.msg.HLMailbox)
			require.NotNil(t, tt.mailbox, tt.msg.HLMailbox)
		})
	}
}
