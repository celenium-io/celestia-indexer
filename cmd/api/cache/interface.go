// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
