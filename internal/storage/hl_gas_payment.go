// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IHLGasPayment interface {
	List(ctx context.Context, limit, offset int) ([]HLGasPayment, error)
}

type HLGasPayment struct {
	bun.BaseModel `bun:"hl_gas_payment" comment:"Table with hyperlane gas payment"`

	Id         uint64          `bun:"id,pk,autoincrement"     comment:"Internal identity"`
	Height     pkgTypes.Level  `bun:"height,notnull"          comment:"The number (height) of this block"`
	Time       time.Time       `bun:"time,pk,notnull"         comment:"The time of block"`
	TransferId uint64          `bun:"transfer_id"             comment:"Transfer identity"`
	GasAmount  decimal.Decimal `bun:"gas_amount,type:numeric" comment:"Gas amount"`
	IgpId      uint64          `bun:"igp_id"                  comment:"IGP identity"`
	Amount     decimal.Decimal `bun:"amount,type:numeric"     comment:"Amount"`

	Igp *HLIGP `bun:"rel:belongs-to,join:igp_id=id"`
}

func (m *HLGasPayment) TableName() string {
	return "hl_gas_payment"
}
