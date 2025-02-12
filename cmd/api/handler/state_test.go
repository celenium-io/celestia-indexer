// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// StateTestSuite -
type StateTestSuite struct {
	suite.Suite
	state      *mock.MockIState
	validators *mock.MockIValidator
	echo       *echo.Echo
	handler    *StateHandler
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *StateTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.state = mock.NewMockIState(s.ctrl)
	s.validators = mock.NewMockIValidator(s.ctrl)
	s.handler = NewStateHandler(s.state, s.validators, testIndexerName)
}

// TearDownSuite -
func (s *StateTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteState_Run(t *testing.T) {
	suite.Run(t, new(StateTestSuite))
}

func (s *StateTestSuite) TestHead() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/head")

	s.validators.EXPECT().
		TotalVotingPower(gomock.Any()).
		Return(decimal.RequireFromString("100"), nil).
		Times(1)

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(storage.State{
			Id:              1,
			Name:            testIndexerName,
			LastHeight:      100,
			LastTime:        testTime,
			TotalTx:         1234,
			TotalAccounts:   123,
			TotalFee:        decimal.RequireFromString("2"),
			TotalBlobsSize:  30,
			TotalValidators: 10,
			TotalStake:      decimal.NewFromInt(100),
			TotalNamespaces: 100,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Head(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var state responses.State
	err := json.NewDecoder(rec.Body).Decode(&state)
	s.Require().NoError(err)
	s.Require().EqualValues(1, state.Id)
	s.Require().EqualValues(testIndexerName, state.Name)
	s.Require().EqualValues(100, state.LastHeight)
	s.Require().EqualValues(1234, state.TotalTx)
	s.Require().EqualValues(123, state.TotalAccounts)
	s.Require().Equal("2", state.TotalFee)
	s.Require().EqualValues(30, state.TotalBlobsSize)
	s.Require().EqualValues(10, state.TotalValidators)
	s.Require().EqualValues(100, state.TotalNamespaces)
	s.Require().Equal(testTime, state.LastTime)
	s.Require().Equal("100", state.TotalVotingPower)
	s.Require().Equal("100", state.TotalStake)
}

func (s *StateTestSuite) TestNoHead() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/head")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(storage.State{}, sql.ErrNoRows).
		Times(1)

	s.state.EXPECT().
		IsNoRows(sql.ErrNoRows).
		Return(true).
		Times(1)

	s.Require().NoError(s.handler.Head(c))
	s.Require().Equal(http.StatusNoContent, rec.Code)
}
