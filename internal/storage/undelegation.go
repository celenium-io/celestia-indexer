// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
type IUndelegation interface {
	storage.Table[*Undelegation]

	ByAddress(ctx context.Context, addressId uint64, limit, offset int) ([]Undelegation, error)
}

// Undelegation -
type Undelegation struct {
	bun.BaseModel `bun:"undelegation" comment:"Table with undelegations"`

	Id             uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal id"`
	Time           time.Time       `bun:"time,notnull"                comment:"The time of block"`
	Height         pkgTypes.Level  `bun:",notnull"                    comment:"The number (height) of this block"`
	AddressId      uint64          `bun:"address_id"                  comment:"Internal address id"`
	ValidatorId    uint64          `bun:"validator_id"                comment:"Internal validator id"`
	Amount         decimal.Decimal `bun:"amount,type:numeric"         comment:"Delegated amount"`
	CompletionTime time.Time       `bun:"completion_time"             comment:"Time when undelegation will be completed"`

	Address   *Address   `bun:"rel:belongs-to,join:address_id=id"`
	Validator *Validator `bun:"rel:belongs-to,join:validator_id=id"`
}

// TableName -
func (Undelegation) TableName() string {
	return "undelegation"
}
