// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decoder

import "strings"

// CleanString drops NUL bytes: Postgres text and bun reject them, so a chain string with U+0000 would halt indexing.
func CleanString(s string) string {
	return strings.ReplaceAll(s, "\x00", "")
}
