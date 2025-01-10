// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package tvl

import (
	"context"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/tvl/l2beat"
	l2beatMock "github.com/celenium-io/celestia-indexer/internal/tvl/l2beat/mock"
	"github.com/celenium-io/celestia-indexer/internal/tvl/lama"
	lamaMock "github.com/celenium-io/celestia-indexer/internal/tvl/lama/mock"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	strg "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
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

	l2beatTestResponse := make([][]interface{}, 2)
	for i := range l2beatTestResponse {
		l2beatTestResponse[i] = make([]interface{}, 4)
		for j := range l2beatTestResponse[i] {
			l2beatTestResponse[i][j] = float64(i * j)
		}
		l2beatTestResponse[i][0] = float64(time.Now().Unix())
	}

	l2BeatApi.EXPECT().
		TVL(gomock.Any(), gomock.Any(), storage.TvlTimeframeMonth).
		Return(l2beat.TVLResponse{{l2beat.Result{Data: l2beat.Data{Json: l2beatTestResponse}}}}, nil).
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
		err = module.save(ctx, testRollups[i], storage.TvlTimeframeMonth)
		require.NoError(t, err)
	}
}
