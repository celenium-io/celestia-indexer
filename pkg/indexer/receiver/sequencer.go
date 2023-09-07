package receiver

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func (r *Module) sequencer(ctx context.Context) {
	orderedBlocks := map[int64]types.BlockData{}
	currentBlock := int64(r.level)

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
				prevB, ok := orderedBlocks[currentBlock-1]

				if ok {
					if !bytes.Equal(b.Block.LastBlockID.Hash, prevB.BlockID.Hash) {
						r.log.Info().
							Str("current.lastBlockHash", hex.EncodeToString(b.Block.LastBlockID.Hash)).
							Str("prevBlockHash", hex.EncodeToString(prevB.BlockID.Hash)).
							Uint64("level", uint64(b.Height)).
							Msg("rollback detected")
						// TODO	call rollback to the rescue and wait
						break
					}
				} // TODO else: check with block from storage?

				r.outputs[BlocksOutput].Push(b)
				r.setLevel(types.Level(currentBlock), b.BlockID.Hash)
				r.log.Debug().Msgf("put in order block=%d", currentBlock)

				currentBlock += 1
			} else {
				break
			}
		}
	}
}
