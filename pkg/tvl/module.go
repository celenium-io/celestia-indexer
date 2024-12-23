// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package tvl

import (
	"context"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/tvl/l2beat"
	"github.com/celenium-io/celestia-indexer/internal/tvl/lama"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	strg "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

const (
	rollupLimit = 100
)

var (
	syncTimestamp = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
)

type Module struct {
	modules.BaseModule
	l2beatApi l2beat.IApi
	lamaApi   lama.IApi
	rollup    storage.IRollup
	tvl       storage.ITvl
	log       zerolog.Logger
}

func New(l2beatCfg config.DataSource, lamaCfg config.DataSource, rollup storage.IRollup, tvl storage.ITvl) *Module {
	module := Module{
		BaseModule: modules.New("tvl"),
		rollup:     rollup,
		tvl:        tvl,
		l2beatApi:  l2beat.NewAPI(l2beatCfg),
		lamaApi:    lama.NewAPI(lamaCfg),
	}
	return &module
}

func (m *Module) Close() error {
	m.Log.Info().Msg("closing TVL scanner...")
	m.G.Wait()

	return nil
}

func (m *Module) Start(ctx context.Context) {
	m.log.Info().Msg("starting TVL scanner...")
	m.G.GoCtx(ctx, m.receive)
}

func (m *Module) getTvl(ctx context.Context, timeframe storage.TvlTimeframe) {
	rollups, err := m.rollup.List(ctx, rollupLimit, 0, strg.SortOrderAsc)
	if err != nil {
		m.Log.Err(err).Msg("receiving rollups")
		return
	}

	for i := range rollups {
		if len(rollups[i].L2Beat) > 0 {
			url := rollups[i].L2Beat
			lastIndex := strings.LastIndex(url, "/")

			if lastIndex == -1 {
				continue
			}

			rollupProject := url[lastIndex+1:]
			tvl, err := m.rollupTvlFromL2Beat(ctx, rollupProject, timeframe)
			if err != nil {
				m.Log.Err(err).Msg("receiving TVL from L2Beat")
				continue
			}

			tvlResponse := tvl[0].Result.Data.Json
			tvlModels := make([]*storage.Tvl, 0)
			for _, t := range tvlResponse {
				rollupTvl := t[1].(float64) + t[2].(float64) + t[3].(float64)
				tvlTs := time.Unix(int64(t[0].(float64)), 0)
				if tvlTs.After(syncTimestamp) {
					tvlModels = append(tvlModels, &storage.Tvl{
						Value:    rollupTvl,
						Time:     tvlTs,
						Rollup:   rollups[i],
						RollupId: rollups[i].Id,
					})
				}
			}

			if len(tvlModels) == 0 {
				continue
			}

			if err := m.tvl.SaveBulk(ctx, tvlModels...); err != nil {
				m.Log.Err(err).Msg("saving tvls")
			}
		}

		if len(rollups[i].DeFiLama) > 0 {
			tvl, err := m.rollupTvlFromLama(ctx, rollups[i].DeFiLama)
			if err != nil {
				m.Log.Err(err).Msg("receiving TVL from DeFi Lama")
				continue
			}

			tvlModels := make([]*storage.Tvl, 0)
			for _, t := range tvl {
				tvlTs := time.Unix(t.Date, 0)
				if tvlTs.After(syncTimestamp) {
					tvlModels = append(tvlModels, &storage.Tvl{
						Value:    t.TVL,
						Time:     tvlTs,
						Rollup:   rollups[i],
						RollupId: rollups[i].Id,
					})
				}
			}

			if len(tvlModels) == 0 {
				continue
			}

			if err := m.tvl.SaveBulk(ctx, tvlModels...); err != nil {
				m.Log.Err(err).Msg("saving tvls")
			}
		}
	}
	syncTimestamp, err = m.lastSyncTimeTvl(ctx)
	if err != nil {
		m.Log.Err(err).Msg("receiving last sync time for TVL")
	}
}

func (m *Module) receive(ctx context.Context) {
	syncTime, err := m.lastSyncTimeTvl(ctx)
	if err != nil {
		m.Log.Err(err).Msg("receiving last sync time for TVL")
		return
	}

	syncTimestamp = syncTime
	m.getTvl(ctx, storage.TvlTimeframeMax)
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.getTvl(ctx, storage.TvlTimeframe6Month)
		}
	}
}

func (m *Module) rollupTvlFromLama(ctx context.Context, rollupName string) ([]lama.TVLResponse, error) {
	requestTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return m.lamaApi.TVL(requestTimeout, rollupName)
}

func (m *Module) rollupTvlFromL2Beat(ctx context.Context, rollupName string, timeframe storage.TvlTimeframe) (l2beat.TVLResponse, error) {
	requestTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return m.l2beatApi.TVL(requestTimeout, rollupName, timeframe)
}

func (m *Module) lastSyncTimeTvl(ctx context.Context) (time.Time, error) {
	requestTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return m.tvl.LastSyncTime(requestTimeout)
}
