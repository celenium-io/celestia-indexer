// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// swagger:enum HLTokenType
/*
	ENUM(
		synthetic,
		collateral
	)
*/
//go:generate go-enum --marshal --sql --values --names
type HLTokenType string

// swagger:enum HLTransferType
/*
	ENUM(
		send,
		receive
	)
*/
//go:generate go-enum --marshal --sql --values --names
type HLTransferType string
