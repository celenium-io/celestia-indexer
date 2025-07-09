// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlockSignature interface {
	storage.Table[*BlockSignature]

	LevelsByValidator(ctx context.Context, validatorId uint64, startHeight types.Level) ([]types.Level, error)
}

type BlockSignature struct {
	bun.BaseModel `bun:"block_signature" comment:"Table with block signatures"`

	Id          uint64      `bun:"id,pk,notnull,autoincrement" comment:"Unique internal id"`
	Height      types.Level `bun:",notnull"                    comment:"The number (height) of this block"`
	Time        time.Time   `bun:"time,pk,notnull"             comment:"The time of block"`
	ValidatorId uint64      `bun:"validator_id"                comment:"Validator's internal identity"`

	Validator *Validator `bun:"rel:belongs-to"`
}

func (BlockSignature) TableName() string {
	return "block_signature"
}
