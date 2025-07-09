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

	"github.com/celenium-io/celestia-indexer/cmd/api/gas"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var testGasState = gas.GasPrice{
	Slow:   "0.02",
	Median: "0.03",
	Fast:   "0.04",
}

// GasTestSuite -
type GasTestSuite struct {
	suite.Suite
	echo       *echo.Echo
	state      *mock.MockIState
	txs        *mock.MockITx
	constants  *mock.MockIConstant
	blockStats *mock.MockIBlockStats
	handler    GasHandler
	tracker    *gas.MockITracker
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *GasTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.state = mock.NewMockIState(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.constants = mock.NewMockIConstant(s.ctrl)
	s.blockStats = mock.NewMockIBlockStats(s.ctrl)
	s.tracker = gas.NewMockITracker(s.ctrl)
	s.handler = NewGasHandler(s.state, s.txs, s.constants, s.blockStats, s.tracker)
}

// TearDownSuite -
func (s *GasTestSuite) TearDownSuite() {
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

	s.constants.EXPECT().
		Get(gomock.Any(), types.ModuleNameAuth, "tx_size_cost_per_byte").
		Return(storage.Constant{
			Module: types.ModuleNameAuth,
			Name:   "tx_size_cost_per_byte",
			Value:  "10",
		}, nil).
		Times(1)

	s.constants.EXPECT().
		Get(gomock.Any(), types.ModuleNameBlob, "gas_per_blob_byte").
		Return(storage.Constant{
			Module: types.ModuleNameBlob,
			Name:   "gas_per_blob_byte",
			Value:  "8",
		}, nil).
		Times(1)

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

	s.tracker.EXPECT().
		State().
		Return(testGasState).
		Times(1)

	s.Require().NoError(s.handler.EstimatePrice(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.GasPrice
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Equal(testGasState.Slow, response.Slow)
	s.Require().Equal(testGasState.Median, response.Median)
	s.Require().Equal(testGasState.Fast, response.Fast)
}

func (s *GasTestSuite) TestEstimatePriceWithPriority() {
	for _, priority := range []string{"slow", "median", "fast"} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)
		c.SetPath("/gas/price/:priority")
		c.SetParamNames("priority")
		c.SetParamValues(priority)

		s.tracker.EXPECT().
			State().
			Return(testGasState).
			Times(1)

		s.Require().NoError(s.handler.EstimatePricePriority(c))
		s.Require().Equal(http.StatusOK, rec.Code, priority)

		var response string
		err := json.NewDecoder(rec.Body).Decode(&response)
		s.Require().NoError(err, priority)

		switch priority {
		case "slow":
			s.Require().Equal(testGasState.Slow, response, priority)
		case "median":
			s.Require().Equal(testGasState.Median, response, priority)
		case "fast":
			s.Require().Equal(testGasState.Fast, response, priority)
		}
	}
}
