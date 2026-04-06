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

func Test_handleCreateSyntheticToken(t *testing.T) {
	ts := time.Now()

	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
		idx    *int
		token  *storage.HLToken
	}{
		{
			name: "test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1036866,
					Time:   ts,
					Type:   "message",
					Data: map[string]string{
						"action": "/hyperlane.warp.v1.MsgCreateSyntheticToken",
					},
				}, {
					Height: 1036866,
					Type:   "hyperlane.warp.v1.EventCreateSyntheticToken",
					Time:   ts,
					Data: map[string]string{ //nolint:gosec
						"msg_index":      "0",
						"origin_denom":   "\"utia\"",
						"origin_mailbox": "\"0x68797065726c616e650000000000000000000000000000000000000000000000\"",
						"owner":          "\"celestia1zvdlcmplx4gdh4hajwlsegnn2xzzfy470gjw4c\"",
						"token_id":       "\"0x726f757465725f61707000000000000000000000000000010000000000000000\"",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgCreateSyntheticToken,
				Height: 1036866,
				Time:   ts,
				Data: types.PackedBytes{
					"OriginDenom":   "utia",
					"OriginMailbox": "0x68797065726c616e650000000000000000000000000000000000000000000000",
					"Owner":         "celestia1zvdlcmplx4gdh4hajwlsegnn2xzzfy470gjw4c",
				},
			},
			idx: testsuite.Ptr(0),
			token: &storage.HLToken{
				Height:           1036866,
				Time:             ts,
				Type:             types.HLTokenTypeSynthetic,
				Denom:            "utia",
				TokenId:          []byte{0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x70, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				SentTransfers:    0,
				ReceiveTransfers: 0,
				Received:         types.NumericZero(),
				Sent:             types.NumericZero(),
				Owner: &storage.Address{
					Address:    "celestia1zvdlcmplx4gdh4hajwlsegnn2xzzfy470gjw4c",
					Height:     1036866,
					LastHeight: 1036866,
					Balance:    storage.EmptyBalance(),
					Hash:       []byte{0x13, 0x1b, 0xfc, 0x6c, 0x3f, 0x35, 0x50, 0xdb, 0xd6, 0xfd, 0x93, 0xbf, 0x0c, 0xa2, 0x73, 0x51, 0x84, 0x24, 0x92, 0xbe},
				},
				Mailbox: &storage.HLMailbox{
					Height:     1036866,
					Time:       ts,
					InternalId: 0,
					Mailbox:    []byte{0x68, 0x79, 0x70, 0x65, 0x72, 0x6c, 0x61, 0x6e, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
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
			err := handleCreateSyntheticToken(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
			require.NotEmpty(t, tt.ctx.HlTokens.Len())

			_ = tt.ctx.HlTokens.Range(func(_ string, value *storage.HLToken) (error, bool) {
				require.NotNil(t, value)
				require.Equal(t, tt.token, value)
				return nil, false
			})
		})
	}
}
