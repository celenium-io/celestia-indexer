// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package hyperlane

import (
	"time"

	"golang.org/x/time/rate"
)

type ApiOption func(api *Api)

func WithRateLimit(rps int) ApiOption {
	return func(api *Api) {
		api.rateLimit = rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), rps)
	}
}

func WithTimeout(timeout time.Duration) ApiOption {
	return func(api *Api) {
		api.timeout = timeout
	}
}
