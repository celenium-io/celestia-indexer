// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package math

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}
