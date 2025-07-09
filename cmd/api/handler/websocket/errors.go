// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

import "errors"

var (
	ErrUnknownMethod     = errors.New("unknown method")
	ErrUnknownChannel    = errors.New("unknown channel")
	ErrUnavailableFilter = errors.New("unknown filter value")
)
