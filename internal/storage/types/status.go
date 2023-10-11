// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// swagger:enum Status
/*
	ENUM(
		success,
		failed
	)
*/
//go:generate go-enum --marshal --sql --values
type Status string
