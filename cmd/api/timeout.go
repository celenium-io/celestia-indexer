package main

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

func RequestTimeout(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			timeoutCtx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			c.SetRequest(c.Request().WithContext(timeoutCtx))
			defer cancel()
			return next(c)
		}
	}
}
