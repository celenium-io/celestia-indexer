// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package tvl

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/internal/tvl/l2beat"
	l2beatMock "github.com/celenium-io/celestia-indexer/internal/tvl/l2beat/mock"
	"github.com/celenium-io/celestia-indexer/internal/tvl/lama"
	lamaMock "github.com/celenium-io/celestia-indexer/internal/tvl/lama/mock"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	strg "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReceiver_tvl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l2BeatApi := l2beatMock.NewMockIApi(ctrl)
	lamaApi := lamaMock.NewMockIApi(ctrl)
	rollupMock := mock.NewMockIRollup(ctrl)
	tvlMock := mock.NewMockITvl(ctrl)
	module := Module{
		BaseModule: modules.New("test"),
		l2beatApi:  l2BeatApi,
		lamaApi:    lamaApi,
		rollup:     rollupMock,
		tvl:        tvlMock,
	}

	var testRollups = []*storage.Rollup{
		{Name: "Ham_test", Slug: "Ham_test_slug", DeFiLama: "Ham"},
		{Name: "Eclipse_test", Slug: "Eclipse_test_slug", L2Beat: "https://l2beat.com/bridges/projects/eclipse"},
	}
	rollupMock.EXPECT().
		List(gomock.Any(), rollupLimit, uint64(0), strg.SortOrderAsc).
		Return(testRollups, nil).
		Times(1)

	lamaApi.EXPECT().
		TVL(gomock.Any(), gomock.Any()).
		Return([]lama.TVLResponse{{Date: 1733529600, TVL: 12345678.90}}, nil).
		Times(1)

	l2beatTestResponse := make([]l2beat.Item, 2)
	for i := range l2beatTestResponse {
		l2beatTestResponse[i].Time = time.Now()
		l2beatTestResponse[i].Canonical = testsuite.RandomDecimal()
		l2beatTestResponse[i].External = testsuite.RandomDecimal()
		l2beatTestResponse[i].Native = testsuite.RandomDecimal()
		l2beatTestResponse[i].EthPrice = testsuite.RandomDecimal()
	}

	l2BeatApi.EXPECT().
		TVL(gomock.Any(), gomock.Any(), l2beat.TvlTimeframe30D).
		Return(l2beat.TVLResponse{
			Success: true,
			Data: l2beat.Data{
				Chart: l2beat.Chart{
					Data: l2beatTestResponse,
				},
			},
		}, nil).
		Times(1)

	tvlMock.EXPECT().
		SaveBulk(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := rollupMock.List(ctx, uint64(100), uint64(0), strg.SortOrderAsc)
	require.NoError(t, err)
	require.Equal(t, testRollups, result)

	for i := range testRollups {
		err = module.save(ctx, testRollups[i], l2beat.TvlTimeframe30D)
		require.NoError(t, err)
	}
}
