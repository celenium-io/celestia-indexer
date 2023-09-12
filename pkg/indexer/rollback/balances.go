package rollback

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

// TODO: rollback account rollbackBalances by events
func (module *Module) rollbackBalances(ctx context.Context, deletedEvents []storage.Event, deletedAddresses []storage.Address) error {
	return nil
}
