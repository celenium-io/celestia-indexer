// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var testSignal = storage.SignalVersion{
	Id:          1,
	Height:      12345,
	ValidatorId: 1,
	Time:        time.Now().UTC(),
	VotingPower: decimal.NewFromFloat(100),
	Version:     1,
	MsgId:       1,
	TxId:        1,
}

var testUpgrade = storage.Upgrade{
	Id:       2,
	Height:   101,
	SignerId: 2,
	Time:     time.Now().UTC(),
	Version:  1,
	MsgId:    1,
	TxId:     1,
}

// SignalTestSuite -
type SignalTestSuite struct {
	suite.Suite
	signals    *mock.MockISignalVersion
	upgrades   *mock.MockIUpgrade
	validators *mock.MockIValidator
	txs        *mock.MockITx
	address    *mock.MockIAddress
	echo       *echo.Echo
	handler    *SignalHandler
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *SignalTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.signals = mock.NewMockISignalVersion(s.ctrl)
	s.upgrades = mock.NewMockIUpgrade(s.ctrl)
	s.validators = mock.NewMockIValidator(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.address = mock.NewMockIAddress(s.ctrl)
	s.handler = NewSignalHandler(s.signals, s.upgrades, s.validators, s.txs, s.address)
}

// TearDownSuite -
func (s *SignalTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(s.T().Context()))
}

func TestSuiteSignal_Run(t *testing.T) {
	suite.Run(t, new(SignalTestSuite))
}

func (s *SignalTestSuite) TestList() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("validator_id", "1")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/signal")

	s.signals.EXPECT().
		List(gomock.Any(), storage.ListSignalsFilter{
			Limit:       10,
			Offset:      0,
			ValidatorId: 1,
			Sort:        sdk.SortOrderDesc,
		}).
		Return([]storage.SignalVersion{
			testSignal,
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var signals []responses.SignalVersion
	err := json.NewDecoder(rec.Body).Decode(&signals)
	s.Require().NoError(err)
	s.Require().Len(signals, 1)
	s.Require().EqualValues(1, signals[0].Id)
}

func (s *SignalTestSuite) TestUpgrades() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("height", "101")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/signal/upgrade")

	s.upgrades.EXPECT().
		List(gomock.Any(), storage.ListUpgradesFilter{
			Limit:  10,
			Offset: 0,
			Height: 101,
			Sort:   sdk.SortOrderDesc,
		}).
		Return([]storage.Upgrade{
			testUpgrade,
		}, nil)

	s.Require().NoError(s.handler.Upgrades(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var upgrades []responses.Upgrade
	err := json.NewDecoder(rec.Body).Decode(&upgrades)
	s.Require().NoError(err)
	s.Require().Len(upgrades, 1)
	s.Require().EqualValues(2, upgrades[0].Id)
}
