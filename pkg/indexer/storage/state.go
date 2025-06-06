// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/shopspring/decimal"
)

func updateState(block *storage.Block, totalAccounts, totalNamespaces, totalProposals, ibcClientsCount int64, totalValidators int, totalVotingPower decimal.Decimal, state *storage.State) {
	if block.Id <= uint64(state.LastHeight) {
		return
	}

	state.LastHeight = block.Height
	state.LastHash = block.Hash
	state.LastTime = block.Time
	state.TotalTx += block.Stats.TxCount
	state.TotalAccounts += totalAccounts
	state.TotalNamespaces += totalNamespaces
	state.TotalProposals += totalProposals
	state.TotalBlobsSize += block.Stats.BlobsSize
	state.TotalValidators += totalValidators
	state.TotalFee = state.TotalFee.Add(block.Stats.Fee)
	state.TotalSupply = state.TotalSupply.Add(block.Stats.SupplyChange)
	state.TotalStake = state.TotalStake.Add(totalVotingPower)
	state.TotalIbcClients += ibcClientsCount
	state.ChainId = block.ChainId
}
