package storage

import (
	"context"
	"time"

	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlockStats interface {
	ByHeight(ctx context.Context, height uint64) (BlockStats, error)
}

type BlockStats struct {
	bun.BaseModel `bun:"table:block_stats" comment:"Table with celestia block stats."`

	Id     uint64         `bun:",pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height pkgTypes.Level `bun:"height"                    comment:"The number (height) of this block" stats:"func:min max,filterable"`
	Time   time.Time      `bun:"time,pk,notnull"           comment:"The time of block"                 stats:"func:min max,filterable"`

	TxCount       uint64          `bun:"tx_count"         comment:"Count of transactions in block"            stats:"func:min max sum avg"`
	EventsCount   uint64          `bun:"events_count"     comment:"Count of events in begin and end of block" stats:"func:min max sum avg"`
	BlobsSize     uint64          `bun:"blobs_size"       comment:"Summary blocks size from pay for blob"     stats:"func:min max sum avg"`
	SupplyChange  decimal.Decimal `bun:",type:numeric"    comment:"Change of total supply in the block"       stats:"func:min max sum avg"`
	InflationRate decimal.Decimal `bun:",type:numeric"    comment:"Inflation rate"                            stats:"func:min max avg"`
	Fee           decimal.Decimal `bun:"fee,type:numeric" comment:"Summary block fee"                         stats:"func:min max sum avg"`
}

func (BlockStats) TableName() string {
	return "block_stats"
}
