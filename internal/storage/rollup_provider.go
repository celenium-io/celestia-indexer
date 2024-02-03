// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IRollupProvider interface {
	storage.Table[*RollupProvider]
}

// RollupProvider -
type RollupProvider struct {
	bun.BaseModel `bun:"rollup_provider" comment:"Table with data providers for rollups."`

	RollupId    uint64 `bun:"rollup_id,pk"    comment:"Unique internal rollup identity"`
	NamespaceId uint64 `bun:"namespace_id,pk" comment:"Namespace identity. May be NULL"`
	AddressId   uint64 `bun:"address_id,pk"   comment:"Celestia address of data provider"`
}

// TableName -
func (RollupProvider) TableName() string {
	return "rollup_provider"
}
