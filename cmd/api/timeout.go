// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

func RequestTimeout(timeout time.Duration, skipper func(echo.Context) bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			timeoutCtx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			c.SetRequest(c.Request().WithContext(timeoutCtx))
			defer cancel()
			return next(c)
		}
	}
}
