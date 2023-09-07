package rollback

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

// TODO: rollback account balances by events
func (module *Module) balances(ctx context.Context, deletedEvents []storage.Event, deletedAddresses []storage.Address) error {
	return nil
}
