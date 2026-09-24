// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum ValidatorStatus
/*
	ENUM(
		active,
		not_active,
		jailed
	)
*/
//go:generate go-enum --marshal --sql --values --names
type ValidatorStatus string
