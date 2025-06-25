// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import "fmt"

type Level int64

func (l Level) String() string {
	return fmt.Sprintf("%d", l)
}
