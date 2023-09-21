package indexer

import (
	"context"
	"sync"

	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"

	internalStorage "github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/genesis"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/parser"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/rollback"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/node"
	"github.com/dipdup-io/celestia-indexer/pkg/node/rpc"
	"github.com/pkg/errors"

	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/receiver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Indexer struct {
	cfg      config.Config
	api      node.Api
	receiver *receiver.Module
	parser   *parser.Module
	storage  *storage.Module
	rollback *rollback.Module
	genesis  *genesis.Module
	stopper  modules.Module
	wg       *sync.WaitGroup
	log      zerolog.Logger
}

func New(ctx context.Context, cfg config.Config, stopperModule modules.Module) (Indexer, error) {
	pg, err := postgres.Create(ctx, cfg.Database)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating pg context")
	}

	api, r, err := createReceiver(ctx, cfg, pg)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating receiver module")
	}

	rb, err := createRollback(r, pg, &api, cfg.Indexer)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating rollback module")
	}

	p, err := createParser(r)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating parser module")
	}

	s, err := createStorage(pg, cfg, p)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating storage module")
	}

	genesisModule, err := createGenesis(pg, cfg, r)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating genesis module")
	}

	err = attachStopper(stopperModule, r, p, s, rb, genesisModule)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating stopper module")
	}

	return Indexer{
		cfg:      cfg,
		api:      &api,
		receiver: r,
		parser:   p,
		storage:  s,
		rollback: rb,
		genesis:  genesisModule,
		stopper:  stopperModule,
		wg:       new(sync.WaitGroup),
		log:      log.With().Str("module", "indexer").Logger(),
	}, nil
}

func (i *Indexer) Start(ctx context.Context) {
	i.log.Info().Msg("starting...")

	i.genesis.Start(ctx)
	i.storage.Start(ctx)
	i.parser.Start(ctx)
	i.receiver.Start(ctx)
}

func (i *Indexer) Close() error {
	i.log.Info().Msg("closing...")
	i.wg.Wait()

	if err := i.receiver.Close(); err != nil {
		log.Err(err).Msg("closing receiver")
	}
	if err := i.genesis.Close(); err != nil {
		log.Err(err).Msg("closing genesis")
	}
	if err := i.parser.Close(); err != nil {
		log.Err(err).Msg("closing parser")
	}
	if err := i.storage.Close(); err != nil {
		log.Err(err).Msg("closing storage")
	}
	if err := i.rollback.Close(); err != nil {
		log.Err(err).Msg("closing rollback")
	}

	return nil
}

func createReceiver(ctx context.Context, cfg config.Config, pg postgres.Storage) (rpc.API, *receiver.Module, error) {
	state, err := loadState(pg, ctx, cfg.Indexer.Name)
	if err != nil {
		return rpc.API{}, nil, errors.Wrap(err, "while loading state")
	}

	api := rpc.NewAPI(cfg.DataSources["node_rpc"])
	receiverModule := receiver.NewModule(cfg.Indexer, &api, state)
	return api, &receiverModule, nil
}

func createRollback(receiverModule modules.Module, pg postgres.Storage, api node.Api, cfg config.Indexer) (*rollback.Module, error) {
	rollbackModule := rollback.NewModule(pg.Transactable, pg.State, pg.Blocks, api, cfg)

	// rollback <- listen signal -- receiver
	if err := rollbackModule.AttachTo(receiverModule, receiver.RollbackOutput, rollback.InputName); err != nil {
		return nil, errors.Wrap(err, "while attaching rollback to receiver")
	}

	// receiver <- listen state -- rollback
	if err := receiverModule.AttachTo(&rollbackModule, rollback.OutputName, receiver.RollbackInput); err != nil {
		return nil, errors.Wrap(err, "while attaching receiver to rollback")
	}

	return &rollbackModule, nil
}

func createParser(receiverModule modules.Module) (*parser.Module, error) {
	parserModule := parser.NewModule()

	if err := parserModule.AttachTo(receiverModule, receiver.BlocksOutput, parser.InputName); err != nil {
		return nil, errors.Wrap(err, "while attaching parser to receiver")
	}

	return &parserModule, nil
}

func createStorage(pg postgres.Storage, cfg config.Config, parserModule modules.Module) (*storage.Module, error) {
	storageModule := storage.NewModule(pg.Transactable, pg.Notificator, cfg.Indexer)

	if err := storageModule.AttachTo(parserModule, parser.OutputName, storage.InputName); err != nil {
		return nil, errors.Wrap(err, "while attaching storage to parser")
	}

	return &storageModule, nil
}

func createGenesis(pg postgres.Storage, cfg config.Config, receiverModule modules.Module) (*genesis.Module, error) {
	genesisModule := genesis.NewModule(pg, cfg.Indexer)

	if err := genesisModule.AttachTo(receiverModule, receiver.GenesisOutput, genesis.InputName); err != nil {
		return nil, errors.Wrap(err, "while attaching genesis to receiver")
	}

	genesisModulePtr := &genesisModule
	if err := receiverModule.AttachTo(genesisModulePtr, genesis.OutputName, receiver.GenesisDoneInput); err != nil {
		return nil, errors.Wrap(err, "while attaching receiver to genesis")
	}

	return genesisModulePtr, nil
}

func attachStopper(stopperModule modules.Module, receiverModule modules.Module, parserModule modules.Module, storageModule modules.Module, rollbackModule modules.Module, genesisModule modules.Module) error {
	if err := stopperModule.AttachTo(receiverModule, receiver.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to receiver")
	}

	if err := stopperModule.AttachTo(parserModule, parser.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to parser")
	}

	if err := stopperModule.AttachTo(storageModule, storage.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to storage")
	}

	if err := stopperModule.AttachTo(rollbackModule, rollback.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to rollback")
	}

	if err := stopperModule.AttachTo(genesisModule, genesis.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to genesis")
	}

	return nil
}

func loadState(pg postgres.Storage, ctx context.Context, indexerName string) (*internalStorage.State, error) {
	state, err := pg.State.ByName(ctx, indexerName)
	if err != nil {
		if pg.State.IsNoRows(err) {
			return nil, nil
		}

		return nil, err
	}

	return &state, nil
}
