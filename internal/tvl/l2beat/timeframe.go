// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package l2beat

/*
	ENUM(
		7d,
		30d,
		90d,
		180d,
		1y,
		max
	)
*/
//go:generate go-enum --marshal --values --names
type TvlTimeframe string
