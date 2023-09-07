package rollback

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/shopspring/decimal"
)

// TODO: compute total supply diff by deleted events
func (module *Module) totalSupplyDiff(ctx context.Context, deletedEvents []storage.Event) (decimal.Decimal, error) {
	return decimal.Zero, nil
}
