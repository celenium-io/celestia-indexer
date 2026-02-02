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

// Test_recvPacket_boundaryChecks tests boundary conditions for recv_packet event handler
func Test_recvPacket_boundaryChecks(t *testing.T) {
	tests := []struct {
		name        string
		ctx         *context.Context
		events      []storage.Event
		msg         *storage.Message
		idx         *int
		expectError bool
	}{
		{
			name: "empty events array",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgRecvPacket,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         testsuite.Ptr(0),
			expectError: false, // should handle gracefully
		},
		{
			name: "index at last element",
			ctx:  context.NewContext(),
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
					Type:   "message",
					Data: map[string]any{
						"action": "",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgRecvPacket,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         testsuite.Ptr(0),
			expectError: false,
		},
		{
			name: "minimal events - boundary check only",
			ctx:  context.NewContext(),
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
					Type:   types.EventTypeMessage,
					Data: map[string]any{
						"action": "",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgRecvPacket,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         testsuite.Ptr(0),
			expectError: false, // should not panic on boundary
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleRecvPacket(tt.ctx, tt.events, tt.msg, tt.idx)
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
		idx         *int
		expectError bool
	}{
		{
			name: "empty events after increment",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgAcknowledgement",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgAcknowledgement,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         testsuite.Ptr(0),
			expectError: false,
		},
		{
			name: "index near end of array",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgAcknowledgement",
					},
				},
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data:   map[string]any{},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgAcknowledgement,
				Height: 100,
				Data:   map[string]any{},
			},
			idx:         testsuite.Ptr(0),
			expectError: false,
		},
		{
			name: "loop boundary check - no transfer case",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgAcknowledgement",
					},
				},
				{
					Height: 100,
					Type:   "some_event",
					Data: map[string]any{
						"action": "",
					},
				},
				{
					Height: 100,
					Type:   "some_event",
					Data: map[string]any{
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
			idx:         testsuite.Ptr(0),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleAcknowledgement(tt.ctx, tt.events, tt.msg, tt.idx)
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

// Test_toTheNextAction_boundaryChecks tests toTheNextAction boundary conditions
func Test_toTheNextAction_boundaryChecks(t *testing.T) {
	tests := []struct {
		name           string
		events         []storage.Event
		initialIdx     int
		expectedIdx    int
		shouldNotPanic bool
	}{
		{
			name:           "index at end of array",
			events:         []storage.Event{},
			initialIdx:     0,
			expectedIdx:    0,
			shouldNotPanic: true,
		},
		{
			name: "index one before end",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]any{
						"action": "",
					},
				},
			},
			initialIdx:     0,
			expectedIdx:    0,
			shouldNotPanic: true,
		},
		{
			name: "normal progression",
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
			initialIdx:     0,
			expectedIdx:    1,
			shouldNotPanic: true,
		},
		{
			name: "all empty actions until end",
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
			initialIdx:     0,
			expectedIdx:    1,
			shouldNotPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tt.initialIdx
			require.NotPanics(t, func() {
				toTheNextAction(tt.events, &idx)
			})
			if tt.shouldNotPanic {
				require.Equal(t, tt.expectedIdx, idx)
			}
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
			Data: map[string]any{
				"action": "/ibc.core.channel.v1.MsgRecvPacket",
			},
		},
	}
	msg := &storage.Message{
		Type:   types.MsgRecvPacket,
		Height: 100,
		Data:   map[string]any{},
	}

	t.Run("nil index pointer", func(t *testing.T) {
		err := handleRecvPacket(ctx, events, msg, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nil event index")
	})

	t.Run("nil message pointer", func(t *testing.T) {
		idx := 0
		err := handleRecvPacket(ctx, events, nil, &idx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nil message")
	})

	t.Run("nil context - should not panic", func(t *testing.T) {
		idx := 0
		// This might panic if context methods are called, but handler should check
		require.NotPanics(t, func() {
			_ = handleRecvPacket(nil, events, msg, &idx)
		})
	})
}

// Test_acknowledgement_loopIncrementSafety tests that loop increments don't cause panics
func Test_acknowledgement_loopIncrementSafety(t *testing.T) {
	tests := []struct {
		name   string
		events []storage.Event
		msg    *storage.Message
		idx    int
	}{
		{
			name: "increment at exact boundary",
			events: []storage.Event{
				{
					Height: 100,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgAcknowledgement",
					},
				},
				{
					Height: 100,
					Type:   types.EventTypeMessage,
					Data: map[string]any{
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
			idx: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			idx := tt.idx
			require.NotPanics(t, func() {
				err := handleAcknowledgement(ctx, tt.events, tt.msg, &idx)
				require.NoError(t, err)
			})
		})
	}
}
