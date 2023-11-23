// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import (
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func headProcessor(block storage.Block) *responses.Block {
	response := responses.NewBlock(block, true)
	return &response
}

func txProcessor(tx storage.Tx) *responses.Tx {
	response := responses.NewTx(tx)
	return &response
}
