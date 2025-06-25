// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package cache

import (
	"context"
	"io"
	"time"
)

type ICache interface {
	io.Closer

	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, data string, f ExpirationFunc) error
}

type ExpirationFunc func() time.Duration
