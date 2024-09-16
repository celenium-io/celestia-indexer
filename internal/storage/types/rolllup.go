// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// swagger:enum RollupCategory
/*
	ENUM(
		finance,
		gaming,
		nft
	)
*/
//go:generate go-enum --marshal --sql --values --names
type RollupCategory string

// swagger:enum RollupType
/*
	ENUM(
		sovereign,
		settled
	)
*/
//go:generate go-enum --marshal --sql --values --names
type RollupType string
