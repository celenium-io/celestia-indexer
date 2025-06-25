// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

type EventFilter struct {
	Limit  int
	Offset int
	Time   time.Time
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IEvent interface {
	storage.Table[*Event]

	ByTxId(ctx context.Context, txId uint64, fltrs EventFilter) ([]Event, error)
	ByBlock(ctx context.Context, height pkgTypes.Level, fltrs EventFilter) ([]Event, error)
}

// Event -
type Event struct {
	bun.BaseModel `bun:"event" comment:"Table with celestia events."`

	Id       uint64          `bun:"id,pk,notnull,autoincrement"      comment:"Unique internal id"`
	Height   pkgTypes.Level  `bun:"height,notnull"                   comment:"The number (height) of this block" stats:"func:min max,filterable"`
	Time     time.Time       `bun:"time,pk,notnull"                  comment:"The time of block"                 stats:"func:min max,filterable"`
	Position int64           `bun:"position"                         comment:"Position in transaction"`
	Type     types.EventType `bun:",type:event_type"                 comment:"Event type"                        stats:"filterable"`
	TxId     *uint64         `bun:"tx_id"                            comment:"Transaction id"`
	Data     map[string]any  `bun:"data,msgpack,type:bytea,nullzero" comment:"Event data"`
}

// TableName -
func (Event) TableName() string {
	return "event"
}
