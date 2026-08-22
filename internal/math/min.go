// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package math

import "cmp"

func Min[T cmp.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}
