// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum UpgradeStatus
/*
	ENUM(
		processing,
		applied,
		waiting_upgrade
	)
*/
//go:generate go-enum --marshal --sql --values --names
type UpgradeStatus string
