// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ITx interface {
	storage.Table[*Tx]

	ByHash(ctx context.Context, hash []byte) (Tx, error)
	Filter(ctx context.Context, fltrs TxFilter) ([]Tx, error)
	ByIdWithRelations(ctx context.Context, id uint64) (Tx, error)
	ByAddress(ctx context.Context, addressId uint64, fltrs TxFilter) ([]Tx, error)
	Genesis(ctx context.Context, limit, offset int, sortOrder storage.SortOrder) ([]Tx, error)
}

type TxFilter struct {
	Limit        int
	Offset       int
	Sort         storage.SortOrder
	Status       []string
	MessageTypes types.MsgTypeBits
	Height       uint64
	TimeFrom     time.Time
	TimeTo       time.Time
	WithMessages bool
}

// Tx -
type Tx struct {
	bun.BaseModel `bun:"tx" comment:"Table with celestia transactions."`

	Id            uint64          `bun:"id,autoincrement,pk,notnull" comment:"Unique internal id"`
	Height        pkgTypes.Level  `bun:",notnull"                    comment:"The number (height) of this block"                 stats:"func:min max,filterable"`
	Time          time.Time       `bun:"time,pk,notnull"             comment:"The time of block"                                 stats:"func:min max,filterable"`
	Position      int64           `bun:"position"                    comment:"Position in block"`
	GasWanted     int64           `bun:"gas_wanted"                  comment:"Gas wanted"                                        stats:"func:min max sum avg"`
	GasUsed       int64           `bun:"gas_used"                    comment:"Gas used"                                          stats:"func:min max sum avg"`
	TimeoutHeight uint64          `bun:"timeout_height"              comment:"Block height until which the transaction is valid" stats:"func:min max avg"`
	EventsCount   int64           `bun:"events_count"                comment:"Events count in transaction"                       stats:"func:min max sum avg"`
	MessagesCount int64           `bun:"messages_count"              comment:"Messages count in transaction"                     stats:"func:min max sum avg"`
	Fee           decimal.Decimal `bun:"fee,type:numeric"            comment:"Paid fee"                                          stats:"func:min max sum avg"`
	Status        types.Status    `bun:"status,type:status"          comment:"Transaction status"                                stats:"filterable"`

	Error        string            `bun:"error,type:text"            comment:"Error string if failed"`
	Codespace    string            `bun:"codespace,type:text"        comment:"Codespace"                                    stats:"filterable"`
	Hash         []byte            `bun:"hash"                       comment:"Transaction hash"`
	Memo         string            `bun:"memo,type:text"             comment:"Note or comment to send with the transaction"`
	MessageTypes types.MsgTypeBits `bun:"message_types,type:bit(73)" comment:"Bit mask with containing messages"            stats:"filterable"`

	Messages []Message `bun:"rel:has-many,join:id=tx_id"`
	Events   []Event   `bun:"rel:has-many"`
	Signers  []Address `bun:"m2m:signer,join:Tx=Address"`

	BlobsSize int64 `bun:"-"`
}

// TableName -
func (Tx) TableName() string {
	return "tx"
}
