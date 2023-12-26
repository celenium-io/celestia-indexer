// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

const (
	sizeColumn       = "size"
	timeColumn       = "time"
	pfbCountColumn   = "pfb_count"
	blobsCountColumn = "blobs_count"
)

func limitScope(q *bun.SelectQuery, limit int) *bun.SelectQuery {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return q.Limit(limit)
}

func sortScope(q *bun.SelectQuery, field string, sort sdk.SortOrder) *bun.SelectQuery {
	if sort != sdk.SortOrderAsc && sort != sdk.SortOrderDesc {
		sort = sdk.SortOrderAsc
	}
	return q.OrderExpr("? ?", bun.Ident(field), bun.Safe(sort))
}

func timeframeScope(q *bun.SelectQuery, tf storage.Timeframe) (*bun.SelectQuery, error) {
	switch tf {
	case storage.TimeframeHour:
		return q.ColumnExpr("time_bucket('1 hour', time) as bucket"), nil
	case storage.TimeframeDay:
		return q.ColumnExpr("time_bucket('1 day', time) as bucket"), nil
	case storage.TimeframeWeek:
		return q.ColumnExpr("time_bucket('1 week', time) as bucket"), nil
	case storage.TimeframeMonth:
		return q.ColumnExpr("time_bucket('1 month', time) as bucket"), nil
	case storage.TimeframeYear:
		return q.ColumnExpr("time_bucket('1 year', time) as bucket"), nil
	default:
		return nil, errors.Errorf("unexpected timeframe %s", tf)
	}
}

func txFilter(query *bun.SelectQuery, fltrs storage.TxFilter) *bun.SelectQuery {
	query = limitScope(query, fltrs.Limit)
	query = sortScope(query, "id", fltrs.Sort)

	if !fltrs.MessageTypes.Empty() {
		query = query.Where("bit_count(message_types & ?::bit(74)) > 0", fltrs.MessageTypes)
	}

	if !fltrs.ExcludedMessageTypes.Empty() {
		query = query.Where("bit_count(message_types & ~(?::bit(74))) > 0", fltrs.ExcludedMessageTypes)
	}

	if len(fltrs.Status) > 0 {
		query = query.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			for i := range fltrs.Status {
				sq = sq.WhereOr("status = ?", fltrs.Status[i])
			}
			return sq
		})
	}
	if fltrs.Height > 0 {
		query = query.Where("height = ?", fltrs.Height)
	}

	if !fltrs.TimeFrom.IsZero() {
		query = query.Where("time >= ?", fltrs.TimeFrom)
	}
	if !fltrs.TimeTo.IsZero() {
		query = query.Where("time < ?", fltrs.TimeTo)
	}
	if fltrs.WithMessages {
		query = query.Relation("Messages")
	}
	return query
}

func addressListFilter(query *bun.SelectQuery, fltrs storage.AddressListFilter) *bun.SelectQuery {
	query = limitScope(query, fltrs.Limit)
	query = sortScope(query, "id", fltrs.Sort)
	return query
}

func messagesFilter(query *bun.SelectQuery, fltrs storage.MessageListWithTxFilters) *bun.SelectQuery {
	query = limitScope(query, fltrs.Limit)

	if len(fltrs.MessageTypes) > 0 {
		query = query.Where("type IN (?)", bun.In(fltrs.MessageTypes))
	}
	if len(fltrs.ExcludedMessageTypes) > 0 {
		query = query.Where("type NOT IN (?)", bun.In(fltrs.ExcludedMessageTypes))
	}
	if fltrs.Height > 0 {
		query = query.Where("message.height = ?", fltrs.Height)
	}

	return query
}
