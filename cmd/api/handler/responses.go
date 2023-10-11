// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func returnArray[T any](c echo.Context, arr []T) error {
	if arr == nil {
		return c.JSON(http.StatusOK, []any{})
	}

	return c.JSON(http.StatusOK, arr)
}
