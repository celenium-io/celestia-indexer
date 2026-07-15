// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"iter"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
)

// Cursor walks a transaction's events, shared across all of that
// transaction's messages. It owns the boundary logic: an event whose Data
// carries a non-empty value for a "stop key" ("action" for top-level
// messages, "module" for IBC-nested ones) marks the start of the next
// message and is never consumed on behalf of the current one.
//
// A Cursor is not safe for concurrent use.
type Cursor struct {
	events []storage.Event
	pos    int
}

// NewCursor returns a Cursor positioned at the first event of events.
func NewCursor(events []storage.Event) *Cursor {
	return &Cursor{
		events: events,
		pos:    0,
	}
}

// Peek reports the event at the current position without consuming it.
// The second result is false once the cursor has reached the end, in
// which case the first result is the zero storage.Event.
func (c *Cursor) Peek() (storage.Event, bool) {
	if c.pos >= len(c.events) {
		return storage.Event{}, false
	}
	return c.events[c.pos], true
}

// Next consumes and returns the event at the current position, advancing
// the cursor by one. The second result is false once the cursor has
// reached the end, in which case the position is left unchanged and the
// first result is the zero storage.Event.
func (c *Cursor) Next() (storage.Event, bool) {
	if c.pos >= len(c.events) {
		return storage.Event{}, false
	}
	event := c.events[c.pos]
	c.pos++
	return event, true
}

// MsgEvents yields events starting at the current position, up to
// (exclusive) the next event whose Data carries a non-empty value for
// stopKey, consuming each event as it is yielded. The boundary event
// itself is never consumed.
//
// If the event at the current position already carries stopKey,
// MsgEvents yields nothing and leaves the position unchanged — call it
// again only after advancing past the boundary (e.g. via Next or
// SkipToNext), otherwise it is a no-op.
func (c *Cursor) MsgEvents(stopKey string) iter.Seq[storage.Event] {
	return func(yield func(storage.Event) bool) {
		for c.pos < len(c.events) {
			if decoder.StringFromMap(c.events[c.pos].Data, stopKey) != "" {
				return
			}
			event := c.events[c.pos]
			c.pos++
			if !yield(event) {
				return
			}
		}
	}
}

// SkipToNext positions the cursor at the next event whose Data carries a
// non-empty value for stopKey (or at the end, if none remains), without
// consuming that event or reporting any of the skipped ones. It is a
// no-op if the event at the current position already carries stopKey.
func (c *Cursor) SkipToNext(stopKey string) {
	for range c.MsgEvents(stopKey) {
	}
}

// Skip advances the cursor by n events without inspecting or reporting
// them, for handlers that unconditionally jump a fixed number of events.
// n must be non-negative. If fewer than n events remain, the cursor is
// left at the end.
func (c *Cursor) Skip(n int) {
	c.pos += n
	if c.pos > len(c.events) {
		c.pos = len(c.events)
	}
}

// Remaining returns the events from the current position to the end,
// without consuming them. It is meant for read-only lookahead (searching
// for an event further ahead while leaving the position where a
// subsequent, separately-computed Skip expects it to be) — prefer Next,
// MsgEvents, or SkipToNext wherever the scan itself should advance the
// cursor.
func (c *Cursor) Remaining() []storage.Event {
	return c.events[c.pos:]
}
