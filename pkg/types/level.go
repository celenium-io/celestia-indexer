// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import "fmt"

type Level int64

func (l Level) String() string {
	return fmt.Sprintf("%d", l)
}
