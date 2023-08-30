package indexer

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/node"
	"github.com/dipdup-io/celestia-indexer/pkg/node/rpc"
	"sync"

	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/receiver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Indexer struct {
	cfg      config.Config
	api      node.API
	receiver *receiver.Receiver
	wg       *sync.WaitGroup
	log      zerolog.Logger
}

func New(cfg config.Config) *Indexer {

	api := rpc.NewAPI(cfg.DataSources["node_rpc"])

	return &Indexer{
		cfg:      cfg,
		api:      &api,
		receiver: receiver.New(cfg, &api),
		wg:       new(sync.WaitGroup),
		log:      log.With().Str("module", "indexer").Logger(),
	}
}

func (i *Indexer) Start(ctx context.Context) error {
	i.log.Info().Msg("starting indexer...")

	i.receiver.Start(ctx)

	return nil
}

func (i *Indexer) Close() error {
	i.log.Info().Msg("closing indexer...")
	i.wg.Wait()

	if err := i.receiver.Close(); err != nil {
		log.Err(err).Msg("closing receiver")
	}

	return nil
}
