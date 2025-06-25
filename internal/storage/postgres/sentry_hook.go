// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"go.opentelemetry.io/otel/trace"

	"github.com/getsentry/sentry-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type SentryHook struct {
	dbName        string
	formatQueries bool
	tracer        trace.Tracer
}

var _ bun.QueryHook = (*SentryHook)(nil)

func NewSentryHook(dbName string, tracer trace.Tracer, formatQueries bool) *SentryHook {
	h := new(SentryHook)
	h.dbName = dbName
	h.formatQueries = formatQueries
	h.tracer = tracer
	return h
}

func (h *SentryHook) Init(db *bun.DB) {}

func (h *SentryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	ctx, _ = h.tracer.Start(ctx, "db", trace.WithSpanKind(trace.SpanKindClient))
	return ctx
}

func (h *SentryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	operation := event.Operation()

	root := sentry.TransactionFromContext(ctx)
	if root == nil {
		return
	}

	span := root.StartChild(operation, sentry.WithOpName("db"))
	defer span.Finish()

	query := h.eventQuery(event)

	span.Description = query
	span.SetData("db.statement", query)
	span.SetTag("db.operation", operation)
	span.SetTag("db.name", h.dbName)
	if event.IQuery != nil {
		if tableName := event.IQuery.GetTableName(); tableName != "" {
			span.SetTag("db.sql.table", tableName)
		}
	}

	span.SetTag("db.system", "postgresql")
	if event.Result != nil {
		if n, _ := event.Result.RowsAffected(); n > 0 {
			span.SetData("db.rows_affected", strconv.Itoa(int(n)))
		}
	}

	switch event.Err {
	case nil, sql.ErrNoRows, sql.ErrTxDone:
		// ignore
	default:
		span.SetTag("exception.message", event.Err.Error())
		span.SetTag("exception.type", fmt.Sprintf("%T", event.Err))
		// span.SetStatus(codes.Error, event.Err.Error())
	}
}

func (h *SentryHook) eventQuery(event *bun.QueryEvent) string {
	const softQueryLimit = 8000
	const hardQueryLimit = 16000

	var query string

	if h.formatQueries && len(event.Query) <= softQueryLimit {
		query = event.Query
	} else {
		query = unformattedQuery(event)
	}

	if len(query) > hardQueryLimit {
		query = query[:hardQueryLimit]
	}

	return query
}

func unformattedQuery(event *bun.QueryEvent) string {
	if event.IQuery != nil {
		if b, err := event.IQuery.AppendQuery(schema.NewNopFormatter(), nil); err == nil {
			return string(b)
		}
	}
	return event.QueryTemplate
}
