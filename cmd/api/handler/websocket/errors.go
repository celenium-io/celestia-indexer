// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import "errors"

var (
	ErrUnknownMethod     = errors.New("unknown method")
	ErrUnknownChannel    = errors.New("unknown channel")
	ErrUnavailableFilter = errors.New("unknown filter value")
)
