// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IValidator interface {
	storage.Table[*Validator]

	ByAddress(ctx context.Context, address string) (Validator, error)
}

type Validator struct {
	bun.BaseModel `bun:"validator" comment:"Table with celestia validators."`

	Id        uint64 `bun:"id,pk,notnull,autoincrement"                comment:"Unique internal identity"`
	Delegator string `bun:"delegator,type:text"                        comment:"Delegator address"`
	Address   string `bun:"address,unique:address_validator,type:text" comment:"Validator address"`

	Moniker  string `bun:"moniker,type:text"  comment:"Human-readable name for the validator"`
	Website  string `bun:"website,type:text"  comment:"Website link"`
	Identity string `bun:"identity,type:text" comment:"Optional identity signature"`
	Contacts string `bun:"contacts,type:text" comment:"Contacts"`
	Details  string `bun:"details,type:text"  comment:"Detailed information about validator"`

	Rate              decimal.Decimal `bun:"rate,type:numeric"                comment:"Commission rate charged to delegators, as a fraction"`
	MaxRate           decimal.Decimal `bun:"max_rate,type:numeric"            comment:"Maximum commission rate which validator can ever charge, as a fraction"`
	MaxChangeRate     decimal.Decimal `bun:"max_change_rate,type:numeric"     comment:"Maximum daily increase of the validator commission, as a fraction"`
	MinSelfDelegation decimal.Decimal `bun:"min_self_delegation,type:numeric" comment:""`

	MsgId  uint64         `bun:"msg_id" comment:"Message id when validator was created"`
	Height pkgTypes.Level `bun:"height" comment:"Height when validator was created"`
}

func (Validator) TableName() string {
	return "validator"
}
