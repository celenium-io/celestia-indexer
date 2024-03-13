// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// VestingTestSuite -
type VestingTestSuite struct {
	suite.Suite
	echo          *echo.Echo
	vestingPeriod *mock.MockIVestingPeriod
	handler       VestingHandler
	ctrl          *gomock.Controller
}

// SetupSuite -
func (s *VestingTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.vestingPeriod = mock.NewMockIVestingPeriod(s.ctrl)
	s.handler = *NewVestingHandler(s.vestingPeriod)
}

// TearDownSuite -
func (s *VestingTestSuite) TearDownSuite() {
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteVesting_Run(t *testing.T) {
	suite.Run(t, new(VestingTestSuite))
}

func (s *VestingTestSuite) TestPeriods() {
	q := make(url.Values)
	q.Set("limit", "5")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/vesting/:id/periods")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.vestingPeriod.EXPECT().
		ByVesting(gomock.Any(), uint64(1), 5, 0).
		Return([]storage.VestingPeriod{
			{
				Time:   testTime,
				Amount: decimal.RequireFromString("1000"),
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Periods(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.VestingPeriod
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	item := response[0]
	s.Require().EqualValues("1000", item.Amount)
	s.Require().Equal(testTime, item.Time)
}
