// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IRedelegation interface {
	storage.Table[*Redelegation]

	ByAddress(ctx context.Context, addressId uint64, limit, offset int) ([]Redelegation, error)
}

// Redelegation -
type Redelegation struct {
	bun.BaseModel `bun:"redelegation" comment:"Table with redelegations"`

	Id             uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal id"`
	Time           time.Time       `bun:"time,notnull"                comment:"The time of block"`
	Height         pkgTypes.Level  `bun:",notnull"                    comment:"The number (height) of this block"`
	AddressId      uint64          `bun:"address_id"                  comment:"Internal address id"`
	SrcId          uint64          `bun:"src_id"                      comment:"Internal source validator id"`
	DestId         uint64          `bun:"dest_id"                     comment:"Internal destination validator id"`
	Amount         decimal.Decimal `bun:"amount,type:numeric"         comment:"Delegated amount"`
	CompletionTime time.Time       `bun:"completion_time"             comment:"Time when redelegation will be completed"`

	Address     *Address   `bun:"rel:belongs-to,join:address_id=id"`
	Source      *Validator `bun:"rel:belongs-to,join:src_id=id"`
	Destination *Validator `bun:"rel:belongs-to,join:dest_id=id"`
}

// TableName -
func (Redelegation) TableName() string {
	return "redelegation"
}
