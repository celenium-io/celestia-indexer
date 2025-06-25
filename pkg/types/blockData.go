// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

type BlockData struct {
	ResultBlock
	ResultBlockResults

	AppVersion uint64
}
