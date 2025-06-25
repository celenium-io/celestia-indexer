// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum Status
/*
	ENUM(
		success,
		failed
	)
*/
//go:generate go-enum --marshal --sql --values --names
type Status string
