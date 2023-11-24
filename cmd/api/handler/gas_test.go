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

	"github.com/celenium-io/celestia-indexer/cmd/api/gas"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// GasTestSuite -
type GasTestSuite struct {
	suite.Suite
	echo       *echo.Echo
	state      *mock.MockIState
	txs        *mock.MockITx
	blockStats *mock.MockIBlockStats
	handler    GasHandler
	tracker    *gas.Tracker
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *GasTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.state = mock.NewMockIState(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.blockStats = mock.NewMockIBlockStats(s.ctrl)
	s.tracker = gas.NewTracker(s.state, s.blockStats, s.txs, nil)
	s.handler = NewGasHandler(s.state, s.txs, s.blockStats, s.tracker)
}

// TearDownSuite -
func (s *GasTestSuite) TearDownSuite() {
	if s.tracker != nil {
		err := s.tracker.Close()
		s.Require().NoError(err)
	}
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteGas_Run(t *testing.T) {
	suite.Run(t, new(GasTestSuite))
}

func (s *GasTestSuite) TestEstimateForPfb() {
	q := make(url.Values)
	q.Set("sizes", "12,34")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/gas/estimate_for_pfb")

	s.Require().NoError(s.handler.EstimateForPfb(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response uint64
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Greater(response, uint64(0))
}

func (s *GasTestSuite) TestEstimatePrice() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/gas/price")

	s.Require().NoError(s.handler.EstimatePrice(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.GasPrice
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Equal("0", response.Fast)
	s.Require().Equal("0", response.Slow)
	s.Require().Equal("0", response.Median)
}
