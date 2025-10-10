// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IRollupProvider interface {
	ByRollupId(ctx context.Context, rollupId uint64) ([]RollupProvider, error)
}

// RollupProvider -
type RollupProvider struct {
	bun.BaseModel `bun:"rollup_provider" comment:"Table with data providers for rollups."`

	RollupId    uint64 `bun:"rollup_id,pk"    comment:"Unique internal rollup identity"`
	NamespaceId uint64 `bun:"namespace_id,pk" comment:"Namespace identity. May be NULL"`
	AddressId   uint64 `bun:"address_id,pk"   comment:"Celestia address of data provider"`

	Namespace *Namespace `bun:"rel:belongs-to,join:namespace_id=id"`
	Address   *Address   `bun:"rel:belongs-to,join:address_id=id"`
}

// TableName -
func (RollupProvider) TableName() string {
	return "rollup_provider"
}
