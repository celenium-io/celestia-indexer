// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IState interface {
	storage.Table[*State]

	ByName(ctx context.Context, name string) (State, error)
}

// State -
type State struct {
	bun.BaseModel `bun:"state" comment:"Current indexer state"`

	Id              uint64          `bun:",pk,autoincrement"         comment:"Unique internal identity"`
	Name            string          `bun:",unique:state_name"        comment:"Indexer name"`
	LastHeight      types.Level     `bun:"last_height"               comment:"Last block height"`
	LastHash        []byte          `bun:"last_hash"                 comment:"Last block hash"`
	LastTime        time.Time       `bun:"last_time"                 comment:"Time of last block"`
	ChainId         string          `bun:"chain_id"                  comment:"Celestia chain id"`
	TotalTx         int64           `bun:"total_tx"                  comment:"Transactions count in celestia"`
	TotalAccounts   int64           `bun:"total_accounts"            comment:"Accounts count in celestia"`
	TotalNamespaces int64           `bun:"total_namespaces"          comment:"Namespaces count in celestia"`
	TotalBlobsSize  int64           `bun:"total_blobs_size"          comment:"Total blobs size"`
	TotalSupply     decimal.Decimal `bun:"total_supply,type:numeric" comment:"Total supply in celestia"`
	TotalFee        decimal.Decimal `bun:"total_fee,type:numeric"    comment:"Total paid fee"`
}

// TableName -
func (State) TableName() string {
	return "state"
}
