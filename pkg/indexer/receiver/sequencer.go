package receiver

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func (r *Receiver) sequencer(ctx context.Context) {
	defer r.wg.Done()

	orderedBlocks := map[int64]types.BlockData{}
	currentBlock := r.cfg.Indexer.StartLevel

	for {
		select {
		case <-ctx.Done():
			return
		case block := <-r.blocks:
			orderedBlocks[block.Block.Height] = block

			if currentBlock == 0 {
				if err := r.receiveGenesis(ctx); err != nil {
					return
					// TODO: handle error on getting genesis, stop indexer
				}

				currentBlock += 1
				break
			}

			if b, ok := orderedBlocks[currentBlock]; ok {
				r.outputs[BlocksOutput].Push(b)

				r.log.Debug().Msgf("put in order block=%d", currentBlock)
				delete(orderedBlocks, currentBlock)
				currentBlock += 1
			} else {
				break
			}
		}
	}
}
