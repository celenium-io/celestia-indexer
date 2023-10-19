// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import "fmt"

type Level int64

func (l Level) String() string {
	return fmt.Sprintf("%d", l)
}
