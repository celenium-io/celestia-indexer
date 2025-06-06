// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package math

import "time"

func TimeFromNano(ts uint64) time.Time {
	return time.Unix(0, int64(ts)).UTC()
}
