// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

import (
	"github.com/celenium-io/celestia-indexer/cmd/api/gas"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func blockProcessor(block storage.Block) Notification[*responses.Block] {
	response := responses.NewBlock(block, true)
	return NewBlockNotification(response)
}

func headProcessor(state storage.State) Notification[*responses.State] {
	response := responses.NewState(state)
	return NewStateNotification(response)
}

func gasPriceProcessor(data gas.GasPrice) Notification[*responses.GasPrice] {
	return NewGasPriceNotification(responses.GasPrice{
		Slow:   data.Slow,
		Median: data.Median,
		Fast:   data.Fast,
	})
}
