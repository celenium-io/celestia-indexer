// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
