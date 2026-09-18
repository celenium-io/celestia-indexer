// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/stretchr/testify/require"
)

// TestMessageTypesColumnWidth keeps the message_types columns in sync with the
// bit mask. The bun tag cannot reference MsgTypeBitsCount, and a mismatch is
// only noticed at runtime: Postgres rejects a bit(n) literal of another width.
func TestMessageTypesColumnWidth(t *testing.T) {
	expected := fmt.Sprintf("type:bit(%d)", types.MsgTypeBitsCount)

	for _, model := range []any{Block{}, Tx{}} {
		field, ok := reflect.TypeOf(model).FieldByName("MessageTypes")
		require.Truef(t, ok, "%T has no MessageTypes field", model)
		require.Containsf(t, field.Tag.Get("bun"), expected, "%T declares a stale mask width", model)
	}
}
