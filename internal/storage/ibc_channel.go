// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ListChannelFilters struct {
	Limit        int
	Offset       int
	Sort         sdk.SortOrder
	ConnectionId string
	ClientId     string
	Status       types.IbcChannelStatus
}

type ChainStats struct {
	Chain    string          `bun:"chain_id"`
	Received decimal.Decimal `bun:"received"`
	Sent     decimal.Decimal `bun:"sent"`
	Flow     decimal.Decimal `bun:"flow"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IIbcChannel interface {
	ById(ctx context.Context, id string) (IbcChannel, error)
	List(ctx context.Context, fltrs ListChannelFilters) ([]IbcChannel, error)
	StatsByChain(ctx context.Context, limit, offset int) ([]ChainStats, error)
}

type IbcChannel struct {
	bun.BaseModel `bun:"ibc_channel" comment:"Table with IBC channels."`

	Id                    string                 `bun:"id,pk"                          comment:"Channel identity"`
	ConnectionId          string                 `bun:"connection_id"                  comment:"Connection identity"`
	PortId                string                 `bun:"port_id"                        comment:"Port id"`
	CounterpartyPortId    string                 `bun:"counterparty_port_id"           comment:"Counterparty port identity"`
	CounterpartyChannelId string                 `bun:"counterparty_channel_id"        comment:"Counterparty channel identity"`
	Version               string                 `bun:"version"                        comment:"Version"`
	ClientId              string                 `bun:"client_id"                      comment:"Client identity"`
	CreatedAt             time.Time              `bun:"created_at"                     comment:"Time whe channel was created"`
	ConfirmedAt           time.Time              `bun:"confirmed_at"                   comment:"Time whe channel was established"`
	Height                pkgTypes.Level         `bun:"height"                         comment:"Block number when channel was created"`
	ConfirmationHeight    pkgTypes.Level         `bun:"confirmation_height"            comment:"Block number when channel was confirmed"`
	CreateTxId            uint64                 `bun:"create_tx_id"                   comment:"Transaction id of creation channel"`
	ConfirmationTxId      uint64                 `bun:"confirmation_tx_id"             comment:"Transaction id of confirmation channel"`
	Ordering              bool                   `bun:"ordering"                       comment:"Ordered or unordered packets in the channel"`
	CreatorId             uint64                 `bun:"creator_id"                     comment:"Internal creator identity"`
	Status                types.IbcChannelStatus `bun:"status,type:ibc_channel_status" comment:"Channel status"`
	Received              decimal.Decimal        `bun:"received,type:numeric"          comment:"Received value"`
	Sent                  decimal.Decimal        `bun:"sent,type:numeric"              comment:"Sent value"`
	TransfersCount        int64                  `bun:"transfers_count"                comment:"Count transfers"`

	Connection     *IbcConnection `bun:"rel:belongs-to,join:connection_id=connection_id"`
	Client         *IbcClient     `bun:"rel:belongs-to,join:client_id=id"`
	CreateTx       *Tx            `bun:"rel:belongs-to,join:create_tx_id=id"`
	ConfirmationTx *Tx            `bun:"rel:belongs-to,join:confirmation_tx_id=id"`
	Creator        *Address       `bun:"rel:belongs-to,join:creator_id=id"`
}

func (IbcChannel) TableName() string {
	return "ibc_channel"
}
