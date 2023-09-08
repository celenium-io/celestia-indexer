package indexer

import (
	"context"
	"sync"

	internalStorage "github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/genesis"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/parser"
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
	api      node.API
	receiver *receiver.Module
	parser   *parser.Module
	storage  *storage.Module
	genesis  *genesis.Module
	wg       *sync.WaitGroup
	log      zerolog.Logger
}

func New(ctx context.Context, cfg config.Config) (Indexer, error) {
	pg, err := postgres.Create(ctx, cfg.Database)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating pg context")
	}

	state, err := LoadState(pg, ctx, cfg.Indexer.Name)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while loading state")
	}

	api := rpc.NewAPI(cfg.DataSources["node_rpc"])
	r := receiver.NewModule(cfg.Indexer, &api, state)

	p := parser.NewModule()
	pInput, err := p.Input(parser.BlocksInput)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "cannot find input in parser")
	}
	if err = r.AttachTo(receiver.BlocksOutput, pInput); err != nil {
		return Indexer{}, err
	}

	s := storage.NewModule(pg, cfg.Indexer)
	sInput, err := s.Input(storage.InputName)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "cannot find input in storage")
	}
	if err = p.AttachTo(parser.DataOutput, sInput); err != nil {
		return Indexer{}, err
	}

	genesisModule := genesis.NewModule(pg, cfg.Indexer)
	gInput, err := genesisModule.Input(genesis.InputName)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "cannot find input in genesis")
	}
	if err = r.AttachTo(receiver.GenesisOutput, gInput); err != nil {
		return Indexer{}, err
	}
	receiverGenesisDone, err := r.Input(receiver.GenesisDoneInput)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "cannot find input in receiver")
	}
	if err = genesisModule.AttachTo(genesis.OutputName, receiverGenesisDone); err != nil {
		return Indexer{}, err
	}

	return Indexer{
		cfg:      cfg,
		api:      &api,
		receiver: &r,
		parser:   &p,
		storage:  &s,
		genesis:  &genesisModule,
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
		log.Err(err).Msg("closing receiver")
	}
	if err := i.parser.Close(); err != nil {
		log.Err(err).Msg("closing receiver")
	}
	if err := i.storage.Close(); err != nil {
		log.Err(err).Msg("closing receiver")
	}

	return nil
}

func LoadState(pg postgres.Storage, ctx context.Context, indexerName string) (*internalStorage.State, error) {
	state, err := pg.State.ByName(ctx, indexerName)
	if err != nil {
		if pg.State.IsNoRows(err) {
			return nil, nil
		}

		return nil, err
	}

	return &state, nil
}
