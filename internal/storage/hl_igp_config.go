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
type IHLIGPConfig interface {
	List(ctx context.Context, limit, offset int) ([]HLIGPConfig, error)
}

type HLIGPConfig struct {
	bun.BaseModel `bun:"hl_igp_config" comment:"Table with hyperlane interchain gas paymaster (IGP) config"`

	Id                uint64          `bun:"id,pk"               comment:"Internal identity"`
	Height            pkgTypes.Level  `bun:"height,notnull"      comment:"The number (height) of this block"`
	Time              time.Time       `bun:"time,pk,notnull"     comment:"The time of block"`
	GasOverhead       decimal.Decimal `bun:"gas_overhead"        comment:"Gas overhead"`
	GasPrice          decimal.Decimal `bun:"gas_price"           comment:"Gas price"`
	RemoteDomain      uint64          `bun:"remote_domain"       comment:"Remote domain"`
	TokenExchangeRate string          `bun:"token_exchange_rate" comment:"Token exchange rate"`
}

func (t *HLIGPConfig) TableName() string {
	return "hl_igp_config"
}
