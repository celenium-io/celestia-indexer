// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum StakingLogType
/*
	ENUM(
		delegation,
		unbonding,
		rewards,
		commissions,
		unbonded,
		slashing
	)
*/
//go:generate go-enum --marshal --sql --values --names
type StakingLogType string
