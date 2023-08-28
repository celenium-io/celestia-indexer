package storage

import (
	"context"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IMessage interface {
	storage.Table[*Message]

	ByTxId(ctx context.Context, txId uint64) ([]Message, error)
}

// Message -
type Message struct {
	bun.BaseModel `bun:"message" comment:"Table with celestia messages." partition:"RANGE(time)"`

	Id       uint64         `bun:"id,type:bigint,pk,notnull" comment:"Unique internal id"`
	Height   uint64         `bun:",notnull"                  comment:"The number (height) of this block"`
	Time     time.Time      `bun:"time,pk,notnull"           comment:"The time of block"`
	Position uint64         `bun:"position"                  comment:"Position in transaction"`
	Type     types.MsgType  `bun:",type:msg_type"            comment:"Message type"`
	TxId     uint64         `bun:"tx_id"                     comment:"Parent transaction id"`
	Data     map[string]any `bun:"data,type:jsonb"           comment:"Message data"`

	Namespace []Namespace `bun:"m2m:namespace_action,join:Message=Namespace"`
}

// TableName -
func (Message) TableName() string {
	return "message"
}
