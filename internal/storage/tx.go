// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
	IdAndTimeByHash(ctx context.Context, hash []byte) (uint64, time.Time, error)
	Filter(ctx context.Context, fltrs TxFilter) ([]Tx, error)
	ByIdWithRelations(ctx context.Context, id uint64) (Tx, error)
	ByAddress(ctx context.Context, addressId uint64, fltrs TxFilter) ([]Tx, error)
	Genesis(ctx context.Context, limit, offset int, sortOrder storage.SortOrder) ([]Tx, error)
	Gas(ctx context.Context, height pkgTypes.Level, ts time.Time) ([]Gas, error)
}

type Gas struct {
	GasWanted int64           `bun:"gas_wanted"`
	GasUsed   int64           `bun:"gas_used"`
	Fee       decimal.Decimal `bun:"fee"`
	GasPrice  decimal.Decimal `bun:"gas_price"`
}

type ByGasPrice []Gas

func (gp ByGasPrice) Len() int           { return len(gp) }
func (gp ByGasPrice) Less(i, j int) bool { return gp[j].GasPrice.GreaterThan(gp[i].GasPrice) }
func (gp ByGasPrice) Swap(i, j int)      { gp[i], gp[j] = gp[j], gp[i] }

type TxFilter struct {
	Limit                int
	Offset               int
	Sort                 storage.SortOrder
	Status               []string
	MessageTypes         types.MsgTypeBits
	ExcludedMessageTypes types.MsgTypeBits
	Height               *uint64
	TimeFrom             time.Time
	TimeTo               time.Time
	WithMessages         bool
}

func (filter *TxFilter) IsEmpty() bool {
	return len(filter.Status) == 0 &&
		filter.MessageTypes.Empty() &&
		filter.ExcludedMessageTypes.Empty() &&
		filter.Height == nil &&
		filter.TimeFrom.IsZero() &&
		filter.TimeTo.IsZero()
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

	Error        string            `bun:"error,type:text"             comment:"Error string if failed"`
	Codespace    string            `bun:"codespace,type:text"         comment:"Codespace"                                    stats:"filterable"`
	Hash         []byte            `bun:"hash"                        comment:"Transaction hash"`
	Memo         string            `bun:"memo,type:text"              comment:"Note or comment to send with the transaction"`
	MessageTypes types.MsgTypeBits `bun:"message_types,type:bit(111)" comment:"Bit mask with containing messages"            stats:"filterable"`

	Messages []Message `bun:"rel:has-many,join:id=tx_id"`
	Events   []Event   `bun:"rel:has-many"`

	Signers    []Address `bun:"-"`
	BlobsSize  int64     `bun:"-"`
	BytesSize  int64     `bun:"-"`
	BlobsCount int       `bun:"-"`
}

// TableName -
func (Tx) TableName() string {
	return "tx"
}
