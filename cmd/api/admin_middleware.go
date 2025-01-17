// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"net/http"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

var accessDeniedErr = echo.Map{
	"error": "access denied",
}

func AdminMiddleware() echo.MiddlewareFunc {
	return checkOnAdminPermission
}

func checkOnAdminPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		val := ctx.Get(handler.ApiKeyName)
		apiKey, ok := val.(storage.ApiKey)
		if !ok {
			return ctx.JSON(http.StatusForbidden, accessDeniedErr)
		}
		if !apiKey.Admin {
			return ctx.JSON(http.StatusForbidden, accessDeniedErr)
		}
		return next(ctx)
	}
}
