// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"strconv"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IConstant interface {
	Get(ctx context.Context, module types.ModuleName, name string) (Constant, error)
	ByModule(ctx context.Context, module types.ModuleName) ([]Constant, error)
	All(ctx context.Context) ([]Constant, error)
}

type Constant struct {
	bun.BaseModel `bun:"table:constant" comment:"Table with celestia constants."`

	Module types.ModuleName `bun:"module,pk,type:module_name" comment:"Module name which declares constant"`
	Name   string           `bun:"name,pk,type:text"          comment:"Constant name"`
	Value  string           `bun:"value,type:text"            comment:"Constant value"`
}

func (Constant) TableName() string {
	return "constant"
}

func (c Constant) MustUint64() uint64 {
	i, err := strconv.ParseUint(c.Value, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (c Constant) MustUint32() uint32 {
	i, err := strconv.ParseUint(c.Value, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(i)
}
