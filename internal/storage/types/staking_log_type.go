// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum StakingLogType
/*
	ENUM(
		delegation,
		unbonding,
		rewards,
		commissions
	)
*/
//go:generate go-enum --marshal --sql --values --names
type StakingLogType string
