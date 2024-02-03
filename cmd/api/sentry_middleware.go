// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

const sdkIdentifier = "sentry.go.http"

func SentryMiddleware() echo.MiddlewareFunc {
	return handleSentry
}

func handleSentry(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		hub := sentry.GetHubFromContext(ctx.Request().Context())
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
		}

		if client := hub.Client(); client != nil {
			client.SetSDKIdentifier(sdkIdentifier)
		}

		req := ctx.Request()
		options := []sentry.SpanOption{
			sentry.WithOpName("http.server"),
			sentry.ContinueFromRequest(req),
			sentry.WithTransactionSource(sentry.SourceURL),
		}

		transaction := sentry.StartTransaction(ctx.Request().Context(),
			fmt.Sprintf("%s %s", req.Method, req.URL.Path),
			options...,
		)
		defer func() {
			transaction.Status = sentry.HTTPtoSpanStatus(ctx.Response().Status)
			transaction.Finish()
		}()

		req = req.WithContext(transaction.Context())
		hub.Scope().SetRequest(req)
		hub.Scope().SetUser(sentry.User{
			IPAddress: ctx.RealIP(),
		})
		ctx.Set("sentry", hub)
		transaction.SetTag("method", req.Method)
		transaction.SetTag("user-agent", req.UserAgent())
		transaction.SetTag("ip", ctx.RealIP())

		defer recoverWithSentry(hub, req)

		return next(ctx)
	}
}

func recoverWithSentry(hub *sentry.Hub, r *http.Request) {
	if err := recover(); err != nil {
		eventID := hub.RecoverWithContext(
			context.WithValue(r.Context(), sentry.RequestContextKey, r),
			err,
		)
		if eventID != nil {
			hub.Flush(time.Second * 10)
		}
	}
}
