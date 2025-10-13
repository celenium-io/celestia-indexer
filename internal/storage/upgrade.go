// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ListUpgradesFilter struct {
	Limit    int
	Offset   int
	Sort     sdk.SortOrder
	Height   uint64
	Signer   string
	TxHash   string
	SignerId *uint64
	TxId     *uint64
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IUpgrade interface {
	List(ctx context.Context, flts ListUpgradesFilter) ([]Upgrade, error)
}

type Upgrade struct {
	bun.BaseModel `bun:"upgrade" comment:"Table with upgrades"`

	Version      uint64          `bun:"version,pk"                comment:"Version"`
	Time         time.Time       `bun:"time"                      comment:"The time of first signal"`
	EndTime      time.Time       `bun:"end_time"                  comment:"The time of upgrade"`
	Height       pkgTypes.Level  `bun:"height"                    comment:"The number (height) of first signal block"`
	EndHeight    pkgTypes.Level  `bun:"end_height"                comment:"The number (height) of upgrade block"`
	SignerId     uint64          `bun:"signer_id"                 comment:"Signer internal identity"`
	MsgId        uint64          `bun:"msg_id,notnull"            comment:"Message internal identity"`
	TxId         uint64          `bun:"tx_id,notnull"             comment:"Transaction internal identity"`
	VotingPower  decimal.Decimal `bun:"voting_power,type:numeric" comment:"Total voting power on upgrade block"`
	VotedPower   decimal.Decimal `bun:"voted_power,type:numeric"  comment:"Total voting power of upgraded validators"`
	SignalsCount int             `bun:"signals_count"             comment:"Count of signals"`

	Signer *Address `bun:"rel:belongs-to,join:signer_id=id"`
	Tx     *Tx      `bun:"rel:belongs-to,join:tx_id=id"`
}

// TableName -
func (Upgrade) TableName() string {
	return "upgrade"
}
