// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package bus

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/workerpool"
)

type Observer struct {
	blocks chan *storage.Block
	txs    chan *storage.Tx

	listenBlocks bool
	listenTxs    bool

	g workerpool.Group
}

func NewObserver(channels ...string) *Observer {
	if len(channels) == 0 {
		return nil
	}

	observer := &Observer{
		blocks: make(chan *storage.Block, 1024),
		txs:    make(chan *storage.Tx, 1024),
		g:      workerpool.NewGroup(),
	}

	for i := range channels {
		switch channels[i] {
		case storage.ChannelHead:
			observer.listenBlocks = true
		case storage.ChannelTx:
			observer.listenTxs = true
		}
	}

	return observer
}

func (observer Observer) Close() error {
	observer.g.Wait()
	close(observer.blocks)
	close(observer.txs)
	return nil
}

func (observer Observer) notifyBlocks(block *storage.Block) {
	if observer.listenBlocks {
		observer.blocks <- block
	}
}

func (observer Observer) notifyTxs(tx *storage.Tx) {
	if observer.listenTxs {
		observer.txs <- tx
	}
}

func (observer Observer) Blocks() <-chan *storage.Block {
	return observer.blocks
}

func (observer Observer) Txs() <-chan *storage.Tx {
	return observer.txs
}
