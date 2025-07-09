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
type IVestingPeriod interface {
	storage.Table[*VestingPeriod]

	ByVesting(ctx context.Context, id uint64, limit, offset int) ([]VestingPeriod, error)
}

type VestingPeriod struct {
	bun.BaseModel `bun:"vesting_period" comment:"Table with vesting periods"`

	Id               uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height           pkgTypes.Level  `bun:"height,notnull"              comment:"The number (height) of this block"`
	VestingAccountId uint64          `bun:"vesting_account_id"          comment:"Vesting account internal identity"`
	Time             time.Time       `bun:"time,notnull"                comment:"The time of periodic vesting"`
	Amount           decimal.Decimal `bun:"amount,type:numeric"         comment:"Vested amount"`
}

func (VestingPeriod) TableName() string {
	return "vesting_period"
}
