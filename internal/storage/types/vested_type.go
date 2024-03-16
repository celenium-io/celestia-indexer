// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
