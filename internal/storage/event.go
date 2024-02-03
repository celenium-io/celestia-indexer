// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IEvent interface {
	storage.Table[*Event]

	ByTxId(ctx context.Context, txId uint64, limit, offset int) ([]Event, error)
	ByBlock(ctx context.Context, height pkgTypes.Level, limit, offset int) ([]Event, error)
}

// Event -
type Event struct {
	bun.BaseModel `bun:"event" comment:"Table with celestia events."`

	Id       uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal id"`
	Height   pkgTypes.Level  `bun:"height,notnull"              comment:"The number (height) of this block" stats:"func:min max,filterable"`
	Time     time.Time       `bun:"time,pk,notnull"             comment:"The time of block"                 stats:"func:min max,filterable"`
	Position int64           `bun:"position"                    comment:"Position in transaction"`
	Type     types.EventType `bun:",type:event_type"            comment:"Event type"                        stats:"filterable"`
	TxId     *uint64         `bun:"tx_id"                       comment:"Transaction id"`
	Data     map[string]any  `bun:"data,type:jsonb,nullzero"    comment:"Event data"`
}

// TableName -
func (Event) TableName() string {
	return "event"
}

func (e Event) Columns() []string {
	return []string{
		"height", "time", "position", "type",
		"tx_id", "data",
	}
}

func (e Event) Flat() []any {
	data := []any{
		e.Height, e.Time, e.Position, e.Type, e.TxId, nil,
	}
	if len(e.Data) > 0 {
		raw, err := json.MarshalToString(e.Data)
		if err == nil {
			data[5] = raw
		}
	}
	return data
}
