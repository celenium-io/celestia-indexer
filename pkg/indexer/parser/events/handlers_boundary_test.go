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
		expectError bool
	}{
		{
			name: "two events - boundary safe",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]string{
						"action": "/unknown.msg.Type",
					},
				},
				{
					Height: 100,
					Type:   "some_event",
					Data: map[string]string{
						"action": "/next.action",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			expectError: false,
		},
		{
			// Previously this panicked: handleSend indexed events[0] directly on
			// an empty slice. Cursor.Peek() can't go out of bounds, so this is
			// now a graceful "unexpected event action" error, not a crash.
			name:        "empty events array",
			events:      []storage.Event{},
			msg:         &storage.Message{Type: types.MsgSend, Height: 100, Data: map[string]any{}},
			expectError: true,
		},
		{
			name: "unhandled message type at end of events",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]string{
						"action": "/unknown.msg.Type",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			expectError: false,
		},
		{
			name: "multiple events with unknown message type",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]string{
						"action": "/unknown.msg.Type",
					},
				},
				{
					Height: 100,
					Type:   "some_event",
					Data:   map[string]string{},
				},
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]string{
						"action": "/another.action",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			c := NewCursor(tt.events)

			var err error
			require.NotPanics(t, func() {
				err = Handle(ctx, c, tt.msg)
			})
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test_handle_sliceIterationSafety tests the internal handle function for slice iteration safety
func Test_handle_sliceIterationSafety(t *testing.T) {
	tests := []struct {
		name     string
		events   []storage.Event
		msg      *storage.Message
		handlers map[types.MsgType]EventHandler
		stopKey  string
	}{
		{
			name: "iteration reaches end safely",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]string{
						"action": "/test",
					},
				},
				{
					Height: 100,
					Type:   "other_event",
					Data:   map[string]string{},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			handlers: map[types.MsgType]EventHandler{},
			stopKey:  "action",
		},
		{
			name: "starts at last element",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]string{
						"action": "/test",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			handlers: map[types.MsgType]EventHandler{},
			stopKey:  "action",
		},
		{
			name: "no message events to stop at",
			events: []storage.Event{
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]string{
						"action": "/test",
					},
				},
				{
					Height: 100,
					Type:   "coin_spent",
					Data:   map[string]string{},
				},
				{
					Height: 100,
					Type:   "transfer",
					Data:   map[string]string{},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnknown,
				Height: 100,
				Data:   map[string]any{},
			},
			handlers: map[types.MsgType]EventHandler{},
			stopKey:  "action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			c := NewCursor(tt.events)

			require.NotPanics(t, func() {
				err := handle(ctx, c, tt.msg, tt.handlers, tt.stopKey)
				require.NoError(t, err)
			})
		})
	}
}

// Test_toTheNextAction_incrementSafety ensures Cursor.SkipToNext("action")
// doesn't cause index out of bounds (this used to test the standalone
// toTheNextAction helper, which SkipToNext replaced).
func Test_toTheNextAction_incrementSafety(t *testing.T) {
	tests := []struct {
		name           string
		events         []storage.Event
		startIdx       int
		wantAtEnd      bool
		wantActionData string
	}{
		{
			name:      "empty slice",
			events:    []storage.Event{},
			startIdx:  0,
			wantAtEnd: true,
		},
		{
			// SkipToNext drains an event with an empty action (it isn't a
			// boundary) instead of stopping short of it, so the cursor lands
			// at the true end here rather than staying put.
			name: "single element with empty action",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]string{
						"action": "",
					},
				},
			},
			startIdx:  0,
			wantAtEnd: true,
		},
		{
			name: "two elements, start at first",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]string{
						"action": "",
					},
				},
				{
					Type: "message",
					Data: map[string]string{
						"action": "/some.action",
					},
				},
			},
			startIdx:       0,
			wantActionData: "/some.action",
		},
		{
			name: "start at penultimate position with empty action",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]string{
						"action": "",
					},
				},
				{
					Type: "message",
					Data: map[string]string{
						"action": "",
					},
				},
				{
					Type: "message",
					Data: map[string]string{
						"action": "/action",
					},
				},
			},
			startIdx:       1,
			wantActionData: "/action",
		},
		{
			// Same canonicalization as the "single element" case above: both
			// remaining events have an empty action, so SkipToNext drains
			// them all and lands at the true end.
			name: "all empty actions",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]string{
						"action": "",
					},
				},
				{
					Type: "message",
					Data: map[string]string{
						"action": "",
					},
				},
			},
			startIdx:  0,
			wantAtEnd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCursor(tt.events)
			c.Skip(tt.startIdx)

			require.NotPanics(t, func() {
				c.SkipToNext("action")
			})

			event, ok := c.Peek()
			if tt.wantAtEnd {
				require.False(t, ok, "cursor should have reached the end")
				return
			}
			require.True(t, ok)
			require.Equal(t, tt.wantActionData, event.Data["action"])
		})
	}
}

// Test_recvPacket_incrementBy2Safety tests the safety of incrementing index by 2
func Test_recvPacket_incrementBy2Safety(t *testing.T) {
	tests := []struct {
		name   string
		events []storage.Event
		msg    *storage.Message
	}{
		{
			name: "exactly 2 elements after current",
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]string{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				},
				{
					Height: 100,
					Type:   types.EventTypeRecvPacket,
					Data:   map[string]string{},
				},
				{
					Height: 100,
					Type:   "some_event",
					Data:   map[string]string{},
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
		},
		{
			name: "only 1 element after increment",
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]string{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				},
				{
					Height: 100,
					Type:   types.EventTypeRecvPacket,
					Data:   map[string]string{},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			c := NewCursor(tt.events)
			require.NotPanics(t, func() {
				err := handleRecvPacket(ctx, c, tt.msg)
				require.NoError(t, err)
			})
		})
	}
}
