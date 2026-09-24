// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

type FilterBondUpdatesListByValidator struct {
	Offset int
	Limit  int
	Sort   storage.SortOrder
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IValidatorBondUpdate interface {
	ListByValidator(ctx context.Context, validatorId uint64, filters FilterBondUpdatesListByValidator) ([]ValidatorBondUpdate, error)
}

type ValidatorBondUpdate struct {
	bun.BaseModel `bun:"validator_bond_update" comment:"Table with validator power updates from validator_updates of block results."`

	Id          uint64         `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Time        time.Time      `bun:"time,pk,notnull"             comment:"The time of block"`
	Height      pkgTypes.Level `bun:"height,notnull"              comment:"The number (height) of this block"`
	ValidatorId uint64         `bun:"validator_id"                comment:"Internal validator id"`
	Power       *types.Numeric `bun:"power,type:numeric"          comment:"Validator's power"`

	Validator *Validator `bun:"-"`
}

func (ValidatorBondUpdate) TableName() string {
	return "validator_bond_update"
}
