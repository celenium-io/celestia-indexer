// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/hex"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IHLIGP interface {
	List(ctx context.Context, limit, offset int) ([]HLIGP, error)
	ByHash(ctx context.Context, hash []byte) (HLIGP, error)
}

type HLIGP struct {
	bun.BaseModel `bun:"hl_igp" comment:"Table with hyperlane interchain gas paymaster (IGP)"`

	Id      uint64         `bun:"id,pk,autoincrement"      comment:"Internal identity"`
	Height  pkgTypes.Level `bun:"height,notnull"           comment:"The number (height) of this block"`
	Time    time.Time      `bun:"time,pk,notnull"          comment:"The time of block"`
	IgpId   []byte         `bun:"igp_id,type:bytea,unique" comment:"IGP id"`
	OwnerId uint64         `bun:"owner_id"                 comment:"Owner identity"`
	Denom   string         `bun:"denom"                    comment:"Denom"`

	Owner  *Address     `bun:"rel:belongs-to,join:owner_id=id"`
	Config *HLIGPConfig `bun:"rel:belongs-to"`
}

func (t *HLIGP) TableName() string {
	return "hl_igp"
}

func (t *HLIGP) String() string {
	return hex.EncodeToString(t.IgpId)
}
