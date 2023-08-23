package postgres

import (
	"context"

	models "github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
)

// Storage -
type Storage struct {
	*postgres.Storage

	cfg config.Database

	Blocks      models.IBlock
	Tx          models.ITx
	Message     models.IMessage
	Event       models.IEvent
	Address     models.IAddress
	Namespace   models.INamespace
	State       models.IState
	Notificator *Notificator

	PartitionManager database.RangePartitionManager
}

// Create -
func Create(ctx context.Context, cfg config.Database) (Storage, error) {
	strg, err := postgres.Create(ctx, cfg, initDatabase)
	if err != nil {
		return Storage{}, err
	}

	s := Storage{
		cfg:         cfg,
		Storage:     strg,
		Blocks:      NewBlocks(strg.Connection()),
		Message:     NewMessage(strg.Connection()),
		Event:       NewEvent(strg.Connection()),
		Address:     NewAddress(strg.Connection()),
		Tx:          NewTx(strg.Connection()),
		State:       NewState(strg.Connection()),
		Namespace:   NewNamespace(strg.Connection()),
		Notificator: NewNotificator(cfg, strg.Connection().DB()),

		PartitionManager: database.NewPartitionManager(strg.Connection(), database.PartitionByMonth),
	}

	return s, nil
}

func initDatabase(ctx context.Context, conn *database.Bun) error {
	if err := createTypes(ctx, conn); err != nil {
		return errors.Wrap(err, "creating custom types")
	}

	if err := database.CreateTables(ctx, conn, models.Models...); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return err
	}

	if err := database.MakeComments(ctx, conn, models.Models...); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return errors.Wrap(err, "make comments")
	}

	return createIndices(ctx, conn)
}

func (s Storage) CreateListener() models.Listener {
	return NewNotificator(s.cfg, s.Notificator.db)
}
