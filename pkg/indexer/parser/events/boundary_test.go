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

// Test_recvPacket_boundaryChecks tests boundary conditions for recv_packet event handler
func Test_recvPacket_boundaryChecks(t *testing.T) {
	tests := []struct {
		name        string
		ctx         *context.Context
		events      []storage.Event
		msg         *storage.Message
		expectError bool
	}{
		{
			name: "empty events array",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]string{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgRecvPacket,
				Height: 100,
				Data:   map[string]any{},
			},
			expectError: false, // should handle gracefully
		},
		{
			name: "index at last element",
			ctx:  context.NewContext(),
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
					Type:   "message",
					Data: map[string]string{
						"action": "",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgRecvPacket,
				Height: 100,
				Data:   map[string]any{},
			},
			expectError: false,
		},
		{
			name: "minimal events - boundary check only",
			ctx:  context.NewContext(),
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
					Type:   types.EventTypeMessage,
					Data: map[string]string{
						"action": "",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgRecvPacket,
				Height: 100,
				Data:   map[string]any{},
			},
			expectError: false, // should not panic on boundary
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCursor(tt.events)
			err := handleRecvPacket(tt.ctx, c, tt.msg)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test_acknowledgement_boundaryChecks tests boundary conditions for acknowledgement event handler
func Test_acknowledgement_boundaryChecks(t *testing.T) {
	tests := []struct {
		name        string
		ctx         *context.Context
		events      []storage.Event
		msg         *storage.Message
		expectError bool
	}{
		{
			name: "empty events after increment",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]string{
						"action": "/ibc.core.channel.v1.MsgAcknowledgement",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgAcknowledgement,
				Height: 100,
				Data:   map[string]any{},
			},
			expectError: false,
		},
		{
			name: "index near end of array",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]string{
						"action": "/ibc.core.channel.v1.MsgAcknowledgement",
					},
				},
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data:   map[string]string{},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgAcknowledgement,
				Height: 100,
				Data:   map[string]any{},
			},
			expectError: false,
		},
		{
			name: "loop boundary check - no transfer case",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]string{
						"action": "/ibc.core.channel.v1.MsgAcknowledgement",
					},
				},
				{
					Height: 100,
					Type:   "some_event",
					Data: map[string]string{
						"action": "",
					},
				},
				{
					Height: 100,
					Type:   "some_event",
					Data: map[string]string{
						"action": "",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgAcknowledgement,
				Height: 100,
				Data: map[string]any{
					"Packet": map[string]any{
						"SourcePort": "icahost",
						"Data":       map[string]any{},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCursor(tt.events)
			err := handleAcknowledgement(tt.ctx, c, tt.msg)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test_exec_outOfBounds tests out of bounds access in exec handler
func Test_exec_outOfBounds(t *testing.T) {
	tests := []struct {
		name        string
		idx         int
		data        map[string]any
		expectError bool
	}{
		{
			name: "valid index",
			idx:  0,
			data: map[string]any{
				"Msgs": []any{
					map[string]any{"key": "value"},
				},
			},
			expectError: false,
		},
		{
			name: "index out of bounds - too large",
			idx:  5,
			data: map[string]any{
				"Msgs": []any{
					map[string]any{"key": "value"},
				},
			},
			expectError: true,
		},
		{
			name: "negative index",
			idx:  -1,
			data: map[string]any{
				"Msgs": []any{
					map[string]any{"key": "value"},
				},
			},
			expectError: true,
		},
		{
			name: "empty array",
			idx:  0,
			data: map[string]any{
				"Msgs": []any{},
			},
			expectError: true,
		},
		{
			name: "missing Msgs key",
			idx:  0,
			data: map[string]any{
				"Other": "data",
			},
			expectError: true,
		},
		{
			name: "Msgs is not an array",
			idx:  0,
			data: map[string]any{
				"Msgs": "not_an_array",
			},
			expectError: true,
		},
		{
			name: "array element is not a map",
			idx:  0,
			data: map[string]any{
				"Msgs": []any{
					"not_a_map",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getInternalDataForExec(tt.data, tt.idx)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test_toTheNextAction_boundaryChecks tests Cursor.SkipToNext("action") boundary
// conditions (this used to test the standalone toTheNextAction helper, which
// SkipToNext replaced).
func Test_toTheNextAction_boundaryChecks(t *testing.T) {
	tests := []struct {
		name           string
		events         []storage.Event
		wantAtEnd      bool // Peek() should report false after SkipToNext
		wantActionData string
	}{
		{
			name:      "index at end of array",
			events:    []storage.Event{},
			wantAtEnd: true,
		},
		{
			name: "single event with empty action is drained, cursor ends at EOF",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]string{
						"action": "",
					},
				},
			},
			wantAtEnd: true,
		},
		{
			name: "normal progression stops at the boundary event, unconsumed",
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
			wantAtEnd:      false,
			wantActionData: "/some.action",
		},
		{
			name: "all empty actions until end drains everything",
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
			wantAtEnd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCursor(tt.events)
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

// Test_handlers_nilChecks tests nil pointer checks
func Test_handlers_nilChecks(t *testing.T) {
	ctx := context.NewContext()
	events := []storage.Event{
		{
			Height: 100,
			Type:   "message",
			Data: map[string]string{
				"action": "/ibc.core.channel.v1.MsgRecvPacket",
			},
		},
	}
	msg := &storage.Message{
		Type:   types.MsgRecvPacket,
		Height: 100,
		Data:   map[string]any{},
	}

	t.Run("nil cursor", func(t *testing.T) {
		err := handleRecvPacket(ctx, nil, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nil event cursor")
	})

	t.Run("nil message pointer", func(t *testing.T) {
		c := NewCursor(events)
		err := handleRecvPacket(ctx, c, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nil message")
	})
}

// Test_acknowledgement_loopIncrementSafety tests that loop increments don't cause panics
func Test_acknowledgement_loopIncrementSafety(t *testing.T) {
	tests := []struct {
		name   string
		events []storage.Event
		msg    *storage.Message
	}{
		{
			name: "increment at exact boundary",
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]string{
						"action": "/ibc.core.channel.v1.MsgAcknowledgement",
					},
				},
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]string{
						"action": "",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgAcknowledgement,
				Height: 100,
				Data: map[string]any{
					"Packet": map[string]any{
						"SourcePort": "other", // not transfer, so no address validation
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
				err := handleAcknowledgement(ctx, c, tt.msg)
				require.NoError(t, err)
			})
		})
	}
}
