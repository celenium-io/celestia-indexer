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
type IJail interface {
	storage.Table[*Jail]

	ByValidator(ctx context.Context, id uint64, limit, offset int) ([]Jail, error)
}

// Jail -
type Jail struct {
	bun.BaseModel `bun:"jail" comment:"Table with all jailed events."`

	Id          uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal id"`
	Time        time.Time       `bun:"time,pk,notnull"             comment:"The time of block"`
	Height      pkgTypes.Level  `bun:"height,notnull"              comment:"The number (height) of this block"`
	ValidatorId uint64          `bun:"validator_id,notnull"        comment:"Internal validator id"`
	Reason      string          `bun:"reason"                      comment:"Reason"`
	Burned      decimal.Decimal `bun:"burned,type:numeric"         comment:"Burned coins"`

	Validator *Validator `bun:"rel:belongs-to,join:validator_id=id"`
}

// TableName -
func (Jail) TableName() string {
	return "jail"
}
