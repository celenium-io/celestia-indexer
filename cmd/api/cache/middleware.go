// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package cache

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type CacheMiddleware struct {
	cache          ICache
	skipper        middleware.Skipper
	expirationFunc ExpirationFunc
}

func Middleware(cache ICache, skipper middleware.Skipper, expirationFunc ExpirationFunc) echo.MiddlewareFunc {
	mdlwr := CacheMiddleware{
		cache:          cache,
		skipper:        skipper,
		expirationFunc: expirationFunc,
	}
	return mdlwr.Handler
}

func (m *CacheMiddleware) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if m.cache == nil {
			return next(c)
		}
		if m.skipper != nil {
			if m.skipper(c) {
				return next(c)
			}
		}
		path := c.Request().URL.String()
		key := strings.ReplaceAll(strings.TrimPrefix(path, "/"), "/", ":")
		if data, ok := m.cache.Get(c.Request().Context(), key); ok {
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
		return m.cacheResult(c.Request().Context(), key, recorder)
	}
}

func (m *CacheMiddleware) cacheResult(ctx context.Context, key string, r *ResponseRecorder) error {
	result := r.Result()
	if !m.isStatusCacheable(result) {
		return nil
	}

	data, err := result.Encode()
	if err != nil {
		return errors.Wrap(err, "unable to read recorded response")
	}

	return m.cache.Set(ctx, key, data, m.expirationFunc)
}

func (m *CacheMiddleware) isStatusCacheable(e *CacheEntry) bool {
	return e.StatusCode == http.StatusOK
}
