package receiver

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/node"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

const (
	name         = "receiver"
	BlocksOutput = "blocks"
)

type Receiver struct {
	api     node.API
	cfg     config.Config
	outputs map[string]*modules.Output
	pool    *workerpool.Pool[storage.Level]
	blocks  chan types.BlockData
	log     zerolog.Logger
	wg      *sync.WaitGroup
}

func NewModule(cfg config.Config, api node.API) Receiver {
	receiver := Receiver{
		api:     api,
		cfg:     cfg,
		outputs: map[string]*modules.Output{BlocksOutput: modules.NewOutput(BlocksOutput)},
		blocks:  make(chan types.BlockData, cfg.Indexer.ThreadsCount*10),
		log:     log.With().Str("module", name).Logger(),
		wg:      new(sync.WaitGroup),
	}

	receiver.pool = workerpool.NewPool(receiver.worker, int(cfg.Indexer.ThreadsCount))

	return receiver
}

// Name -
func (*Receiver) Name() string {
	return name
}

func (r *Receiver) Start(ctx context.Context) {
	r.log.Info().Msg("starting receiver...")
	r.pool.Start(ctx)

	r.wg.Add(1)
	go r.sequencer(ctx)

	if err := r.readBlocks(ctx); err != nil {
		r.log.Err(err).Msg("read block")
		return
	}
}

func (r *Receiver) Close() error {
	r.log.Info().Msg("closing...")
	r.wg.Wait()

	if err := r.pool.Close(); err != nil {
		return err
	}

	close(r.blocks)

	return nil
}

func (r *Receiver) Output(name string) (*modules.Output, error) {
	output, ok := r.outputs[name]
	if !ok {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return output, nil
}

func (r *Receiver) Input(name string) (*modules.Input, error) {
	return nil, errors.Wrap(modules.ErrUnknownInput, name)
}

func (r *Receiver) AttachTo(outputName string, input *modules.Input) error {
	output, err := r.Output(outputName)
	if err != nil {
		return err
	}

	output.Attach(input)
	return nil
}
