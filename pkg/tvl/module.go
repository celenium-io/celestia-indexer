// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package tvl

import (
	"context"
	"net/url"
	"strings"
	"time"

	"cosmossdk.io/errors"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/tvl/l2beat"
	"github.com/celenium-io/celestia-indexer/internal/tvl/lama"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	strg "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

const (
	rollupLimit = uint64(100)
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
	m.Log.Info().Msg("starting TVL scanner...")
	m.G.GoCtx(ctx, m.receive)
}

func (m *Module) getTvl(ctx context.Context, timeframe l2beat.TvlTimeframe) {
	rollups, err := m.rollup.List(ctx, rollupLimit, 0, strg.SortOrderAsc)
	if err != nil {
		m.Log.Err(err).Msg("receiving rollups")
		return
	}

	for i := range rollups {
		if err = m.save(ctx, rollups[i], timeframe); err != nil {
			m.Log.Err(err).Msg("saving TVL")
		}
	}

	syncTimestamp, err = m.lastSyncTimeTvl(ctx)
	if err != nil {
		m.Log.Err(err).Msg("receiving last sync time for TVL")
	}

	m.Log.Info().Msg("Sync rollup TVL is completed")
}

func (m *Module) save(ctx context.Context, rollup *storage.Rollup, timeframe l2beat.TvlTimeframe) error {
	if len(rollup.L2Beat) > 0 {
		if _, err := url.Parse(rollup.L2Beat); err != nil {
			return errors.Wrap(err, "invalid L2Beat url")
		}
		urlParts := strings.Split(rollup.L2Beat, "/")
		rollupProject := urlParts[len(urlParts)-1]

		tvl, err := m.rollupTvlFromL2Beat(ctx, rollupProject, timeframe)
		if err != nil {
			m.Log.Err(err).Msg("receiving TVL from L2Beat")
			return err
		}

		m.Log.Info().Str("rollup", rollup.Name).Msg("receiving TVL from L2Beat")

		tvlModels := make([]*storage.Tvl, 0)
		for _, t := range tvl.Data.Chart.Data {
			if t.Time.After(syncTimestamp) {
				tvlModels = append(tvlModels, &storage.Tvl{
					Value:    decimal.Sum(t.Canonical, t.External, t.Native),
					Time:     t.Time,
					Rollup:   rollup,
					RollupId: rollup.Id,
				})
			}
		}

		if len(tvlModels) == 0 {
			return nil
		}

		if err := m.tvl.SaveBulk(ctx, tvlModels...); err != nil {
			m.Log.Err(err).Msg("saving TVL")
			return err
		}

		m.Log.Info().Str("rollup", rollup.Name).Msg("successfully saving TVL")
		return nil
	}

	if len(rollup.DeFiLama) > 0 {
		tvl, err := m.rollupTvlFromLama(ctx, rollup.DeFiLama)
		if err != nil {
			m.Log.Err(err).Msg("receiving TVL from DeFi Lama")
			return err
		}

		m.Log.Info().Str("rollup", rollup.Name).Msg("receiving TVL from DeFi Lama")
		tvlModels := make([]*storage.Tvl, 0)
		for _, t := range tvl {
			tvlTs := time.Unix(t.Date, 0)
			if tvlTs.After(syncTimestamp) {
				tvlModels = append(tvlModels, &storage.Tvl{
					Value:    decimal.NewFromFloat(t.TVL),
					Time:     tvlTs,
					Rollup:   rollup,
					RollupId: rollup.Id,
				})
			}
		}

		if len(tvlModels) == 0 {
			return nil
		}

		if err := m.tvl.SaveBulk(ctx, tvlModels...); err != nil {
			m.Log.Err(err).Msg("saving TVL")
			return err
		}

		m.Log.Info().Str("rollup", rollup.Name).Msg("successfully saving TVL")
	}

	return nil
}

func (m *Module) receive(ctx context.Context) {
	syncTime, err := m.lastSyncTimeTvl(ctx)
	if err != nil {
		m.Log.Err(err).Msg("receiving last sync time for TVL")
	}

	syncTimestamp = syncTime
	m.getTvl(ctx, l2beat.TvlTimeframeMax)
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.getTvl(ctx, l2beat.TvlTimeframe180D)
		}
	}
}

func (m *Module) rollupTvlFromLama(ctx context.Context, rollupName string) ([]lama.TVLResponse, error) {
	requestTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return m.lamaApi.TVL(requestTimeout, rollupName)
}

func (m *Module) rollupTvlFromL2Beat(ctx context.Context, rollupName string, timeframe l2beat.TvlTimeframe) (l2beat.TVLResponse, error) {
	requestTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return m.l2beatApi.TVL(requestTimeout, rollupName, timeframe)
}

func (m *Module) lastSyncTimeTvl(ctx context.Context) (time.Time, error) {
	requestTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return m.tvl.LastSyncTime(requestTimeout)
}
