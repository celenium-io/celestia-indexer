// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"strings"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IDelegation interface {
	storage.Table[*Delegation]

	ByAddress(ctx context.Context, addressId uint64, limit, offset int, showZero bool) ([]Delegation, error)
	ByValidator(ctx context.Context, validatorId uint64, limit, offset int, showZero bool) ([]Delegation, error)
}

// Delegation -
type Delegation struct {
	bun.BaseModel `bun:"delegation" comment:"Table with delegations"`

	Id          uint64          `bun:"id,pk,notnull,autoincrement"         comment:"Unique internal id"`
	AddressId   uint64          `bun:"address_id,unique:delegation_pair"   comment:"Internal address id"`
	ValidatorId uint64          `bun:"validator_id,unique:delegation_pair" comment:"Internal validator id"`
	Amount      decimal.Decimal `bun:"amount,type:numeric"                 comment:"Delegated amount"`

	Address   *Address   `bun:"rel:belongs-to,join:address_id=id"`
	Validator *Validator `bun:"rel:belongs-to,join:validator_id=id"`
}

// TableName -
func (Delegation) TableName() string {
	return "delegation"
}

func (d Delegation) String() string {
	sb := new(strings.Builder)
	if d.Address != nil {
		sb.WriteString(d.Address.Address)
	}
	if d.Validator != nil {
		sb.WriteString(d.Validator.Address)
	}
	return sb.String()
}
