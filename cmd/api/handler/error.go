// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"net/http"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var (
	errInvalidHashLength = errors.New("invalid hash: should be 32 bytes length")
	errInvalidAddress    = errors.New("invalid address")
)

type NoRows interface {
	IsNoRows(err error) bool
}

type Error struct {
	Message string `json:"message"`
}

func badRequestError(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, Error{
		Message: err.Error(),
	})
}

func internalServerError(c echo.Context, err error) error {
	if hub := sentryecho.GetHubFromContext(c); hub != nil {
		if !errors.Is(err, context.Canceled) {
			hub.CaptureMessage(err.Error())
		}
	}
	return c.JSON(http.StatusInternalServerError, Error{
		Message: err.Error(),
	})
}

func handleError(c echo.Context, err error, noRows NoRows) error {
	if err == nil {
		return nil
	}
	if noRows.IsNoRows(err) {
		return c.NoContent(http.StatusNoContent)
	}
	if errors.Is(err, errInvalidAddress) {
		return badRequestError(c, err)
	}
	return internalServerError(c, err)
}

func success(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
