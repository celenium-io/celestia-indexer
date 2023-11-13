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

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

// GasTestSuite -
type GasTestSuite struct {
	suite.Suite
	echo    *echo.Echo
	handler GasHandler
}

// SetupSuite -
func (s *GasTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.handler = NewGasHandler()
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

	s.Require().NoError(s.handler.EstimateForPfb(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response uint64
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Greater(response, uint64(0))
}
