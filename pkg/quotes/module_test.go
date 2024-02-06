// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package quotes

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/binance"
	binanceMock "github.com/celenium-io/celestia-indexer/internal/binance/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReceiver_get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prices := mock.NewMockIPrice(ctrl)
	api := binanceMock.NewMockIApi(ctrl)
	module := Module{
		BaseModule: modules.New("test"),
		api:        api,
		storage:    prices,
	}

	prices.EXPECT().
		Last(gomock.Any()).
		Return(storage.Price{
			Time: time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
		}, nil).
		MaxTimes(1).
		MinTimes(1)

	prices.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil).
		MaxTimes(3).
		MinTimes(3)

	api.EXPECT().
		OHLC(gomock.Any(), symbol, interval, &binance.OHLCArgs{
			Start: 1698710460000,
		}).
		Return([]binance.OHLC{
			{
				Time: time.Date(2023, 10, 31, 0, 1, 0, 0, time.UTC),
			},
			{
				Time: time.Date(2023, 10, 31, 0, 2, 0, 0, time.UTC),
			},
			{
				Time: time.Date(2023, 10, 31, 0, 3, 0, 0, time.UTC),
			},
		}, nil).
		MaxTimes(1).
		MinTimes(1)

	api.EXPECT().
		OHLC(gomock.Any(), symbol, interval, &binance.OHLCArgs{
			Start: 1698710640000,
		}).
		Return([]binance.OHLC{}, nil).
		MaxTimes(1).
		MinTimes(1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := module.init(ctx)
	require.NoError(t, err)

	err = module.get(ctx)
	require.NoError(t, err)

	current := time.Date(2023, 10, 31, 0, 3, 0, 0, time.UTC)
	require.Equal(t, current, module.currentTime)
}
