package node

import (
	"context"

	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"

	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Api interface {
	Status(ctx context.Context) (types.Status, error)
	Head(ctx context.Context) (pkgTypes.ResultBlock, error)
	Block(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlock, error)
	BlockResults(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlockResults, error)
	Genesis(ctx context.Context) (types.Genesis, error)
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type DalApi interface {
	Blobs(ctx context.Context, height uint64, hash ...string) ([]types.Blob, error)
	Blob(ctx context.Context, height uint64, namespace, commitment string) (types.Blob, error)
}
