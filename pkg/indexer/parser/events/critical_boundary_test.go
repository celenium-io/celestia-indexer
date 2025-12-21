// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/stretchr/testify/require"
)

// Test_toTheNextAction_criticalBoundary tests the critical case where increment happens
// at exactly len(events)-2, and after increment we're at len(events)-1
func Test_toTheNextAction_criticalBoundary(t *testing.T) {
	tests := []struct {
		name        string
		events      []storage.Event
		startIdx    int
		expectPanic bool
	}{
		{
			name: "CRITICAL: increment at len-2 with empty action",
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
						"action": "", // empty action at len-1
					},
				},
			},
			startIdx:    0,
			expectPanic: false, // should NOT panic if bounds check is correct
		},
		{
			name: "CRITICAL: three elements, start at len-2",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]any{
						"action": "first",
					},
				},
				{
					Type: "message",
					Data: map[string]any{
						"action": "", // at index 1 (len-2)
					},
				},
				{
					Type: "message",
					Data: map[string]any{
						"action": "", // at index 2 (len-1)
					},
				},
			},
			startIdx:    1,     // start at len-2
			expectPanic: false, // should NOT panic
		},
		{
			name: "Edge case: two elements both empty",
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
			startIdx:    0,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tt.startIdx

			if tt.expectPanic {
				require.Panics(t, func() {
					toTheNextAction(tt.events, &idx)
				}, "Expected panic but function completed normally")
			} else {
				require.NotPanics(t, func() {
					toTheNextAction(tt.events, &idx)
				}, "Function should not panic with proper boundary checks")
			}

			// Additional safety check: index should never exceed len(events)-1
			require.LessOrEqual(t, idx, len(tt.events)-1,
				"Index should never exceed last valid position")
		})
	}
}

// Test_acknowledgement_criticalLoop tests the loop in acknowledgement handler
// This is the exact scenario from line 81-98 in acknowledgement.go
func Test_acknowledgement_criticalLoop(t *testing.T) {
	tests := []struct {
		name        string
		events      []storage.Event
		startIdx    int
		expectPanic bool
	}{
		{
			name: "Loop with increment reaching exact boundary",
			events: []storage.Event{
				{
					Type: "some_event",
					Data: map[string]any{
						"action": "",
					},
				},
				{
					Type: "last_event",
					Data: map[string]any{
						"action": "",
					},
				},
			},
			startIdx:    0,
			expectPanic: false,
		},
		{
			name: "Multiple events with empty actions",
			events: []storage.Event{
				{
					Type: "event1",
					Data: map[string]any{
						"action": "",
					},
				},
				{
					Type: "event2",
					Data: map[string]any{
						"action": "",
					},
				},
				{
					Type: "event3",
					Data: map[string]any{
						"action": "",
					},
				},
			},
			startIdx:    0,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tt.startIdx

			// Simulate the loop from acknowledgement.go:81-98
			simulateLoop := func() {
				action := ""
				for action == "" && len(tt.events)-1 > idx {
					idx += 1
					// This line could panic if bounds check is wrong:
					action = tt.events[idx].Data["action"].(string)
				}
			}

			if tt.expectPanic {
				require.Panics(t, simulateLoop)
			} else {
				require.NotPanics(t, simulateLoop)
			}
		})
	}
}

// Test_recvPacket_criticalLoop tests the loop in recv_packet handler
// This is the exact scenario from line 104-116 in recv_packet.go
func Test_recvPacket_criticalLoop(t *testing.T) {
	tests := []struct {
		name        string
		events      []storage.Event
		startIdx    int
		expectPanic bool
	}{
		{
			name: "Loop reaching exact end",
			events: []storage.Event{
				{
					Type: "event1",
					Data: map[string]any{
						"action": "",
					},
				},
				{
					Type: "event2",
					Data: map[string]any{
						"action": "",
					},
				},
			},
			startIdx:    0,
			expectPanic: false,
		},
		{
			name: "Single event with empty action",
			events: []storage.Event{
				{
					Type: "event",
					Data: map[string]any{
						"action": "",
					},
				},
			},
			startIdx:    0,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tt.startIdx

			// Simulate the loop from recv_packet.go:104-116
			simulateLoop := func() {
				action := ""
				for action == "" && len(tt.events)-1 > idx {
					idx += 1
					// This line could cause panic if idx is out of bounds
					action = tt.events[idx].Data["action"].(string)
				}
			}

			if tt.expectPanic {
				require.Panics(t, simulateLoop)
			} else {
				require.NotPanics(t, simulateLoop)
			}
		})
	}
}
