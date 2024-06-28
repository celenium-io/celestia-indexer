package cache

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type CacheMiddleware struct {
	cache   ICache
	skipper middleware.Skipper
}

func Middleware(cache ICache, skipper middleware.Skipper) echo.MiddlewareFunc {
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
