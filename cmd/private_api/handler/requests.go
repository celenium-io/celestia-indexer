// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import "github.com/labstack/echo/v4"

func bindAndValidate[T any](c echo.Context) (*T, error) {
	req := new(T)
	if err := c.Bind(req); err != nil {
		return req, err
	}
	if err := c.Validate(req); err != nil {
		return req, err
	}
	return req, nil
}
