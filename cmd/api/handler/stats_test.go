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
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// StatsTestSuite -
type StatsTestSuite struct {
	suite.Suite
	stats   *mock.MockIStats
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
	s.handler = NewStatsHandler(s.stats)
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
