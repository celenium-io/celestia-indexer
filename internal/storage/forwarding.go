// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

type ForwardingFilter struct {
	Height    *uint64
	AddressId *uint64
	TxId      *uint64
	From      time.Time
	To        time.Time
	Sort      storage.SortOrder
	Limit     int
	Offset    int
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IForwarding interface {
	storage.Table[*Forwarding]

	Filter(ctx context.Context, filters ForwardingFilter) ([]Forwarding, error)
	ById(ctx context.Context, id uint64) (Forwarding, time.Time, error)
	Inputs(ctx context.Context, addressId uint64, from, to time.Time) ([]ForwardingInput, error)
}

type Forwarding struct {
	bun.BaseModel `bun:"forwarding" comment:"Table with forwarding events."`

	Id            uint64          `bun:"id,pk,notnull,autoincrement"       comment:"Unique internal id"`
	Height        pkgTypes.Level  `bun:"height,notnull"                    comment:"The number (height) of this block"`
	Time          time.Time       `bun:"time,pk,notnull"                   comment:"The time of block"`
	AddressId     uint64          `bun:"address_id,notnull"                comment:"Foreign key to addresses table"`
	DestDomain    uint64          `bun:"dest_domain,notnull"               comment:"The destination domain of the forwarding"`
	DestRecipient []byte          `bun:"dest_recipient,notnull,type:bytea" comment:"The destination recipient of the forwarding"`
	SuccessCount  uint64          `bun:"success_count"                     comment:"The number of successful forwarded tokens"`
	FailedCount   uint64          `bun:"failed_count"                      comment:"The number of failed forwarded tokens"`
	TxId          uint64          `bun:"tx_id,notnull"                     comment:"Foreign key to transactions table"`
	Transfers     json.RawMessage `bun:"transfers,type:jsonb"              comment:"The list of transfers included in the forwarding"`

	Address *Address `bun:"rel:has-one,join:address_id=id"`
	Tx      *Tx      `bun:"rel:has-one,join:tx_id=id"`
}

func (Forwarding) TableName() string {
	return "forwarding"
}

type ForwardingInput struct {
	Height       pkgTypes.Level    `bun:"height"`
	Time         time.Time         `bun:"time"`
	TxHash       []byte            `bun:"hash"`
	From         string            `bun:"src"`
	Amount       string            `bun:"amount"`
	Denom        string            `bun:"denom"`
	Counterparty uint64            `bun:"counterparty"`
	Data         types.PackedBytes `bun:"data"`
}
