// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import (
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func blockProcessor(block storage.Block) *responses.Block {
	response := responses.NewBlock(block, true)
	return &response
}

func headProcessor(state storage.State) *responses.State {
	response := responses.NewState(state)
	return &response
}
