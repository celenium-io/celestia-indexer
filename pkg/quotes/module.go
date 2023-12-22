// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package quotes

import (
	"context"
	"database/sql"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/binance"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
)

const (
	PricesOutput = "prices"

	symbol   = "TIAUSDT"
	interval = "1m"
)

var (
	startOfTime = time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
)

type Module struct {
	modules.BaseModule
	api         binance.IApi
	storage     storage.IPrice
	currentTime time.Time
}

func New(cfg config.DataSource, storage storage.IPrice) *Module {
	module := Module{
		BaseModule:  modules.New("quotes"),
		storage:     storage,
		api:         binance.NewAPI(cfg),
		currentTime: startOfTime,
	}
	module.CreateOutput(PricesOutput)
	return &module
}

func (m *Module) Start(ctx context.Context) {
	m.Log.Info().Msg("starting receiver...")

	if err := m.init(ctx); err != nil {
		m.Log.Err(err).Msg("initialization")
		return
	}

	m.G.GoCtx(ctx, m.receive)
}

func (m *Module) init(ctx context.Context) error {
	last, err := m.storage.Last(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return errors.Wrap(err, "get prices from database")
	}
	m.currentTime = last.Time
	return nil
}

func (m *Module) Close() error {
	m.Log.Info().Msg("closing...")
	m.G.Wait()

	return nil
}

func (m *Module) receive(ctx context.Context) {
	if err := m.get(ctx); err != nil {
		m.Log.Err(err).Msg("receiving prices")
		return
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.get(ctx); err != nil {
				m.Log.Err(err).Msg("receiving prices")
			}
		}
	}
}

func (m *Module) getPrices(ctx context.Context) (bool, error) {
	requestCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	start := m.currentTime.Add(time.Minute).UnixMilli()
	candles, err := m.api.OHLC(requestCtx, symbol, interval, &binance.OHLCArgs{
		Start: start,
	})
	if err != nil {
		return true, err
	}

	for i := range candles {
		if err := m.storage.Save(ctx, &storage.Price{
			Time:  candles[i].Time,
			Open:  candles[i].Open,
			High:  candles[i].High,
			Low:   candles[i].Low,
			Close: candles[i].Close,
		}); err != nil {
			return true, errors.Wrap(err, "saving price")
		}
		m.currentTime = candles[i].Time
	}

	end := len(candles) == 0
	if !end {
		m.Log.Info().Str("current_time", m.currentTime.String()).Msg("received quotes")
	}
	return end, nil
}

func (m *Module) get(ctx context.Context) error {
	var (
		end bool
		err error
	)
	for !end {
		end, err = m.getPrices(ctx)
		if err != nil {
			return err
		}
		if !end {
			time.Sleep(time.Millisecond * 500)
		}
	}
	return nil
}
