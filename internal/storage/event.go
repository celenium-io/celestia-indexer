package storage

import (
	"context"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IEvent interface {
	storage.Table[*Event]

	ByTxId(ctx context.Context, txId uint64) ([]Event, error)
	ByBlock(ctx context.Context, height uint64) ([]Event, error)
}

// Event -
type Event struct {
	bun.BaseModel `bun:"event" comment:"Table with celestia events."`

	Id       uint64          `bun:"id,type:bigint,pk,notnull" comment:"Unique internal id"`
	Height   uint64          `bun:",notnull"                  comment:"The number (height) of this block" stats:"func:min max,filterable"`
	Time     time.Time       `bun:"time,pk,notnull"           comment:"The time of block"                 stats:"func:min max,filterable"`
	Position uint64          `bun:"position"                  comment:"Position in transaction"`
	Type     types.EventType `bun:",type:event_type"          comment:"Event type"                        stats:"filterable"`
	TxId     *uint64         `bun:"tx_id"                     comment:"Transaction id"`
	Data     map[string]any  `bun:"data,type:jsonb"           comment:"Event data"`
}

// TableName -
func (Event) TableName() string {
	return "event"
}
