// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

/*
	ENUM(
		uncategorized,
		finance,
		gaming,
		nft,
		social
	)
*/
//go:generate go-enum --marshal --sql --values --names
type RollupCategory string

// swagger:enum RollupType
/*
	ENUM(
		sovereign,
		settled,
		other
	)
*/
//go:generate go-enum --marshal --sql --values --names
type RollupType string
