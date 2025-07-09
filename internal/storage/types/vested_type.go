// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum VestingType
/*
	ENUM(
		delayed,
		periodic,
		permanent,
		continuous
	)
*/
//go:generate go-enum --marshal --sql --values
type VestingType string
