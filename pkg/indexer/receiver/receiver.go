package receiver

import (
	"context"
	"sync"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/node"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/rs/zerolog/log"
)

const (
	BlocksOutput     = "blocks"
	RollbackOutput   = "signal"
	RollbackInput    = "state"
	GenesisOutput    = "genesis"
	GenesisDoneInput = "genesis_done"
	StopOutput       = "stop"
)

// Module - runs through a chain with aim ti catch-up head and identifies either block is fits in sequence or signals of rollback.
//
//		|----------------|
//		|                | -- types.BlockData -> BlocksOutput
//		|     MODULE     |
//		|    Receiver    | -- struct{}        -> RollbackOutput
//		|                | <- storage.State   -- RollbackInput
//	    |----------------|
type Module struct {
	modules.BaseModule
	api              node.API
	cfg              config.Indexer
	pool             *workerpool.Pool[types.Level]
	blocks           chan types.BlockData
	level            types.Level
	hash             []byte
	needGenesis      bool
	taskQueue        *sdkSync.Map[types.Level, struct{}]
	mx               *sync.RWMutex
	rollbackSync     *sync.WaitGroup
	cancelWorkers    context.CancelFunc
	cancelReadBlocks context.CancelFunc
}

var _ modules.Module = (*Module)(nil)

func NewModule(cfg config.Indexer, api node.API, state *storage.State) Module {
	level := types.Level(cfg.StartLevel)
	var lastHash []byte
	if state != nil {
		level = state.LastHeight
		lastHash = state.LastHash
	}

	receiver := Module{
		BaseModule:   modules.New("receiver"),
		api:          api,
		cfg:          cfg,
		blocks:       make(chan types.BlockData, cfg.ThreadsCount*10),
		needGenesis:  state == nil,
		level:        level,
		hash:         lastHash,
		taskQueue:    sdkSync.NewMap[types.Level, struct{}](),
		mx:           new(sync.RWMutex),
		rollbackSync: new(sync.WaitGroup),
	}

	receiver.CreateInput(RollbackInput)
	receiver.CreateInput(GenesisDoneInput)

	receiver.CreateOutput(BlocksOutput)
	receiver.CreateOutput(RollbackOutput)
	receiver.CreateOutput(GenesisOutput)
	receiver.CreateOutput(StopOutput)

	receiver.pool = workerpool.NewPool(receiver.worker, int(cfg.ThreadsCount))

	return receiver
}

func (r *Module) Start(ctx context.Context) {
	r.Log.Info().Msg("starting receiver...")
	workersCtx, cancelWorkers := context.WithCancel(ctx)
	r.cancelWorkers = cancelWorkers
	r.pool.Start(workersCtx)

	if r.needGenesis {
		if err := r.receiveGenesis(ctx); err != nil {
			log.Err(err).Msg("receiving genesis error")
			return
		}
	}

	r.G.GoCtx(ctx, r.sequencer)
	r.G.GoCtx(ctx, r.sync)
	r.G.GoCtx(ctx, r.rollback)
}

func (r *Module) Close() error {
	r.Log.Info().Msg("closing...")
	r.G.Wait()

	if err := r.pool.Close(); err != nil {
		return err
	}

	close(r.blocks)

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
	rollbackInput := r.MustInput(RollbackInput)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-rollbackInput.Listen():
			if !ok {
				r.Log.Warn().Msg("can't read message from rollback input, channel is closed and drained")
				continue
			}

			state, ok := msg.(storage.State)
			if !ok {
				r.Log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			r.taskQueue.Clear()
			r.setLevel(state.LastHeight, state.LastHash)
			r.Log.Info().Msgf("caught return from rollback to level=%d", state.LastHeight)
			r.rollbackSync.Done()
		}
	}
}
