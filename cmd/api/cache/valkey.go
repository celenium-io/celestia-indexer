// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"context"
	"time"

	"github.com/pkg/errors"
	valkey "github.com/valkey-io/valkey-go"
)

var _ ICache = (*ValKey)(nil)

type ValKey struct {
	client     valkey.Client
	ttlSeconds int64
}

func NewValKey(url string, ttl time.Duration) (*ValKey, error) {
	opts, err := valkey.ParseURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "parse valkey url")
	}
	client, err := valkey.NewClient(opts)
	if err != nil {
		return nil, errors.Wrap(err, "create valkey client")
	}

	return &ValKey{
		client:     client,
		ttlSeconds: int64(ttl.Seconds()),
	}, nil
}

func (c *ValKey) Get(ctx context.Context, key string) (data string, found bool) {
	val, err := c.client.Do(
		ctx, c.client.B().Get().Key(key).Build(),
	).ToString()
	return val, err == nil
}

func (c *ValKey) Set(ctx context.Context, key string, data string, expirationFunc ExpirationFunc) error {
	expiredAt := c.ttlSeconds
	if expirationFunc != nil {
		expiredAt = int64(expirationFunc().Seconds())
	}

	return c.client.Do(
		ctx,
		c.client.B().Set().Key(key).Value(data).ExSeconds(expiredAt).Build(),
	).Error()
}

func (c *ValKey) Close() error {
	c.client.Close()
	return nil
}
