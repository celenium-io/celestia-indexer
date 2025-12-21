// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

// Test_Handle_boundarySafety tests that Handle function properly handles boundary conditions
func Test_Handle_boundarySafety(t *testing.T) {
	tests := []struct {
		name        string
		events      []storage.Event
		msg         *storage.Message
		idx         int
		expectError bool
	}{
		{
			name: "two events - boundary safe",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]any{
						"action": "/unknown.msg.Type",
					},
				},
				{
					Height: 100,
					Type:   "some_event",
					Data: map[string]any{
						"action": "/next.action",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         0,
			expectError: false,
		},
		{
			name:   "empty events array",
			events: []storage.Event{},
			msg: &storage.Message{
				Type:   types.MsgSend,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         0,
			expectError: true, // will fail because events[0] doesn't exist
		},
		{
			name: "unhandled message type at end of events",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]any{
						"action": "/unknown.msg.Type",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         0,
			expectError: false,
		},
		{
			name: "multiple events with unknown message type",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]any{
						"action": "/unknown.msg.Type",
					},
				},
				{
					Height: 100,
					Type:   "some_event",
					Data:   map[string]any{},
				},
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]any{
						"action": "/another.action",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			idx := tt.idx

			if tt.expectError {
				require.Panics(t, func() {
					_ = Handle(ctx, tt.events, tt.msg, &idx)
				})
			} else {
				require.NotPanics(t, func() {
					err := Handle(ctx, tt.events, tt.msg, &idx)
					require.NoError(t, err)
				})
			}
		})
	}
}

// Test_handle_sliceIterationSafety tests the internal handle function for slice iteration safety
func Test_handle_sliceIterationSafety(t *testing.T) {
	tests := []struct {
		name        string
		events      []storage.Event
		msg         *storage.Message
		idx         int
		handlers    map[types.MsgType]EventHandler
		stopKey     string
		expectPanic bool
	}{
		{
			name: "iteration reaches end safely",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]any{
						"action": "/test",
					},
				},
				{
					Height: 100,
					Type:   "other_event",
					Data:   map[string]any{},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         0,
			handlers:    map[types.MsgType]EventHandler{},
			stopKey:     "action",
			expectPanic: false,
		},
		{
			name: "starts at last element",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]any{
						"action": "/test",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         0,
			handlers:    map[types.MsgType]EventHandler{},
			stopKey:     "action",
			expectPanic: false,
		},
		{
			name: "no message events to stop at",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]any{
						"action": "/test",
					},
				},
				{
					Height: 100,
					Type:   "coin_spent",
					Data:   map[string]any{},
				},
				{
					Height: 100,
					Type:   "transfer",
					Data:   map[string]any{},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         0,
			handlers:    map[types.MsgType]EventHandler{},
			stopKey:     "action",
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			idx := tt.idx

			if tt.expectPanic {
				require.Panics(t, func() {
					_ = handle(ctx, tt.events, tt.msg, &idx, tt.handlers, tt.stopKey)
				})
			} else {
				require.NotPanics(t, func() {
					err := handle(ctx, tt.events, tt.msg, &idx, tt.handlers, tt.stopKey)
					require.NoError(t, err)
				})
			}
		})
	}
}

// Test_toTheNextAction_incrementSafety ensures toTheNextAction doesn't cause index out of bounds
func Test_toTheNextAction_incrementSafety(t *testing.T) {
	tests := []struct {
		name      string
		events    []storage.Event
		startIdx  int
		expectIdx int
	}{
		{
			name:      "empty slice",
			events:    []storage.Event{},
			startIdx:  0,
			expectIdx: 0, // should not increment
		},
		{
			name: "single element with empty action",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]any{
						"action": "",
					},
				},
			},
			startIdx:  0,
			expectIdx: 0, // at boundary, should not increment
		},
		{
			name: "two elements, start at first",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]any{
						"action": "",
					},
				},
				{
					Type: "message",
					Data: map[string]any{
						"action": "/some.action",
					},
				},
			},
			startIdx:  0,
			expectIdx: 1,
		},
		{
			name: "start at penultimate position with empty action",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]any{
						"action": "",
					},
				},
				{
					Type: "message",
					Data: map[string]any{
						"action": "",
					},
				},
				{
					Type: "message",
					Data: map[string]any{
						"action": "/action",
					},
				},
			},
			startIdx:  1,
			expectIdx: 2,
		},
		{
			name: "all empty actions",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]any{
						"action": "",
					},
				},
				{
					Type: "message",
					Data: map[string]any{
						"action": "",
					},
				},
			},
			startIdx:  0,
			expectIdx: 1, // stops at last element
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tt.startIdx
			require.NotPanics(t, func() {
				toTheNextAction(tt.events, &idx)
			})
			require.Equal(t, tt.expectIdx, idx, "index should match expected value")
		})
	}
}

// Test_recvPacket_incrementBy2Safety tests the safety of incrementing index by 2
func Test_recvPacket_incrementBy2Safety(t *testing.T) {
	tests := []struct {
		name   string
		events []storage.Event
		msg    *storage.Message
		idx    int
	}{
		{
			name: "exactly 2 elements after current",
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				},
				{
					Height: 100,
					Type:   types.EventTypeRecvPacket,
					Data:   map[string]any{},
				},
				{
					Height: 100,
					Type:   "some_event",
					Data:   map[string]any{},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgRecvPacket,
				Height: 100,
				Data: map[string]any{
					"Packet": map[string]any{
						"DestinationPort": "other",
					},
				},
			},
			idx: 0,
		},
		{
			name: "only 1 element after increment",
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				},
				{
					Height: 100,
					Type:   types.EventTypeRecvPacket,
					Data:   map[string]any{},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgRecvPacket,
				Height: 100,
				Data: map[string]any{
					"Packet": map[string]any{
						"DestinationPort": "other",
					},
				},
			},
			idx: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			idx := tt.idx
			require.NotPanics(t, func() {
				err := handleRecvPacket(ctx, tt.events, tt.msg, &idx)
				require.NoError(t, err)
			})
		})
	}
}
