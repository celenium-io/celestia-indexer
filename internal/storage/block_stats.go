// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlockStats interface {
	ByHeight(ctx context.Context, height pkgTypes.Level) (BlockStats, error)
}

type BlockStats struct {
	bun.BaseModel `bun:"table:block_stats" comment:"Table with celestia block stats."`

	Id     uint64         `bun:",pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height pkgTypes.Level `bun:"height"                    comment:"The number (height) of this block" stats:"func:min max,filterable"`
	Time   time.Time      `bun:"time,pk,notnull"           comment:"The time of block"                 stats:"func:min max,filterable"`

	TxCount       int64           `bun:"tx_count"         comment:"Count of transactions in block"                          stats:"func:min max sum avg"`
	EventsCount   int64           `bun:"events_count"     comment:"Count of events in begin and end of block"               stats:"func:min max sum avg"`
	BlobsSize     int64           `bun:"blobs_size"       comment:"Summary blocks size from pay for blob"                   stats:"func:min max sum avg"`
	BlockTime     uint64          `bun:"block_time"       comment:"Time in milliseconds between current and previous block" stats:"func:min max sum avg"`
	SupplyChange  decimal.Decimal `bun:",type:numeric"    comment:"Change of total supply in the block"                     stats:"func:min max sum avg"`
	InflationRate decimal.Decimal `bun:",type:numeric"    comment:"Inflation rate"                                          stats:"func:min max avg"`
	Fee           decimal.Decimal `bun:"fee,type:numeric" comment:"Summary block fee"                                       stats:"func:min max sum avg"`

	MessagesCounts map[types.MsgType]int64 `bun:"-"`
}

func (BlockStats) TableName() string {
	return "block_stats"
}
