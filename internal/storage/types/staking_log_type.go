// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
