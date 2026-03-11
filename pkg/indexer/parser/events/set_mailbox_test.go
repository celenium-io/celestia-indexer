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
					Address:    "celestia1zvdlcmplx4gdh4hajwlsegnn2xzzfy470gjw4c",
					Height:     1036866,
					LastHeight: 1036866,
					Balance:    storage.EmptyBalance(),
					Hash:       []byte{0x13, 0x1b, 0xfc, 0x6c, 0x3f, 0x35, 0x50, 0xdb, 0xd6, 0xfd, 0x93, 0xbf, 0x0c, 0xa2, 0x73, 0x51, 0x84, 0x24, 0x92, 0xbe},
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
			require.NotEmpty(t, tt.ctx.HlMailboxes.Len())

			mailbox, ok := tt.ctx.HlMailboxes.Get(0)
			require.True(t, ok)
			require.EqualValues(t, tt.mailbox, mailbox)
		})
	}
}

func Test_handleSetMailbox_newOwner(t *testing.T) {
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
						"new_owner":          "\"celestia1lg0e9n4pt29lpq2k4ptue4ckw09dx0aujlpe4j\"",
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
					Address:    "celestia1lg0e9n4pt29lpq2k4ptue4ckw09dx0aujlpe4j",
					Height:     1036866,
					LastHeight: 1036866,
					Balance:    storage.EmptyBalance(),
					Hash:       []byte{0xfa, 0x1f, 0x92, 0xce, 0xa1, 0x5a, 0x8b, 0xf0, 0x81, 0x56, 0xa8, 0x57, 0xcc, 0xd7, 0x16, 0x73, 0xca, 0xd3, 0x3f, 0xbc},
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
			require.NotEmpty(t, tt.ctx.HlMailboxes.Len())

			mailbox, ok := tt.ctx.HlMailboxes.Get(0)
			require.True(t, ok)
			require.EqualValues(t, tt.mailbox, mailbox)
		})
	}
}
