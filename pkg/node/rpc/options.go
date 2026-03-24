// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

type ApiOption func(api *API)

func WithDisableGzip() ApiOption {
	return func(api *API) {
		api.disableGzip = true
	}
}
