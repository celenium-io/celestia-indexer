// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/stretchr/testify/require"
)

// event returns a storage.Event carrying data[stopKey] = value, or no
// stopKey data at all when stopKey is empty. label is stashed under
// "label" purely so test failures can identify which event was involved.
func event(label string, stopKey, value string) storage.Event {
	data := map[string]string{"label": label}
	if stopKey != "" {
		data[stopKey] = value
	}
	return storage.Event{Data: data}
}

func labels(events []storage.Event) []string {
	out := make([]string, len(events))
	for i, e := range events {
		out[i] = e.Data["label"]
	}
	return out
}

func collect(seq func(yield func(storage.Event) bool)) []storage.Event {
	var out []storage.Event
	for e := range seq {
		out = append(out, e)
	}
	return out
}

func TestCursor_PeekNext_EmptySlice(t *testing.T) {
	c := NewCursor(nil)

	_, ok := c.Peek()
	require.False(t, ok)

	_, ok = c.Next()
	require.False(t, ok)
}

func TestCursor_Peek_DoesNotAdvance(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
		event("b", "", ""),
	})

	first, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "a", first.Data["label"])

	// Peeking again must return the same event.
	again, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "a", again.Data["label"])

	next, ok := c.Next()
	require.True(t, ok)
	require.Equal(t, "a", next.Data["label"])

	second, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "b", second.Data["label"])
}

func TestCursor_Next_ExhaustsThenReportsFalse(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
		event("b", "", ""),
	})

	first, ok := c.Next()
	require.True(t, ok)
	require.Equal(t, "a", first.Data["label"])

	second, ok := c.Next()
	require.True(t, ok)
	require.Equal(t, "b", second.Data["label"])

	// Once exhausted, both Next and Peek keep reporting false without
	// panicking, and the position does not run past the end.
	_, ok = c.Next()
	require.False(t, ok)
	_, ok = c.Peek()
	require.False(t, ok)
	_, ok = c.Next()
	require.False(t, ok)
}

func TestCursor_MsgEvents_StopsBeforeBoundaryWithoutConsumingIt(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
		event("b", "", ""),
		event("boundary", "action", "/some.Msg"),
		event("d", "", ""),
	})

	got := collect(c.MsgEvents("action"))
	require.Equal(t, []string{"a", "b"}, labels(got))

	// The boundary event must still be there, unconsumed.
	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "boundary", next.Data["label"])
}

func TestCursor_MsgEvents_NoBoundary_YieldsEverythingAndEndsAtEOF(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
		event("b", "", ""),
		event("c", "", ""),
	})

	got := collect(c.MsgEvents("action"))
	require.Equal(t, []string{"a", "b", "c"}, labels(got))

	_, ok := c.Peek()
	require.False(t, ok, "cursor must be at the end when no boundary event exists")
}

func TestCursor_MsgEvents_EmptySlice_YieldsNothing(t *testing.T) {
	c := NewCursor(nil)

	got := collect(c.MsgEvents("action"))
	require.Empty(t, got)
}

func TestCursor_MsgEvents_NoOpWhenAlreadyAtBoundary(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("boundary", "action", "/some.Msg"),
		event("a", "", ""),
	})

	got := collect(c.MsgEvents("action"))
	require.Empty(t, got, "MsgEvents must yield nothing when the current event is itself a boundary")

	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "boundary", next.Data["label"], "position must be unchanged")
}

func TestCursor_MsgEvents_EarlyBreakLeavesRemainingEventsUnconsumed(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
		event("b", "", ""),
		event("boundary", "action", "/some.Msg"),
	})

	var got []storage.Event
	for e := range c.MsgEvents("action") {
		got = append(got, e)
		if len(got) == 1 {
			break
		}
	}
	require.Equal(t, []string{"a"}, labels(got))

	// "b" was not consumed by the aborted iteration.
	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "b", next.Data["label"])
}

func TestCursor_SkipToNext_LandsOnBoundaryWithoutConsumingIt(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
		event("b", "", ""),
		event("boundary", "module", "staking"),
		event("d", "", ""),
	})

	c.SkipToNext("module")

	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "boundary", next.Data["label"])

	// The boundary event is still there to be consumed by the caller.
	consumed, ok := c.Next()
	require.True(t, ok)
	require.Equal(t, "boundary", consumed.Data["label"])
}

func TestCursor_SkipToNext_NoBoundary_EndsAtEOF(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
		event("b", "", ""),
	})

	c.SkipToNext("action")

	_, ok := c.Peek()
	require.False(t, ok)
}

func TestCursor_SkipToNext_NoOpWhenAlreadyAtBoundary(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("boundary", "action", "/some.Msg"),
		event("a", "", ""),
	})

	c.SkipToNext("action")

	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "boundary", next.Data["label"], "position must not move past an already-current boundary")
}

func TestCursor_SkipToNext_EmptySlice_NoOp(t *testing.T) {
	c := NewCursor(nil)

	c.SkipToNext("action")

	_, ok := c.Peek()
	require.False(t, ok)
}

func TestCursor_DifferentStopKeysAreIndependent(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "module", "staking"), // carries "module" but not "action"
		event("b", "", ""),
		event("boundary", "action", "/some.Msg"),
	})

	// Asking for the "action" boundary must not stop early on an event
	// that only carries "module".
	got := collect(c.MsgEvents("action"))
	require.Equal(t, []string{"a", "b"}, labels(got))
}

func TestCursor_Skip(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
		event("b", "", ""),
		event("c", "", ""),
	})

	c.Skip(2)

	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "c", next.Data["label"])
}

func TestCursor_Skip_Zero_NoOp(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
	})

	c.Skip(0)

	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "a", next.Data["label"])
}

func TestCursor_Skip_PastEnd_ClampsToEnd(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
		event("b", "", ""),
	})

	c.Skip(10)

	_, ok := c.Peek()
	require.False(t, ok)

	// Clamped, not overflowed: a further Skip must not panic or move
	// past the end.
	require.NotPanics(t, func() { c.Skip(5) })
	_, ok = c.Peek()
	require.False(t, ok)
}

func TestCursor_Skip_EmptySlice_NoOp(t *testing.T) {
	c := NewCursor(nil)

	require.NotPanics(t, func() { c.Skip(3) })
	_, ok := c.Peek()
	require.False(t, ok)
}

func TestCursor_Remaining_DoesNotConsume(t *testing.T) {
	c := NewCursor([]storage.Event{
		event("a", "", ""),
		event("b", "", ""),
		event("c", "", ""),
	})
	c.Next()

	require.Equal(t, []string{"b", "c"}, labels(c.Remaining()))

	// Remaining must be read-only: position is unchanged by the call.
	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "b", next.Data["label"])
}

func TestCursor_Remaining_EmptyAtEnd(t *testing.T) {
	c := NewCursor([]storage.Event{event("a", "", "")})
	c.Next()

	require.Empty(t, c.Remaining())
}
