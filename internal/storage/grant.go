// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IGrant interface {
	storage.Table[*Grant]

	ByGranter(ctx context.Context, id uint64, limit, offset int) ([]Grant, error)
	ByGrantee(ctx context.Context, id uint64, limit, offset int) ([]Grant, error)
}

type Grant struct {
	bun.BaseModel `bun:"grant" comment:"Table with grants"`

	Id            uint64         `bun:"id,pk,notnull,autoincrement"    comment:"Unique internal identity"`
	Height        types.Level    `bun:"height"                         comment:"Block height"`
	RevokeHeight  *types.Level   `bun:"revoke_height"                  comment:"Block height when grant was revoked"`
	Time          time.Time      `bun:"time"                           comment:"The time of block"`
	GranterId     uint64         `bun:"granter_id,unique:grant_key"    comment:"Granter internal identity"`
	GranteeId     uint64         `bun:"grantee_id,unique:grant_key"    comment:"Grantee internal identity"`
	Authorization string         `bun:"authorization,unique:grant_key" comment:"Authorization type"`
	Expiration    *time.Time     `bun:"expiration"                     comment:"Expiration time"`
	Revoked       bool           `bun:"revoked"                        comment:"Is grant revoked"`
	Params        map[string]any `bun:"params,type:jsonb,nullzero"     comment:"Authorization parameters"`

	Granter *Address `bun:"rel:has-one"`
	Grantee *Address `bun:"rel:has-one"`
}

func (Grant) TableName() string {
	return "grant"
}

func (g Grant) String() string {
	return fmt.Sprintf("%s_%s_%s", g.Authorization, g.Granter.Address, g.Grantee.Address)
}
