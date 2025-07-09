// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

type ListConnectionFilters struct {
	Limit    int
	Offset   int
	Sort     sdk.SortOrder
	ClientId string
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IIbcConnection interface {
	ById(ctx context.Context, id string) (IbcConnection, error)
	List(ctx context.Context, fltrs ListConnectionFilters) ([]IbcConnection, error)
}

type IbcConnection struct {
	bun.BaseModel `bun:"ibc_connection" comment:"Table with IBC connections."`

	ConnectionId             string         `bun:"connection_id,pk"           comment:"Connection identity"`
	ClientId                 string         `bun:"client_id"                  comment:"Client identity"`
	CounterpartyConnectionId string         `bun:"counterparty_connection_id" comment:"Counterparty connection identity"`
	CounterpartyClientId     string         `bun:"counterparty_client_id"     comment:"Counterparty client identity"`
	CreatedAt                time.Time      `bun:"created_at"                 comment:"Time whe connection was created"`
	ConnectedAt              time.Time      `bun:"connected_at"               comment:"Time whe connection was established"`
	Height                   pkgTypes.Level `bun:"height"                     comment:"Block number when connection was created"`
	ConnectionHeight         pkgTypes.Level `bun:"connection_height"          comment:"Block number when connection was established"`
	CreateTxId               uint64         `bun:"create_tx_id"               comment:"Transaction id of creation connection"`
	ConnectionTxId           uint64         `bun:"connection_tx_id"           comment:"Transaction id of establishing connection"`
	ChannelsCount            int64          `bun:"channels_count"             comment:"Count of channels which was opened in the connection"`

	Client       *IbcClient `bun:"rel:belongs-to,join:client_id=id"`
	CreateTx     *Tx        `bun:"rel:belongs-to,join:create_tx_id=id"`
	ConnectionTx *Tx        `bun:"rel:belongs-to,join:connection_tx_id=id"`
}

func (IbcConnection) TableName() string {
	return "ibc_connection"
}
