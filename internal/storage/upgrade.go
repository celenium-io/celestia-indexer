// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
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

	Id       uint64         `bun:"id,pk"           comment:"Unique identity"`
	Height   pkgTypes.Level `bun:"height"          comment:"The number (height) of this block"`
	SignerId uint64         `bun:"signer_id"       comment:"Signer internal identity"`
	Time     time.Time      `bun:"time,pk,notnull" comment:"The time of upgrade"`
	Version  uint64         `bun:"version"         comment:"Version"`
	MsgId    uint64         `bun:"msg_id,notnull"  comment:"Message internal identity"`
	TxId     uint64         `bun:"tx_id,notnull"   comment:"Transaction internal identity"`

	Signer *Address `bun:"rel:belongs-to,join:signer_id=id"`
}

// TableName -
func (Upgrade) TableName() string {
	return "upgrade"
}
