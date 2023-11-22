// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type Cache struct {
	maxEntitiesCount int

	m     map[string][]byte
	queue []string
	mx    *sync.RWMutex
}

type Config struct {
	MaxEntitiesCount int
}

func NewCache(cfg Config) *Cache {
	return &Cache{
		maxEntitiesCount: cfg.MaxEntitiesCount,
		m:                make(map[string][]byte),
		queue:            make([]string, cfg.MaxEntitiesCount),
		mx:               new(sync.RWMutex),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mx.RLock()
	data, ok := c.m[key]
	c.mx.RUnlock()
	return data, ok
}

func (c *Cache) Set(key string, data []byte) {
	c.mx.Lock()
	queueIdx := len(c.m)
	c.m[key] = data
	if queueIdx == c.maxEntitiesCount {
		keyForRemove := c.queue[queueIdx-1]
		c.queue = append([]string{key}, c.queue[:queueIdx-1]...)
		delete(c.m, keyForRemove)
	} else {
		c.queue[queueIdx] = key
	}
	c.mx.Unlock()
}

func (c *Cache) Clear() {
	c.mx.Lock()
	for key := range c.m {
		delete(c.m, key)
	}
	c.queue = make([]string, c.maxEntitiesCount)
	c.mx.Unlock()
}

type CacheMiddleware struct {
	cache   *Cache
	skipper middleware.Skipper
}

func Middleware(cache *Cache, skipper middleware.Skipper) echo.MiddlewareFunc {
	mdlwr := CacheMiddleware{
		cache:   cache,
		skipper: skipper,
	}
	return mdlwr.Handler
}

func (m *CacheMiddleware) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if m.skipper != nil {
			if m.skipper(c) {
				return next(c)
			}
		}
		path := c.Request().URL.String()

		if data, ok := m.cache.Get(path); ok {
			entry := new(CacheEntry)
			if err := entry.Decode(data); err != nil {
				return err
			}
			return entry.Replay(c.Response())
		}

		recorder := NewResponseRecorder(c.Response().Writer)
		c.Response().Writer = recorder

		if err := next(c); err != nil {
			return err
		}
		return m.cacheResult(path, recorder)
	}
}

func (m *CacheMiddleware) cacheResult(key string, r *ResponseRecorder) error {
	result := r.Result()
	if !m.isStatusCacheable(result) {
		return nil
	}

	data, err := result.Encode()
	if err != nil {
		return errors.Wrap(err, "unable to read recorded response")
	}

	m.cache.Set(key, data)
	return nil
}

func (m *CacheMiddleware) isStatusCacheable(e *CacheEntry) bool {
	return e.StatusCode == http.StatusOK || e.StatusCode == http.StatusNoContent
}
