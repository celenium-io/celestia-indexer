// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
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
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// StatsTestSuite -
type StatsTestSuite struct {
	suite.Suite
	stats   *mock.MockIStats
	ns      *mock.MockINamespace
	price   *mock.MockIPrice
	state   *mock.MockIState
	echo    *echo.Echo
	handler StatsHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *StatsTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.stats = mock.NewMockIStats(s.ctrl)
	s.price = mock.NewMockIPrice(s.ctrl)
	s.ns = mock.NewMockINamespace(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.handler = NewStatsHandler(s.stats, s.ns, s.price, s.state)
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

func (s *StatsTestSuite) TestHistogramCountBlocks() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/summary/:table/:function/:timeframe")
	c.SetParamNames("table", "function", "timeframe")
	c.SetParamValues("block", "count", "day")

	s.stats.EXPECT().
		HistogramCount(gomock.Any(), storage.HistogramCountRequest{
			CountRequest: storage.CountRequest{
				Table: "block",
			},
			Timeframe: "day",
		}).
		Return([]storage.HistogramItem{
			{
				Value: "123123",
				Time:  testTime,
			},
		}, nil)

	s.Require().NoError(s.handler.Histogram(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.HistogramItem
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	item := response[0]
	s.Require().Equal("123123", item.Value)
	s.Require().True(testTime.Equal(item.Time))
}

func (s *StatsTestSuite) TestHistogramSumFeeBlocks() {
	q := make(url.Values)
	q.Set("column", "fee")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/summary/:table/:function/:timeframe")
	c.SetParamNames("table", "function", "timeframe")
	c.SetParamValues("block", "sum", "day")

	s.stats.EXPECT().
		Histogram(gomock.Any(), storage.HistogramRequest{
			SummaryRequest: storage.SummaryRequest{
				CountRequest: storage.CountRequest{
					Table: "block",
				},
				Function: "sum",
				Column:   "fee",
			},
			Timeframe: "day",
		}).
		Return([]storage.HistogramItem{
			{
				Value: "123123",
				Time:  testTime,
			},
		}, nil)

	s.Require().NoError(s.handler.Histogram(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.HistogramItem
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	item := response[0]
	s.Require().Equal("123123", item.Value)
	s.Require().True(testTime.Equal(item.Time))
}

func (s *StatsTestSuite) TestHistogramCountBlocksBadRequest() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/histogram/:table/:function/:timeframe")
	c.SetParamNames("table", "function", "timeframe")
	c.SetParamValues("unknown", "count", "day")

	s.Require().NoError(s.handler.Histogram(c))
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

func (s *StatsTestSuite) TestTxCount24h() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/tx_count_24h")

	s.stats.EXPECT().
		TxCountForLast24h(gomock.Any()).
		Return([]storage.TxCountForLast24hItem{
			{
				Time:    testTime,
				TxCount: 100,
				TPS:     0.01,
			},
		}, nil)

	s.Require().NoError(s.handler.TxCountHourly24h(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.TxCountHistogramItem
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	item := response[0]
	s.Require().EqualValues(100, item.Count)
	s.Require().EqualValues(0.01, item.TPS)
	s.Require().True(testTime.Equal(item.Time))
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

func (s *StatsTestSuite) TestPriceSeries() {
	for _, tf := range []string{
		storage.PriceTimeframeDay,
		storage.PriceTimeframeHour,
		storage.PriceTimeframeMinute,
	} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)
		c.SetPath("/v1/stats/price/series/:timeframe")
		c.SetParamNames("timeframe")
		c.SetParamValues(tf)

		s.price.EXPECT().
			Get(gomock.Any(), tf, int64(0), int64(0), 100).
			Return([]storage.Price{
				{
					Time:  testTime,
					Open:  decimal.RequireFromString("0.1"),
					High:  decimal.RequireFromString("0.2"),
					Low:   decimal.RequireFromString("0.01"),
					Close: decimal.RequireFromString("0.15"),
				},
			}, nil)

		s.Require().NoError(s.handler.PriceSeries(c))
		s.Require().Equal(http.StatusOK, rec.Code)

		var response []responses.Price
		err := json.NewDecoder(rec.Body).Decode(&response)
		s.Require().NoError(err)
		s.Require().Len(response, 1)

		item := response[0]
		s.Require().Equal("0.1", item.Open)
		s.Require().Equal("0.2", item.High)
		s.Require().Equal("0.01", item.Low)
		s.Require().Equal("0.15", item.Close)
	}
}

func (s *StatsTestSuite) TestPriceCurrent() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/price/current")

	s.price.EXPECT().
		Last(gomock.Any()).
		Return(storage.Price{
			Time:  testTime,
			Open:  decimal.RequireFromString("0.1"),
			High:  decimal.RequireFromString("0.2"),
			Low:   decimal.RequireFromString("0.01"),
			Close: decimal.RequireFromString("0.15"),
		}, nil)

	s.Require().NoError(s.handler.PriceCurrent(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.Price
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().Equal("0.1", response.Open)
	s.Require().Equal("0.2", response.High)
	s.Require().Equal("0.01", response.Low)
	s.Require().Equal("0.15", response.Close)
}
