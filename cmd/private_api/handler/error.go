// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
	errInvalidAddress   = errors.New("invalid address")
	errUnknownAddress   = errors.New("unknown address")
	errUnknownNamespace = errors.New("unknown namespace")
	errInvalidApiKey    = errors.New("invalid api key")
	errCancelRequest    = "pq: canceling statement due to user request"
)

type NoRows interface {
	IsNoRows(err error) bool
}

type Error struct {
	Message string `json:"message"`
}

func handleError(c echo.Context, err error, noRows NoRows) error {
	if err == nil {
		return nil
	}
	if err.Error() == errCancelRequest {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return c.JSON(http.StatusRequestTimeout, Error{
			Message: "timeout",
		})
	}
	if errors.Is(err, context.Canceled) {
		return c.JSON(http.StatusBadGateway, Error{
			Message: err.Error(),
		})
	}
	if noRows.IsNoRows(err) {
		return c.NoContent(http.StatusNoContent)
	}
	if errors.Is(err, errInvalidAddress) || errors.Is(err, errUnknownAddress) {
		return badRequestError(c, err)
	}
	return internalServerError(c, err)
}

func internalServerError(c echo.Context, err error) error {
	if hub := sentryecho.GetHubFromContext(c); hub != nil {
		hub.CaptureMessage(err.Error())
	}
	return c.JSON(http.StatusInternalServerError, Error{
		Message: err.Error(),
	})
}

func badRequestError(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, Error{
		Message: err.Error(),
	})
}

func success(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
