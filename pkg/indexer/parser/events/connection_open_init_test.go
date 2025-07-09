// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func Test_handleConnectionOpenInit(t *testing.T) {
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
						"action": "/ibc.core.connection.v1.MsgConnectionOpenTry",
					},
				}, {
					Height: 1036866,
					Type:   "connection_open_try",
					Data: map[string]any{
						"client_id":                  "07-tendermint-184",
						"connection_id":              "connection-143",
						"counterparty_client_id":     "07-tendermint-0",
						"counterparty_connection_id": "connection-0",
					},
				}, {
					Height: 1036866,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_connection",
					},
				},
				{
					Height: 1036866,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.connection.v1.MsgConnectionOpenInit",
					},
				}, {
					Height: 1036866,
					Type:   "connection_open_init",
					Data: map[string]any{
						"client_id":                  "07-tendermint-184",
						"connection_id":              "connection-144",
						"counterparty_client_id":     "07-tendermint-0",
						"counterparty_connection_id": "",
					},
				}, {
					Height: 1036866,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_connection",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgConnectionOpenTry,
					Height: 1036866,
				}, {
					Type:   types.MsgConnectionOpenInit,
					Height: 1036866,
				},
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.msg {
				err := handleConnectionOpenInit(tt.ctx, tt.events, tt.msg[i], tt.idx)
				require.NoError(t, err)
				require.NotNil(t, tt.msg[i].IbcConnection)
			}
		})
	}
}
