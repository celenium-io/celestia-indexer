package indexer

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type Indexer struct {
	cfg Config
	wg  *sync.WaitGroup
	log zerolog.Logger
}

func New(cfg Config) *Indexer {
	return &Indexer{
		cfg: cfg,
		wg:  new(sync.WaitGroup),
		log: log.With().Str("module", "indexer").Logger(),
	}
}

func (i *Indexer) Start(ctx context.Context) error {
	i.log.Info().Msg("starting indexer...")
	return nil
}

func (i *Indexer) Stop() error {
	i.log.Info().Msg("stopping indexer...")
	i.wg.Wait()

	return nil
}
