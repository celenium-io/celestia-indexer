// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/stretchr/testify/require"
)

// Test_toTheNextAction_criticalBoundary tests the critical case where the
// cursor sits at or near the last event when SkipToNext("action") is called
// (this used to test the standalone toTheNextAction helper, which
// Cursor.SkipToNext replaced — Peek/Next make out-of-bounds access
// impossible by construction, so these cases can no longer panic).
func Test_toTheNextAction_criticalBoundary(t *testing.T) {
	tests := []struct {
		name     string
		events   []storage.Event
		startIdx int
	}{
		{
			name: "CRITICAL: increment at len-2 with empty action",
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
						"action": "", // empty action at len-1
					},
				},
			},
			startIdx: 0,
		},
		{
			name: "CRITICAL: three elements, start at len-2",
			events: []storage.Event{
				{
					Type: "message",
					Data: map[string]string{
						"action": "first",
					},
				},
				{
					Type: "message",
					Data: map[string]string{
						"action": "", // at index 1 (len-2)
					},
				},
				{
					Type: "message",
					Data: map[string]string{
						"action": "", // at index 2 (len-1)
					},
				},
			},
			startIdx: 1, // start at len-2
		},
		{
			name: "Edge case: two elements both empty",
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
			startIdx: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCursor(tt.events)
			c.Skip(tt.startIdx)

			require.NotPanics(t, func() {
				c.SkipToNext("action")
			}, "Function should not panic with proper boundary checks")
		})
	}
}
