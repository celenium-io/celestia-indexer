// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum ModuleName
/*
	ENUM(
		auth,
		blob,
		crisis,
		distribution,
		indexer,
		gov,
		slashing,
		staking,
		consensus,
		baseapp,
		icahost
	)
*/
//go:generate go-enum --marshal --sql --values
type ModuleName string
