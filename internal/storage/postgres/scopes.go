// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"fmt"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

const (
	sizeColumn       = "size"
	timeColumn       = "time"
	pfbCountColumn   = "pfb_count"
	blobsCountColumn = "blobs_count"
	feeColumn        = "fee"
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
	if fltrs.Height != nil {
		query = query.Where("height = ?", *fltrs.Height)
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
	query = query.Offset(fltrs.Offset)

	switch fltrs.SortField {
	case "id", "spendable", "unbonding", "delegated", "last_height":
		query = sortScope(query, fltrs.SortField, fltrs.Sort)
	case "first_height":
		query = sortScope(query, "height", fltrs.Sort)
	default:
		query = sortScope(query, "id", fltrs.Sort)
	}

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

func blobLogFilters(query *bun.SelectQuery, fltrs storage.BlobLogFilters) *bun.SelectQuery {
	if fltrs.Offset > 0 {
		query = query.Offset(fltrs.Offset)
	}

	if !fltrs.From.IsZero() {
		query = query.Where("time >= ?", fltrs.From)
	}
	if !fltrs.To.IsZero() {
		query = query.Where("time < ?", fltrs.To)
	}
	if fltrs.Commitment != "" {
		query = query.Where("commitment = ?", fltrs.Commitment)
	}

	query = limitScope(query, fltrs.Limit)
	return blobLogSort(query, fltrs)
}

func blobLogSort(query *bun.SelectQuery, fltrs storage.BlobLogFilters) *bun.SelectQuery {
	switch fltrs.SortBy {
	case sizeColumn, timeColumn:
		query = sortScope(query, fmt.Sprintf("blob_log.%s", fltrs.SortBy), fltrs.Sort)
	case "":
		query = sortScope(query, "blob_log.id", fltrs.Sort)
	default:
		return query
	}
	return query
}
