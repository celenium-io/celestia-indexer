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
	name             = "receiver"
	BlocksOutput     = "blocks"
	RollbackOutput   = "signal"
	RollbackInput    = "state"
	GenesisOutput    = "genesis"
	GenesisDoneInput = "genesis_done"
)

// Module - runs through chain with aim ti catch up head and identifies either block is fits in sequence or signals of rollback.
//
//		|----------------|
//		|                | -- types.BlockData -> BlocksOutput
//		|     MODULE     |
//		|    Receiver    | -- struct{}        -> RollbackOutput
//		|                | <- storage.State   -- RollbackInput
//	    |----------------|
type Module struct {
	api              node.API
	cfg              config.Indexer
	inputs           map[string]*modules.Input
	outputs          map[string]*modules.Output
	pool             *workerpool.Pool[types.Level]
	blocks           chan types.BlockData
	level            types.Level
	hash             []byte
	needGenesis      bool
	mx               *sync.RWMutex
	log              zerolog.Logger
	rollbackSync     *sync.WaitGroup
	g                workerpool.Group
	cancelWorkers    context.CancelFunc
	cancelReadBlocks context.CancelFunc
}

func NewModule(cfg config.Indexer, api node.API, state *storage.State) Module {
	level := types.Level(cfg.StartLevel)
	var lastHash []byte
	if state != nil {
		level = state.LastHeight
		lastHash = state.LastHash
	}

	receiver := Module{
		api: api,
		cfg: cfg,
		inputs: map[string]*modules.Input{
			RollbackInput:    modules.NewInput(RollbackInput),
			GenesisDoneInput: modules.NewInput(GenesisDoneInput),
		},
		outputs: map[string]*modules.Output{
			BlocksOutput:   modules.NewOutput(BlocksOutput),
			RollbackOutput: modules.NewOutput(RollbackOutput),
			GenesisOutput:  modules.NewOutput(GenesisOutput),
		},
		blocks:       make(chan types.BlockData, cfg.ThreadsCount*10),
		needGenesis:  state == nil,
		level:        level,
		hash:         lastHash,
		mx:           new(sync.RWMutex),
		log:          log.With().Str("module", name).Logger(),
		rollbackSync: new(sync.WaitGroup),
		g:            workerpool.NewGroup(),
	}

	receiver.pool = workerpool.NewPool(receiver.worker, int(cfg.ThreadsCount))

	return receiver
}

// Name -
func (*Module) Name() string {
	return name
}

func (r *Module) Start(ctx context.Context) {
	r.log.Info().Msg("starting receiver...")
	workersCtx, cancelWorkers := context.WithCancel(ctx)
	r.cancelWorkers = cancelWorkers
	r.pool.Start(workersCtx)

	if r.needGenesis {
		if err := r.receiveGenesis(ctx); err != nil {
			log.Err(err).Msg("receiving genesis error")
			return
		}
	}

	r.g.GoCtx(ctx, r.sequencer)
	r.g.GoCtx(ctx, r.sync)
	r.g.GoCtx(ctx, r.rollback)
}

func (r *Module) Close() error {
	r.log.Info().Msg("closing...")
	r.g.Wait()

	if err := r.pool.Close(); err != nil {
		return err
	}

	close(r.blocks)

	for name, input := range r.inputs {
		if err := input.Close(); err != nil {
			return errors.Wrapf(err, "closing error of '%s' input", name)
		}
	}

	return nil
}

func (r *Module) Output(name string) (*modules.Output, error) {
	output, ok := r.outputs[name]
	if !ok {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return output, nil
}

func (r *Module) Input(name string) (*modules.Input, error) {
	input, ok := r.inputs[name]
	if !ok {
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
	return input, nil
}

func (r *Module) AttachTo(outputName string, input *modules.Input) error {
	output, err := r.Output(outputName)
	if err != nil {
		return err
	}

	output.Attach(input)
	return nil
}

func (r *Module) Level() (types.Level, []byte) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	return r.level, r.hash
}

func (r *Module) setLevel(level types.Level, hash []byte) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.level = level
	r.hash = hash
}

func (r *Module) rollback(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-r.inputs[RollbackInput].Listen():
			if !ok {
				r.log.Warn().Msg("can't read message from rollback input")
				continue
			}

			state, ok := msg.(storage.State)
			if !ok {
				r.log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			r.setLevel(state.LastHeight, state.LastHash)
			r.log.Info().Msgf("caught return from rollback to level=%d", state.LastHeight)

			r.rollbackSync.Done()
		}
	}
}
