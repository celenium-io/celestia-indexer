// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package node

import (
	"context"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Api interface {
	Status(ctx context.Context) (types.Status, error)
	Head(ctx context.Context) (pkgTypes.ResultBlock, error)
	Block(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlock, error)
	BlockResults(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlockResults, error)
	Genesis(ctx context.Context) (types.Genesis, error)
	BlockData(ctx context.Context, level pkgTypes.Level) (pkgTypes.BlockData, error)
	BlockDataGet(ctx context.Context, level pkgTypes.Level) (pkgTypes.BlockData, error)
	BlockBulkData(ctx context.Context, levels ...pkgTypes.Level) ([]pkgTypes.BlockData, error)
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type DalApi interface {
	Blobs(ctx context.Context, height pkgTypes.Level, hash ...string) ([]types.Blob, error)
	Blob(ctx context.Context, height pkgTypes.Level, namespace, commitment string) (types.Blob, error)
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type CosmosApi interface {
	ModuleAccounts(ctx context.Context) ([]types.Account, error)
}
