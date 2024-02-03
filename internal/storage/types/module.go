// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
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
		consensus
	)
*/
//go:generate go-enum --marshal --sql --values
type ModuleName string
