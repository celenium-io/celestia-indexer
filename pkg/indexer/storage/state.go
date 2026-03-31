// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func updateState(block *storage.Block, totalAccounts, totalNamespaces, totalProposals, ibcClientsCount int64, totalValidators int, version uint64, state *storage.State) {
	state.LastHeight = block.Height
	state.LastHash = block.Hash
	state.LastTime = block.Time
	state.TotalTx += block.Stats.TxCount
	state.TotalAccounts += totalAccounts
	state.TotalNamespaces += totalNamespaces
	state.TotalProposals += totalProposals
	state.TotalBlobsSize += block.Stats.BlobsSize
	state.TotalValidators += totalValidators
	state.TotalFee = state.TotalFee.Add(block.Stats.Fee.Decimal)
	state.TotalSupply = state.TotalSupply.Add(block.Stats.SupplyChange.Decimal)
	state.TotalIbcClients += ibcClientsCount
	state.ChainId = block.ChainId
	state.Version = version
}
