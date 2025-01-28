// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package api

import (
	"time"

	"golang.org/x/time/rate"
)

type ApiOption func(api *Api)

func WithTimeout(timeout time.Duration) ApiOption {
	return func(api *Api) {
		api.timeout = timeout
	}
}
func WithRateLimit(rps int) ApiOption {
	return func(api *Api) {
		api.rateLimiter = rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), rps)
	}
}
