// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum IbcChannelStatus
/*
	ENUM(
		initialization,
		opened,
		closed
	)
*/
//go:generate go-enum --marshal --sql --values --names
type IbcChannelStatus string
