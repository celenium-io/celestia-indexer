// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

type GasPrice struct {
	Slow   string `example:"0.1234" format:"string" json:"slow"   swaggertype:"string"`
	Median string `example:"0.1234" format:"string" json:"median" swaggertype:"string"`
	Fast   string `example:"0.1234" format:"string" json:"fast"   swaggertype:"string"`
}
