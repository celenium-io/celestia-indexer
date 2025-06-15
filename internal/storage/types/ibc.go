// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// swagger:enum IbcChannelStatus
/*
	ENUM(
		initialization,
		opened,
		closed
	)
*/
//go:generate go-enum --marshal --sql --values --names
type IbcChannelStatus string
