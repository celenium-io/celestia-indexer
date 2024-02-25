// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
type IStakingLog interface {
	storage.Table[*StakingLog]

	ByValidator(ctx context.Context, validatorId uint64, limit, offset int) ([]StakingLog, error)
}

// Delegation -
type StakingLog struct {
	bun.BaseModel `bun:"staking_log" comment:"Table with staking events log"`

	Id          uint64               `bun:"id,pk,notnull,autoincrement" comment:"Unique internal id"`
	Time        time.Time            `bun:"time,pk,notnull"             comment:"The time of block"                 stats:"func:min max,filterable"`
	Height      pkgTypes.Level       `bun:"height,notnull"              comment:"The number (height) of this block" stats:"func:min max,filterable"`
	AddressId   *uint64              `bun:"address_id"                  comment:"Internal address id"`
	ValidatorId uint64               `bun:"validator_id"                comment:"Internal validator id"`
	Change      decimal.Decimal      `bun:"change,type:numeric"         comment:"Change amount"`
	Type        types.StakingLogType `bun:"type,type:staking_log_type"  comment:"Staking log type"`

	Address   *Address   `bun:"rel:belongs-to,join:address_id=id"`
	Validator *Validator `bun:"rel:belongs-to,join:validator_id=id"`
}

// TableName -
func (StakingLog) TableName() string {
	return "staking_log"
}
