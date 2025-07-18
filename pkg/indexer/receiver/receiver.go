// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/celenium-io/celestia-indexer/pkg/node"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/cometbft/cometbft/rpc/client/http"
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

// Module - runs through a chain aiming to catch up the head
//
//			and identifies whether the block fits in sequence or signals of rollback.
//
//		|----------------|
//		|                | -- types.BlockData -> BlocksOutput
//		|     MODULE     |
//		|    Receiver    | -- struct{}        -> RollbackOutput
//		|                | <- storage.State   -- RollbackInput
//	    |----------------|
type Module struct {
	modules.BaseModule
	api              node.Api
	cosmosApi        node.CosmosApi
	ws               *http.HTTP
	cfg              config.Indexer
	blocks           chan types.BlockData
	level            types.Level
	hash             []byte
	needGenesis      bool
	taskQueue        *sdkSync.Map[types.Level, struct{}]
	mx               *sync.RWMutex
	rollbackSync     *sync.WaitGroup
	cancelReadBlocks context.CancelFunc
	appVersion       *atomic.Uint64
	w                *Worker
}

var _ modules.Module = (*Module)(nil)

func NewModule(cfg config.Indexer, api node.Api, cosmosApi node.CosmosApi, ws *http.HTTP, state *storage.State) Module {
	level := types.Level(cfg.StartLevel)
	var lastHash []byte
	if state != nil {
		level = state.LastHeight
		lastHash = state.LastHash
	}

	receiver := Module{
		BaseModule:   modules.New("receiver"),
		api:          api,
		cosmosApi:    cosmosApi,
		ws:           ws,
		cfg:          cfg,
		blocks:       make(chan types.BlockData, 128),
		needGenesis:  state == nil,
		level:        level,
		hash:         lastHash,
		taskQueue:    sdkSync.NewMap[types.Level, struct{}](),
		mx:           new(sync.RWMutex),
		rollbackSync: new(sync.WaitGroup),
		appVersion:   new(atomic.Uint64),
	}

	receiver.w = NewWorker(api, receiver.Log, receiver.blocks, cfg.RequestBulkSize)

	receiver.CreateInput(RollbackInput)
	receiver.CreateInput(GenesisDoneInput)

	receiver.CreateOutput(BlocksOutput)
	receiver.CreateOutput(RollbackOutput)
	receiver.CreateOutput(GenesisOutput)
	receiver.CreateOutput(StopOutput)

	return receiver
}

func (r *Module) Start(ctx context.Context) {
	r.Log.Info().Msg("starting receiver...")

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

func (r *Module) stopAll() {
	r.MustOutput(StopOutput).Push(struct{}{})
}
