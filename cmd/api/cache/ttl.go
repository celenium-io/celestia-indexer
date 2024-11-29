// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"time"

	"github.com/cespare/xxhash"
	"github.com/elastic/go-freelru"
)

type TTLCache struct {
	db         *freelru.LRU[string, []byte]
	expiration time.Duration
}

func hashStringXXHASH(s string) uint32 {
	return uint32(xxhash.Sum64String(s))
}

func NewTTLCache(expiration time.Duration) (*TTLCache, error) {
	lru, err := freelru.New[string, []byte](8192*32, hashStringXXHASH)
	if err != nil {
		return nil, err
	}
	lru.SetLifetime(expiration)

	return &TTLCache{
		db:         lru,
		expiration: expiration,
	}, nil
}

func (c *TTLCache) Get(key string) (data []byte, found bool) {
	return c.db.Get(key)
}

func (c *TTLCache) Set(key string, data []byte) {
	c.db.Add(key, data)
}

func (c *TTLCache) Clear() {
	c.db.Purge()
}
