// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IVestingAccount interface {
	storage.Table[*VestingAccount]

	ByAddress(ctx context.Context, addressId uint64, limit, offset int, showEnded bool) ([]VestingAccount, error)
}

type VestingAccount struct {
	bun.BaseModel `bun:"vesting_account" comment:"Table with vesting accounts"`

	Id        uint64            `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height    pkgTypes.Level    `bun:"height,notnull"              comment:"The number (height) of this block"`
	Time      time.Time         `bun:"time,notnull"                comment:"The time of block"`
	TxId      *uint64           `bun:"tx_id"                       comment:"Transaction internal identity"`
	AddressId uint64            `bun:"address_id,notnull"          comment:"Address internal id"`
	Type      types.VestingType `bun:"type,type:vesting_type"      comment:"Type vesting account"`
	Amount    decimal.Decimal   `bun:"amount,type:numeric"         comment:"Vested amount"`
	StartTime *time.Time        `bun:"start_time"                  comment:"Start time of unlock value"`
	EndTime   *time.Time        `bun:"end_time"                    comment:"End time of unlock value"`

	Address *Address `bun:"rel:has-one"`
	Tx      *Tx      `bun:"rel:has-one"`

	VestingPeriods []VestingPeriod `bun:"rel:has-many"`
}

func (VestingAccount) TableName() string {
	return "vesting_account"
}
