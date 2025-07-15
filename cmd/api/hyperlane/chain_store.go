// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package hyperlane

import (
	"context"
	"sync"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/node/hyperlane"
	"github.com/dipdup-io/workerpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Store *ChainStore

type ChainStore struct {
	data map[uint64]hyperlane.ChainMetadata
	api  hyperlane.IApi
	mx   *sync.RWMutex
	g    workerpool.Group
	log  zerolog.Logger
}

func NewChainStore(url string) *ChainStore {
	api := hyperlane.NewApi(
		url,
		hyperlane.WithRateLimit(1),
		hyperlane.WithTimeout(time.Second*time.Duration(1)))

	cs := &ChainStore{
		data: make(map[uint64]hyperlane.ChainMetadata),
		api:  api,
		mx:   new(sync.RWMutex),
		g:    workerpool.NewGroup(),
		log:  log.With().Str("module", "chain_store").Logger(),
	}

	Store = cs
	return cs
}

func (cs *ChainStore) Start(ctx context.Context) {
	cs.g.GoCtx(ctx, cs.sync)
}

func (cs *ChainStore) Get(domainId uint64) (hyperlane.ChainMetadata, bool) {
	cs.mx.RLock()
	defer cs.mx.RUnlock()
	val, ok := cs.data[domainId]
	return val, ok
}

func (cs *ChainStore) Set(metadata map[uint64]hyperlane.ChainMetadata) {
	cs.mx.Lock()
	cs.data = metadata
	cs.mx.Unlock()
}

func (cs *ChainStore) sync(ctx context.Context) {
	metadata, err := cs.api.ChainMetadata(ctx)
	if err != nil {
		cs.log.Error().Err(err).Msg("sync hyperlane chain metadata failed")
	}

	cs.log.Info().Int("chain metadata count", len(metadata)).Msg("sync hyperlane chain metadata")
	cs.Set(metadata)

	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ticker.C:
			if metadata, err = cs.api.ChainMetadata(ctx); err != nil {
				cs.log.Error().Err(err).Msg("sync hyperlane chain metadata failed")
			}

			cs.Set(metadata)
		}
	}
}

func (cs *ChainStore) Close() error {
	cs.g.Wait()

	return nil
}
