package node

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type API interface {
	Head(ctx context.Context) (types.ResultBlock, error)
	Block(ctx context.Context, level storage.Level) (types.ResultBlock, error)
	BlockResults(ctx context.Context, level storage.Level) (types.ResultBlockResults, error)
	Genesis(ctx context.Context) (types.Genesis, error)
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type CelestiaNodeApi interface {
	Blobs(ctx context.Context, height uint64, hash ...string) ([]types.Blob, error)
}
