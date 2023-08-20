package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// ITx -
type ITx interface {
	storage.Table[*Tx]
}

// Tx -
type Tx struct {
	bun.BaseModel `bun:"tx" comment:"Table with celestia transactions." partition:"RANGE(time)"`

	Id            uint64          `bun:"id,type:bigint,pk,notnull" comment:"Unique internal id"`
	Height        uint64          `bun:",notnull"                  comment:"The number (height) of this block"`
	Time          time.Time       `bun:"time,pk,notnull"           comment:"The time of block"`
	Position      uint64          `bun:"position"                  comment:"Position in block"`
	GasWanted     uint64          `bun:"gas_wanted"                comment:"Gas wanted"`
	GasUsed       uint64          `bun:"gas_used"                  comment:"Gas used"`
	TimeoutHeight uint64          `bun:"timeout_height"            comment:"Block height until which the transaction is valid"`
	EventsCount   uint64          `bun:"events_count"              comment:"Events count in transaction"`
	MessagesCount uint64          `bun:"messages_count"            comment:"Messages count in transaction"`
	Fee           decimal.Decimal `bun:",type:numeric"             comment:"Paid fee"`
	Status        Status          `bun:",type:status"              comment:"Transaction status"`
	Error         string          `bun:"error"                     comment:"Error string if failed"`
	Codespace     string          `bun:"codespace"                 comment:"Codespace"`
	Hash          []byte          `bun:"hash"                      comment:"Transaction hash"`
	Memo          string          `bun:"memo"                      comment:"Note or comment to send with the transaction"`

	Messages []Message `bun:"rel:has-many"`
	Events   []Event   `bun:"rel:has-many"`
}

// TableName -
func (Tx) TableName() string {
	return "tx"
}
