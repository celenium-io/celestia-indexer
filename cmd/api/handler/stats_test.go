// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/cmd/api/hyperlane"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// StatsTestSuite -
type StatsTestSuite struct {
	suite.Suite
	stats      *mock.MockIStats
	ns         *mock.MockINamespace
	state      *mock.MockIState
	ibc        *mock.MockIIbcTransfer
	hyperlane  *mock.MockIHLTransfer
	channels   *mock.MockIIbcChannel
	chainStore *hyperlane.MockIChainStore
	echo       *echo.Echo
	handler    StatsHandler
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *StatsTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.stats = mock.NewMockIStats(s.ctrl)
	s.ns = mock.NewMockINamespace(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.ibc = mock.NewMockIIbcTransfer(s.ctrl)
	s.hyperlane = mock.NewMockIHLTransfer(s.ctrl)
	s.channels = mock.NewMockIIbcChannel(s.ctrl)
	s.chainStore = hyperlane.NewMockIChainStore(s.ctrl)
	s.handler = NewStatsHandler(s.stats, s.ns, s.ibc, s.channels, s.hyperlane, s.chainStore, s.state)
}

// TearDownSuite -
func (s *StatsTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteStats_Run(t *testing.T) {
	suite.Run(t, new(StatsTestSuite))
}

func (s *StatsTestSuite) TestCountBlocks() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/summary/:table/:function")
	c.SetParamNames("table", "function")
	c.SetParamValues("block", "count")

	s.stats.EXPECT().
		Count(gomock.Any(), storage.CountRequest{
			Table: "block",
		}).
		Return("21000", nil)

	s.Require().NoError(s.handler.Summary(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response string
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Equal("21000", response)
}

func (s *StatsTestSuite) TestSumFeeBlocks() {
	q := make(url.Values)
	q.Set("column", "fee")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/summary/:table/:function")
	c.SetParamNames("table", "function")
	c.SetParamValues("block", "sum")

	s.stats.EXPECT().
		Summary(gomock.Any(), storage.SummaryRequest{
			CountRequest: storage.CountRequest{
				Table: "block",
			},
			Function: "sum",
			Column:   "fee",
		}).
		Return("21000", nil)

	s.Require().NoError(s.handler.Summary(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response string
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Equal("21000", response)
}

func (s *StatsTestSuite) TestCountBlocksBadRequest() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/summary/:table/:function")
	c.SetParamNames("table", "function")
	c.SetParamValues("unknown", "count")

	s.Require().NoError(s.handler.Summary(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)
}

func (s *StatsTestSuite) TestTPS() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/tps")

	s.stats.EXPECT().
		TPS(gomock.Any()).
		Return(storage.TPS{
			Current:           0.3,
			High:              1,
			Low:               0.1,
			ChangeLastHourPct: 0.12,
		}, nil)

	s.Require().NoError(s.handler.TPS(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.TPS
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().EqualValues(0.3, response.Current)
	s.Require().EqualValues(1, response.High)
	s.Require().EqualValues(0.1, response.Low)
	s.Require().EqualValues(0.12, response.ChangeLastHourPct)
}

func (s *StatsTestSuite) TestChanges24hBlockStats() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/changes_24h")

	s.stats.EXPECT().
		Change24hBlockStats(gomock.Any()).
		Return(storage.Change24hBlockStats{
			BytesInBlock: 0.123,
			BlobsSize:    0.321,
			Fee:          0.432,
			TxCount:      0.345,
		}, nil)

	s.Require().NoError(s.handler.Change24hBlockStats(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var item responses.Change24hBlockStats
	err := json.NewDecoder(rec.Body).Decode(&item)
	s.Require().NoError(err)

	s.Require().EqualValues(0.123, item.BytesInBlock)
	s.Require().EqualValues(0.321, item.BlobsSize)
	s.Require().EqualValues(0.432, item.Fee)
	s.Require().EqualValues(0.345, item.TxCount)
}

func (s *StatsTestSuite) TestNamespaceUsage() {
	q := make(url.Values)
	q.Set("top", "1")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/namespace/usage")

	s.ns.EXPECT().
		ListWithSort(gomock.Any(), "size", sdk.SortOrderDesc, 1, 0).
		Return([]storage.Namespace{
			testNamespace,
		}, nil)

	s.state.EXPECT().
		List(gomock.Any(), uint64(1), uint64(0), sdk.SortOrderAsc).
		Return([]*storage.State{
			&testState,
		}, nil)

	s.Require().NoError(s.handler.NamespaceUsage(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.NamespaceUsage
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 2)

	item0 := response[0]
	s.Require().Equal("fc7443b155920156", item0.Name)
	s.Require().Equal(testNamespace.Size, item0.Size)
	s.Require().Equal("0000000000000000000000000000000000000000fc7443b155920156", item0.NamespaceID)
	s.Require().NotNil(item0.Version)
	s.Require().EqualValues(0, *item0.Version)

	item1 := response[1]
	s.Require().Equal("others", item1.Name)
	s.Require().EqualValues(900, item1.Size)
	s.Require().Nil(item1.Version)
}

func (s *StatsTestSuite) TestBlockStatsHistogram() {
	for _, name := range []string{
		storage.SeriesBPS,
		storage.SeriesBlobsSize,
		storage.SeriesBlockTime,
		storage.SeriesEventsCount,
		storage.SeriesFee,
		storage.SeriesSupplyChange,
		storage.SeriesTPS,
		storage.SeriesTxCount,
		storage.SeriesGasEfficiency,
		storage.SeriesGasLimit,
		storage.SeriesGasPrice,
		storage.SeriesGasUsed,
		storage.SeriesBytesInBlock,
		storage.SeriesBlobsCount,
	} {

		for _, tf := range []storage.Timeframe{
			storage.TimeframeHour,
			storage.TimeframeDay,
			storage.TimeframeWeek,
			storage.TimeframeMonth,
			storage.TimeframeYear,
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/v1/stats/series/:name/:timeframe")
			c.SetParamNames("name", "timeframe")
			c.SetParamValues(name, string(tf))

			s.stats.EXPECT().
				Series(gomock.Any(), tf, name, gomock.Any()).
				Return([]storage.SeriesItem{
					{
						Time:  testTime,
						Value: "11234",
						Max:   "782634",
						Min:   "69.6665479793",
					},
				}, nil)

			s.Require().NoError(s.handler.Series(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var response []responses.SeriesItem
			err := json.NewDecoder(rec.Body).Decode(&response)
			s.Require().NoError(err)
			s.Require().Len(response, 1)

			item := response[0]
			s.Require().Equal("11234", item.Value)
			s.Require().Equal("782634", item.Max)
			s.Require().Equal("69.6665479793", item.Min)
		}
	}
}

func (s *StatsTestSuite) TestBlockCumulativeStatsHistogram() {
	for _, name := range []string{
		storage.SeriesBlobsSize,
		storage.SeriesFee,
		storage.SeriesTxCount,
		storage.SeriesGasLimit,
		storage.SeriesGasUsed,
		storage.SeriesBytesInBlock,
		storage.SeriesBlobsCount,
	} {

		for _, tf := range []storage.Timeframe{
			storage.TimeframeDay,
			storage.TimeframeWeek,
			storage.TimeframeMonth,
			storage.TimeframeYear,
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/v1/stats/series/:name/:timeframe/cumulative")
			c.SetParamNames("name", "timeframe")
			c.SetParamValues(name, string(tf))

			s.stats.EXPECT().
				Series(gomock.Any(), tf, name, gomock.Any()).
				Return([]storage.SeriesItem{
					{
						Time:  testTime,
						Value: "11234",
					},
				}, nil)

			s.Require().NoError(s.handler.Series(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var response []responses.SeriesItem
			err := json.NewDecoder(rec.Body).Decode(&response)
			s.Require().NoError(err)
			s.Require().Len(response, 1)

			item := response[0]
			s.Require().Equal("11234", item.Value)
		}
	}
}

func (s *StatsTestSuite) TestNamespaceStatsHistogram() {
	for _, name := range []string{
		storage.SeriesNsPfbCount,
		storage.SeriesNsSize,
	} {

		for _, tf := range []storage.Timeframe{
			storage.TimeframeHour,
			storage.TimeframeDay,
			storage.TimeframeWeek,
			storage.TimeframeMonth,
			storage.TimeframeYear,
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/v1/stats/namespace/series/:id/:name/:timeframe")
			c.SetParamNames("id", "name", "timeframe")
			c.SetParamValues("000000000000000000000000000000000000000008E5F679BF7116CB", name, string(tf))

			s.ns.EXPECT().
				ByNamespaceId(gomock.Any(), gomock.Any()).
				Return([]storage.Namespace{
					testNamespace,
				}, nil)

			s.stats.EXPECT().
				NamespaceSeries(gomock.Any(), tf, name, testNamespace.Id, gomock.Any()).
				Return([]storage.SeriesItem{
					{
						Time:  testTime,
						Value: "11234",
					},
				}, nil)

			s.Require().NoError(s.handler.NamespaceSeries(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var response []responses.SeriesItem
			err := json.NewDecoder(rec.Body).Decode(&response)
			s.Require().NoError(err)
			s.Require().Len(response, 1)

			item := response[0]
			s.Require().Equal("11234", item.Value)
		}
	}
}

func (s *StatsTestSuite) TestSquareSize() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/square_size")

	s.stats.EXPECT().
		SquareSize(gomock.Any(), nil, nil).
		Return(map[int][]storage.SeriesItem{
			2: {
				{
					Time:  testTime,
					Value: "100",
				},
			},
		}, nil)

	s.Require().NoError(s.handler.SquareSize(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.SquareSizeResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)
	s.Require().Contains(response, 2)
	s.Require().Len(response[2], 1)
}

func (s *StatsTestSuite) TestCumulativeSeries() {
	for _, name := range []string{
		storage.SeriesBlobsSize,
		storage.SeriesFee,
		storage.SeriesSupplyChange,
		storage.SeriesTxCount,
		storage.SeriesGasLimit,
		storage.SeriesGasUsed,
		storage.SeriesBytesInBlock,
		storage.SeriesBlobsCount,
	} {

		for _, tf := range []storage.Timeframe{
			storage.TimeframeHour,
			storage.TimeframeDay,
			storage.TimeframeWeek,
			storage.TimeframeMonth,
			storage.TimeframeYear,
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/v1/stats/series/:name/:timeframe")
			c.SetParamNames("name", "timeframe")
			c.SetParamValues(name, string(tf))

			s.stats.EXPECT().
				CumulativeSeries(gomock.Any(), tf, name, storage.NewSeriesRequest(0, 0)).
				Return([]storage.SeriesItem{
					{
						Value: "1000",
						Time:  testTime,
					},
				}, nil)

			s.Require().NoError(s.handler.SeriesCumulative(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var response []responses.SeriesItem
			err := json.NewDecoder(rec.Body).Decode(&response)
			s.Require().NoError(err)
			s.Require().Len(response, 1)

			item := response[0]
			s.Require().Equal("1000", item.Value)
		}
	}
}

func (s *StatsTestSuite) TestRollupStats24h() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/rollup_stats_24h")

	s.stats.EXPECT().
		RollupStats24h(gomock.Any()).
		Return([]storage.RollupStats24h{
			{
				Name:       "name",
				Logo:       "logo",
				RollupId:   1,
				Size:       12,
				Fee:        43.01,
				BlobsCount: 123,
			},
		}, nil)

	s.Require().NoError(s.handler.RollupStats24h(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.RollupStats24h
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	item := response[0]
	s.Require().EqualValues(1, item.Id)
	s.Require().EqualValues(12, item.Size)
	s.Require().EqualValues(43, item.Fee)
	s.Require().EqualValues(123, item.BlobsCount)
	s.Require().EqualValues("name", item.Name)
	s.Require().EqualValues("logo", item.Logo)
}

func (s *StatsTestSuite) TestMessgaesCount24h() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/messages_count_24h")

	s.stats.EXPECT().
		MessagesCount24h(gomock.Any()).
		Return([]storage.CountItem{
			{
				Name:  "test",
				Value: 100,
			},
		}, nil)

	s.Require().NoError(s.handler.MessagesCount24h(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.CountItem
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	s.Require().EqualValues("test", response[0].Name)
	s.Require().EqualValues(100, response[0].Value)
}

func (s *StatsTestSuite) TestSizeGroups() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/size_groups")

	s.stats.EXPECT().
		SizeGroups(gomock.Any(), nil).
		Return([]storage.SizeGroup{
			{
				Name:    "test",
				Size:    100,
				AvgSize: 10,
				Count:   10,
			},
		}, nil)

	s.Require().NoError(s.handler.SizeGroups(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.SizeGroup
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	s.Require().EqualValues("test", response[0].Name)
	s.Require().EqualValues(100, response[0].Size)
	s.Require().EqualValues(10, response[0].AvgSize)
	s.Require().EqualValues(10, response[0].Count)
}

func (s *StatsTestSuite) TestIbcSeries() {
	for _, name := range []string{
		"count",
		"amount",
	} {

		for _, tf := range []storage.Timeframe{
			storage.TimeframeHour,
			storage.TimeframeDay,
			storage.TimeframeMonth,
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/v1/stats/ibc/series/:id/:name/:timeframe")
			c.SetParamNames("id", "name", "timeframe")
			c.SetParamValues("channel-1", name, string(tf))

			s.ibc.EXPECT().
				Series(gomock.Any(), "channel-1", tf, name, storage.NewSeriesRequest(0, 0)).
				Return([]storage.HistogramItem{
					{
						Value: "1000",
						Time:  testTime,
					},
				}, nil)

			s.Require().NoError(s.handler.IbcSeries(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var response []responses.HistogramItem
			err := json.NewDecoder(rec.Body).Decode(&response)
			s.Require().NoError(err)
			s.Require().Len(response, 1)

			item := response[0]
			s.Require().Equal("1000", item.Value)
		}
	}
}

func (s *StatsTestSuite) TestIbcChainStats() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/ibc/chains")

	s.channels.EXPECT().
		StatsByChain(gomock.Any(), 10, 0).
		Return([]storage.ChainStats{
			{
				Chain:    "test",
				Received: decimal.RequireFromString("101"),
				Sent:     decimal.RequireFromString("99"),
				Flow:     decimal.RequireFromString("200"),
			},
		}, nil)

	s.Require().NoError(s.handler.IbcByChains(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.IbcChainStats
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	s.Require().EqualValues("test", response[0].Chain)
	s.Require().EqualValues("101", response[0].Received)
	s.Require().EqualValues("99", response[0].Sent)
	s.Require().EqualValues("200", response[0].Flow)
}

func (s *StatsTestSuite) TestIbcSummaryStats() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/ibc/summary")

	s.ibc.EXPECT().
		LargestTransfer24h(gomock.Any()).Return(
		storage.IbcTransfer{
			Id:              1,
			Time:            testTime,
			Height:          1000,
			Timeout:         &testTime,
			ChannelId:       "channel-1",
			ConnectionId:    "connection-1",
			Amount:          decimal.RequireFromString("101"),
			Denom:           currency.Utia,
			Memo:            "memo",
			ReceiverAddress: testsuite.Ptr("osmo1mj37s3mmv78tj0ke3yely7zwmzl5rkh9gx9ma2"),
			Sender: &storage.Address{
				Hash:    testHashAddress,
				Address: testAddress,
			},
			Sequence: 123456,
			Tx:       &testTx,
			Connection: &storage.IbcConnection{
				Client: &storage.IbcClient{
					ChainId: "chain-id",
				},
			},
		}, nil)

	s.channels.EXPECT().
		BusiestChannel1m(gomock.Any()).
		Return(storage.BusiestChannel{
			ChannelId:      "channel-111",
			TransfersCount: 1000,
			ChainId:        "chain-1",
		}, nil)

	s.Require().NoError(s.handler.IbcSummary(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.IbcSummaryStats
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().EqualValues(1, response.LargestTransfer.Id)
	s.Require().EqualValues("channel-1", response.LargestTransfer.ChannelId)
	s.Require().EqualValues("101", response.LargestTransfer.Amount)
	s.Require().EqualValues(currency.Utia, response.LargestTransfer.Denom)

	s.Require().EqualValues("channel-111", response.BusiestChannel.ChannelId)
	s.Require().EqualValues("chain-1", response.BusiestChannel.ChainId)
	s.Require().EqualValues(1000, response.BusiestChannel.TransfersCount)
}

func (s *StatsTestSuite) TestHlDomainStats() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/hyperlane/chains")

	s.chainStore.EXPECT().
		Set(testChainStore).
		Times(1)

	s.chainStore.EXPECT().
		All().
		Return(testChainStore).
		Times(1)

	s.chainStore.EXPECT().
		Get(gomock.Any()).
		Return(testChainMetadata, true).
		Times(len(testChainStore))

	s.hyperlane.EXPECT().
		StatsByDomain(gomock.Any(), 10, 0).
		Return([]storage.DomainStats{
			{
				Domain:  1,
				Amount:  decimal.RequireFromString("1010.101"),
				TxCount: 1488,
			},
		}, nil)

	s.chainStore.Set(testChainStore)
	s.chainStore.All()
	s.Require().NoError(s.handler.HlByDomain(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.HlDomainStats
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	result := response[0]
	s.Require().EqualValues(1, result.Domain)
	s.Require().EqualValues("1010.101", result.Amount)
	s.Require().EqualValues(1488, result.TransfersCount)
	s.Require().NotNil(result.ChainMetadata)
	s.Require().EqualValues(testChainMetadata.DisplayName, result.ChainMetadata.Name)
	s.Require().EqualValues(testChainMetadata.NativeToken.Decimals, result.ChainMetadata.NativeToken.Decimals)
	s.Require().EqualValues(testChainMetadata.NativeToken.Name, result.ChainMetadata.NativeToken.Name)
	s.Require().EqualValues(testChainMetadata.NativeToken.Symbol, result.ChainMetadata.NativeToken.Symbol)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].Name, result.ChainMetadata.BlockExplorers[0].Name)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].ApiUrl, result.ChainMetadata.BlockExplorers[0].ApiUrl)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].Url, result.ChainMetadata.BlockExplorers[0].Url)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].Family, result.ChainMetadata.BlockExplorers[0].Family)
}

func (s *StatsTestSuite) TestHlSeries() {
	for _, name := range []string{
		"count",
		"amount",
	} {

		for _, tf := range []storage.Timeframe{
			storage.TimeframeHour,
			storage.TimeframeDay,
			storage.TimeframeMonth,
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/v1/stats/hyperlane/series/:id/:name/:timeframe")
			c.SetParamNames("id", "name", "timeframe")
			c.SetParamValues("123", name, string(tf))

			s.hyperlane.EXPECT().
				Series(gomock.Any(), uint64(123), tf, name, storage.NewSeriesRequest(0, 0)).
				Return([]storage.HistogramItem{
					{
						Value: "1111",
						Time:  testTime,
					},
				}, nil)

			s.Require().NoError(s.handler.HlSeries(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var response []responses.HistogramItem
			err := json.NewDecoder(rec.Body).Decode(&response)
			s.Require().NoError(err)
			s.Require().Len(response, 1)

			item := response[0]
			s.Require().Equal("1111", item.Value)
		}
	}
}

func (s *StatsTestSuite) TestHlDomainStatsWithoutChainStore() {
	s.chainStore = nil
	s.handler = NewStatsHandler(s.stats, s.ns, s.ibc, s.channels, s.hyperlane, nil, s.state)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/hyperlane/chains")

	s.hyperlane.EXPECT().
		StatsByDomain(gomock.Any(), 10, 0).
		Return([]storage.DomainStats{
			{
				Domain:  1,
				Amount:  decimal.RequireFromString("1010.101"),
				TxCount: 1488,
			},
		}, nil)

	s.Require().NoError(s.handler.HlByDomain(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.HlDomainStats
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	result := response[0]
	s.Require().EqualValues(1, result.Domain)
	s.Require().EqualValues("1010.101", result.Amount)
	s.Require().EqualValues(1488, result.TransfersCount)
	s.Require().Nil(result.ChainMetadata)
}

func (s *StatsTestSuite) TestHlTotalSeries() {
	for _, name := range []string{
		"count",
		"amount",
	} {

		for _, tf := range []storage.Timeframe{
			storage.TimeframeHour,
			storage.TimeframeDay,
			storage.TimeframeMonth,
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/v1/stats/hyperlane/chains/:name/:timeframe")
			c.SetParamNames("name", "timeframe")
			c.SetParamValues(name, string(tf))

			s.hyperlane.EXPECT().
				TotalSeries(gomock.Any(), tf, name, storage.NewSeriesRequest(0, 0)).
				Return([]storage.HistogramItem{
					{
						Value: "1234",
						Time:  testTime,
					},
				}, nil)

			s.Require().NoError(s.handler.HlTotalSeries(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var response []responses.HistogramItem
			err := json.NewDecoder(rec.Body).Decode(&response)
			s.Require().NoError(err)
			s.Require().Len(response, 1)

			item := response[0]
			s.Require().Equal("1234", item.Value)
		}
	}
}

func (s *StatsTestSuite) TestStakingSeries() {
	for _, name := range []string{
		storage.SeriesDelegations,
		storage.SeriesDelegationsCount,
		storage.SeriesUnbondings,
		storage.SeriesUnbondingsCount,
		storage.SeriesCommissions,
		storage.SeriesRewards,
		storage.SeriesFlow,
		storage.SeriesCumulativeFlow,
	} {

		for _, tf := range []storage.Timeframe{
			storage.TimeframeHour,
			storage.TimeframeDay,
			storage.TimeframeMonth,
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/v1/stats/staking/series/:id/:name/:timeframe")
			c.SetParamNames("id", "name", "timeframe")
			c.SetParamValues("1", name, string(tf))

			s.stats.EXPECT().
				StakingSeries(gomock.Any(), tf, name, uint64(1), storage.NewSeriesRequest(0, 0)).
				Return([]storage.SeriesItem{
					{
						Value: "1111",
						Time:  testTime,
					},
				}, nil)

			s.Require().NoError(s.handler.StakingSeries(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var response []responses.SeriesItem
			err := json.NewDecoder(rec.Body).Decode(&response)
			s.Require().NoError(err)
			s.Require().Len(response, 1)

			item := response[0]
			s.Require().Equal("1111", item.Value)
		}
	}
}
