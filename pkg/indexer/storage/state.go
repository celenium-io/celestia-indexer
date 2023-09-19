package storage

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func updateState(block *storage.Block, totalAccounts uint64, totalNamespaces uint64, state *storage.State) {
	if types.Level(block.Id) <= state.LastHeight {
		return
	}

	state.LastHeight = block.Height
	state.LastHash = block.Hash
	state.LastTime = block.Time
	state.TotalTx += block.Stats.TxCount
	state.TotalAccounts += totalAccounts
	state.TotalNamespaces += totalNamespaces
	state.TotalBlobsSize += block.Stats.BlobsSize
	state.TotalFee = state.TotalFee.Add(block.Stats.Fee)
	state.TotalSupply = state.TotalSupply.Add(block.Stats.SupplyChange)
	state.ChainId = block.ChainId
}
