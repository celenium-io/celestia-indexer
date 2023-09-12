package indexer

import (
	"context"
	"sync"

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
	cfg          config.Config
	api          node.API
	receiver     *receiver.Module
	parser       *parser.Module
	storage      *storage.Module
	rollback     *rollback.Module
	genesis      *genesis.Module
	stopperInput *modules.Input
	wg           *sync.WaitGroup
	log          zerolog.Logger
}

func New(ctx context.Context, cfg config.Config, stopperInput *modules.Input) (Indexer, error) {
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

	err = attachStopper(stopperInput, r, p, s, rb, genesisModule)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating stopper module")
	}

	return Indexer{
		cfg:          cfg,
		api:          &api,
		receiver:     &r,
		parser:       &p,
		storage:      &s,
		rollback:     &rb,
		genesis:      &genesisModule,
		stopperInput: stopperInput,
		wg:           new(sync.WaitGroup),
		log:          log.With().Str("module", "indexer").Logger(),
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

func createReceiver(ctx context.Context, cfg config.Config, pg postgres.Storage) (rpc.API, receiver.Module, error) {
	state, err := loadState(pg, ctx, cfg.Indexer.Name)
	if err != nil {
		return rpc.API{}, receiver.Module{}, errors.Wrap(err, "while loading state")
	}

	api := rpc.NewAPI(cfg.DataSources["node_rpc"])
	receiverModule := receiver.NewModule(cfg.Indexer, &api, state)
	return api, receiverModule, nil
}

func createRollback(r receiver.Module, pg postgres.Storage, api node.API, cfg config.Indexer) (rollback.Module, error) {
	rollbackModule := rollback.NewModule(pg.Transactable, pg.State, pg.Blocks, api, cfg)

	// rollback <- listen signal -- receiver
	rbInput, err := rollbackModule.Input(rollback.InputName)
	if err != nil {
		return rollback.Module{}, err
	}

	if err = r.AttachTo(receiver.RollbackOutput, rbInput); err != nil {
		return rollback.Module{}, errors.Wrap(err, "while attaching rollback to receiver")
	}

	// receiver <- listen state -- rollback
	rInput, err := r.Input(receiver.RollbackInput)
	if err != nil {
		return rollback.Module{}, err
	}

	if err = rollbackModule.AttachTo(rollback.OutputName, rInput); err != nil {
		return rollback.Module{}, errors.Wrap(err, "while attaching receiver to rollback")
	}

	return rollbackModule, nil
}

func createParser(r receiver.Module) (parser.Module, error) {
	parserModule := parser.NewModule()
	pInput, err := parserModule.Input(parser.InputName)
	if err != nil {
		return parser.Module{}, err
	}

	if err = r.AttachTo(receiver.BlocksOutput, pInput); err != nil {
		return parser.Module{}, errors.Wrap(err, "while attaching parser to receiver")
	}

	return parserModule, nil
}

func createStorage(pg postgres.Storage, cfg config.Config, p parser.Module) (storage.Module, error) {
	s := storage.NewModule(pg, cfg.Indexer)
	sInput, err := s.Input(storage.InputName)
	if err != nil {
		return storage.Module{}, err
	}

	if err = p.AttachTo(parser.OutputName, sInput); err != nil {
		return storage.Module{}, err
	}

	return s, nil
}

func createGenesis(pg postgres.Storage, cfg config.Config, r receiver.Module) (genesis.Module, error) {
	genesisModule := genesis.NewModule(pg, cfg.Indexer)
	gInput, err := genesisModule.Input(genesis.InputName)
	if err != nil {
		return genesis.Module{}, errors.Wrap(err, "cannot find input in genesis")
	}
	if err = r.AttachTo(receiver.GenesisOutput, gInput); err != nil {
		return genesis.Module{}, err
	}
	receiverGenesisDone, err := r.Input(receiver.GenesisDoneInput)
	if err != nil {
		return genesis.Module{}, errors.Wrap(err, "cannot find input in receiver")
	}
	if err = genesisModule.AttachTo(genesis.OutputName, receiverGenesisDone); err != nil {
		return genesis.Module{}, err
	}
	return genesisModule, nil
}

func attachStopper(stopperInput *modules.Input, r receiver.Module, p parser.Module, s storage.Module, rb rollback.Module, g genesis.Module) error {
	if err := r.AttachTo(receiver.StopOutput, stopperInput); err != nil {
		return errors.Wrap(err, "while attaching stopper to receiver")
	}

	if err := p.AttachTo(parser.StopOutput, stopperInput); err != nil {
		return errors.Wrap(err, "while attaching stopper to parser")
	}

	if err := s.AttachTo(storage.StopOutput, stopperInput); err != nil {
		return errors.Wrap(err, "while attaching stopper to storage")
	}

	if err := rb.AttachTo(rollback.StopOutput, stopperInput); err != nil {
		return errors.Wrap(err, "while attaching stopper to rollback")
	}

	if err := g.AttachTo(genesis.StopOutput, stopperInput); err != nil {
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
