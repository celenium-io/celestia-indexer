// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
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

type MessageListWithTxFilters struct {
	Height               pkgTypes.Level
	Limit                int
	Offset               int
	ExcludedMessageTypes []string
	MessageTypes         []string
}

type AddressMsgsFilter struct {
	Limit        int
	Offset       int
	Sort         storage.SortOrder
	MessageTypes []string
}

type MessageWithTx struct {
	bun.BaseModel `bun:"message,alias:message"`

	Message
	Tx *Tx `bun:"rel:belongs-to"`
}

type AddressMessageWithTx struct {
	bun.BaseModel `bun:"message,alias:message"`

	MsgAddress
	Msg *Message `bun:"rel:belongs-to,join:msg_id=id"`
	Tx  *Tx      `bun:"rel:belongs-to,join:msg__tx_id=id"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IMessage interface {
	storage.Table[*Message]

	ByTxId(ctx context.Context, txId uint64, limit, offset int) ([]Message, error)
	ListWithTx(ctx context.Context, filters MessageListWithTxFilters) ([]MessageWithTx, error)
	ByAddress(ctx context.Context, id uint64, filters AddressMsgsFilter) ([]AddressMessageWithTx, error)
}

// Message -
type Message struct {
	bun.BaseModel `bun:"message" comment:"Table with celestia messages."`

	Id       uint64         `bun:"id,pk,notnull,autoincrement" comment:"Unique internal id"`
	Height   pkgTypes.Level `bun:",notnull"                    comment:"The number (height) of this block" stats:"func:min max,filterable"`
	Time     time.Time      `bun:"time,pk,notnull"             comment:"The time of block"                 stats:"func:min max,filterable"`
	Position int64          `bun:"position"                    comment:"Position in transaction"`
	Type     types.MsgType  `bun:",type:msg_type"              comment:"Message type"                      stats:"filterable"`
	TxId     uint64         `bun:"tx_id"                       comment:"Parent transaction id"`
	Data     map[string]any `bun:"data,type:jsonb,nullzero"    comment:"Message data"`

	Namespace []Namespace       `bun:"m2m:namespace_message,join:Message=Namespace"`
	Validator *Validator        `bun:"rel:belongs-to"`
	Addresses []AddressWithType `bun:"-"`
	BlobLogs  []*BlobLog        `bun:"-"`
}

// TableName -
func (Message) TableName() string {
	return "message"
}
