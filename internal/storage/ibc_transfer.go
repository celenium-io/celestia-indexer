// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ListIbcTransferFilters struct {
	Limit      int
	Offset     int
	Sort       sdk.SortOrder
	ReceiverId *uint64
	SenderId   *uint64
	AddressId  *uint64
	ChannelId  string
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IIbcTransfer interface {
	List(ctx context.Context, fltrs ListIbcTransferFilters) ([]IbcTransfer, error)
	Series(ctx context.Context, channelId string, timeframe Timeframe, column string, req SeriesRequest) (items []HistogramItem, err error)
}

type IbcTransfer struct {
	bun.BaseModel `bun:"ibc_transfer" comment:"Table with IBC transfers."`

	Id              uint64          `bun:"id,pk,autoincrement"     comment:"Transfer internal identity"`
	Time            time.Time       `bun:"time,notnull,pk"         comment:"Message time"`
	Height          pkgTypes.Level  `bun:"height"                  comment:"Block number"`
	Amount          decimal.Decimal `bun:"amount,type:numeric"     comment:"Transferred amount"`
	Denom           string          `bun:"denom"                   comment:"Currency"`
	Memo            string          `bun:"memo"                    comment:"Memo"`
	ReceiverAddress *string         `bun:"receiver_address"        comment:"Receiver string. It's not null if it's not celestia address."`
	ReceiverId      *uint64         `bun:"receiver_id"             comment:"Receiver id. It's not null if it's celestia address."`
	SenderAddress   *string         `bun:"sender_address"          comment:"Sender string. It's not null if it's not celestia address."`
	SenderId        *uint64         `bun:"sender_id"               comment:"Sender id. It's not null if it's celestia address."`
	ConnectionId    string          `bun:"connection_id"           comment:"Connection identity"`
	ChannelId       string          `bun:"channel_id"              comment:"Channel identity"`
	Port            string          `bun:"port"                    comment:"Port"`
	Timeout         *time.Time      `bun:"timeout,nullzero"        comment:"Date-time timeout"`
	HeightTimeout   uint64          `bun:"height_timeout,nullzero" comment:"Height timeout"`
	Sequence        uint64          `bun:"sequence"                comment:"Sequence number of packet"`
	TxId            uint64          `bun:"tx_id"                   comment:"Transaction id where transfer occurred"`

	Tx         *Tx            `bun:"rel:belongs-to,join:tx_id=id"`
	Receiver   *Address       `bun:"rel:belongs-to,join:receiver_id=id"`
	Sender     *Address       `bun:"rel:belongs-to,join:sender_id=id"`
	Connection *IbcConnection `bun:"rel:belongs-to,join:connection_id=connection_id"`
}

func (IbcTransfer) TableName() string {
	return "ibc_transfer"
}
