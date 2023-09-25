package receiver

import (
	"bytes"
	"context"
	"encoding/hex"

	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func (r *Module) sequencer(ctx context.Context) {
	orderedBlocks := map[int64]types.BlockData{}
	l, prevBlockHash := r.Level()
	currentBlock := int64(l + 1)

	for {
		select {
		case <-ctx.Done():
			return
		case block, ok := <-r.blocks:
			if !ok {
				r.Log.Warn().Msg("can't read message from blocks input, channel was dried and closed")
				r.stopAll()
				return
			}

			orderedBlocks[block.Block.Height] = block

			b, ok := orderedBlocks[currentBlock]
			for ok {
				if prevBlockHash != nil {
					if !bytes.Equal(b.Block.LastBlockID.Hash, prevBlockHash) {
						prevBlockHash, currentBlock, orderedBlocks = r.startRollback(ctx, b, prevBlockHash)
						break
					}
				}

				r.MustOutput(BlocksOutput).Push(b)
				r.setLevel(types.Level(currentBlock), b.BlockID.Hash)
				r.Log.Debug().
					Uint64("height", uint64(currentBlock)).
					Msg("put in order block")

				prevBlockHash = b.BlockID.Hash
				delete(orderedBlocks, currentBlock)
				currentBlock += 1

				b, ok = orderedBlocks[currentBlock]
			}
		}
	}
}

func (r *Module) startRollback(
	ctx context.Context,
	b types.BlockData,
	prevBlockHash []byte,
) ([]byte, int64, map[int64]types.BlockData) {
	r.Log.Info().
		Str("current.lastBlockHash", hex.EncodeToString(b.Block.LastBlockID.Hash)).
		Str("prevBlockHash", hex.EncodeToString(prevBlockHash)).
		Uint64("level", uint64(b.Height)).
		Msg("rollback detected")

	// Pause all receiver routines
	r.rollbackSync.Add(1)

	// Stop readBlocks
	if r.cancelReadBlocks != nil {
		r.cancelReadBlocks()
	}

	// Stop pool workers
	if r.cancelWorkers != nil {
		r.cancelWorkers()
	}

	clearChannel(r.blocks)

	// Start rollback
	r.MustOutput(RollbackOutput).Push(struct{}{})

	// Wait until rollback will be finished
	r.rollbackSync.Wait()

	// Reset empty state
	level, hash := r.Level()
	currentBlock := int64(level)
	prevBlockHash = hash
	orderedBlocks := map[int64]types.BlockData{}

	// Restart workers pool that read blocks
	workersCtx, cancelWorkers := context.WithCancel(ctx)
	r.cancelWorkers = cancelWorkers
	r.pool.Start(workersCtx)

	return prevBlockHash, currentBlock, orderedBlocks
}

func clearChannel(blocks <-chan types.BlockData) {
	for len(blocks) > 0 {
		<-blocks
	}
}
