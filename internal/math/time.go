// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package math

import "time"

func TimeFromNano(ts uint64) time.Time {
	return time.Unix(0, int64(ts)).UTC()
}
