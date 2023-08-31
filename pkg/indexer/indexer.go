package indexer

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/parser"
	"github.com/dipdup-io/celestia-indexer/pkg/node"
	"github.com/dipdup-io/celestia-indexer/pkg/node/rpc"
	"github.com/dipdup-io/celestia-indexer/pkg/storage"
	"github.com/pkg/errors"
	"sync"

	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/receiver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Indexer struct {
	cfg      config.Config
	api      node.API
	receiver *receiver.Receiver
	parser   *parser.Parser
	storage  *storage.Module
	wg       *sync.WaitGroup
	log      zerolog.Logger
}

func New(ctx context.Context, cfg config.Config) (Indexer, error) {

	api := rpc.NewAPI(cfg.DataSources["node_rpc"])
	r := receiver.NewModule(cfg, &api)

	p := parser.NewModule()
	pInput, err := p.Input(parser.BlocksInput)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "cannot find input in parser")
	}
	if err = r.AttachTo(receiver.BlocksOutput, pInput); err != nil {
		return Indexer{}, err
	}

	pg, err := postgres.Create(ctx, cfg.Database)
	if err != nil {
		log.Err(err).Msg("creating pg context in indexer")
	}

	s := storage.NewModule(pg)
	sInput, err := s.Input(storage.InputName)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "cannot find input in storage")
	}
	if err = p.AttachTo(parser.DataOutput, sInput); err != nil {
		return Indexer{}, err
	}

	return Indexer{
		cfg:      cfg,
		api:      &api,
		receiver: &r,
		parser:   &p,
		storage:  &s,
		wg:       new(sync.WaitGroup),
		log:      log.With().Str("module", "indexer").Logger(),
	}, nil
}

func (i *Indexer) Start(ctx context.Context) {
	i.log.Info().Msg("starting...")

	go i.storage.Start(ctx)
	go i.parser.Start(ctx)
	go i.receiver.Start(ctx)
}

func (i *Indexer) Close() error {
	i.log.Info().Msg("closing...")
	i.wg.Wait()

	if err := i.receiver.Close(); err != nil {
		log.Err(err).Msg("closing receiver")
	}

	return nil
}
