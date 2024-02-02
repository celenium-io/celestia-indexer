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
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// ValidatorTestSuite -
type ValidatorTestSuite struct {
	suite.Suite
	validators      *mock.MockIValidator
	blocks          *mock.MockIBlock
	blockSignatures *mock.MockIBlockSignature
	state           *mock.MockIState
	echo            *echo.Echo
	handler         *ValidatorHandler
	ctrl            *gomock.Controller
}

// SetupSuite -
func (s *ValidatorTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.validators = mock.NewMockIValidator(s.ctrl)
	s.blocks = mock.NewMockIBlock(s.ctrl)
	s.blockSignatures = mock.NewMockIBlockSignature(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.handler = NewValidatorHandler(s.validators, s.blocks, s.blockSignatures, s.state, testIndexerName)
}

// TearDownSuite -
func (s *ValidatorTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteValidator_Run(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}

func (s *ValidatorTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validators/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.validators.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&testValidator, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var validator responses.Validator
	err := json.NewDecoder(rec.Body).Decode(&validator)
	s.Require().NoError(err)
	s.Require().EqualValues(1, validator.Id)
	s.Require().EqualValues("moniker", validator.Moniker)
	s.Require().EqualValues("012345", validator.ConsAddress)
}

func (s *ValidatorTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validator")

	s.validators.EXPECT().
		List(gomock.Any(), uint64(10), uint64(0), pgSort("asc")).
		Return([]*storage.Validator{
			&testValidator,
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var validators []responses.Validator
	err := json.NewDecoder(rec.Body).Decode(&validators)
	s.Require().NoError(err)
	s.Require().Len(validators, 1)
	s.Require().EqualValues(1, validators[0].Id)
	s.Require().EqualValues("moniker", validators[0].Moniker)
	s.Require().EqualValues("012345", validators[0].ConsAddress)
}

func (s *ValidatorTestSuite) TestByProposer() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validator/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.blocks.EXPECT().
		ByProposer(gomock.Any(), uint64(1), 10, 0).
		Return([]storage.Block{
			testBlock,
		}, nil)

	s.Require().NoError(s.handler.Blocks(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blocks []responses.Block
	err := json.NewDecoder(rec.Body).Decode(&blocks)
	s.Require().NoError(err)
	s.Require().Len(blocks, 1)

	block := blocks[0]
	s.Require().EqualValues(1, block.Id)
	s.Require().EqualValues(100, block.Height)
	s.Require().Equal("1", block.VersionApp)
	s.Require().Equal("11", block.VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", block.Hash.String())
	s.Require().Equal(testTime, block.Time)
	s.Require().NotNil(block.Stats)
}

func (s *ValidatorTestSuite) TestUptime() {
	q := make(url.Values)
	q.Add("limit", "4")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validators/:id/uptime")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(storage.State{
			LastHeight: 1000,
		}, nil)

	s.blockSignatures.EXPECT().
		LevelsByValidator(gomock.Any(), uint64(1), types.Level(995)).
		Return([]types.Level{999, 998, 997, 996}, nil)

	s.Require().NoError(s.handler.Uptime(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var uptime responses.ValidatorUptime
	err := json.NewDecoder(rec.Body).Decode(&uptime)
	s.Require().NoError(err)
	s.Require().EqualValues("1.0000", uptime.Uptime)
	s.Require().Len(uptime.Blocks, 4)

	block := uptime.Blocks[0]
	s.Require().True(block.Signed)
	s.Require().EqualValues(999, block.Height)
}

func (s *ValidatorTestSuite) TestUptimeUnusual() {
	q := make(url.Values)
	q.Add("limit", "10")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validators/:id/uptime")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(storage.State{
			LastHeight: 4,
		}, nil)

	s.blockSignatures.EXPECT().
		LevelsByValidator(gomock.Any(), uint64(1), types.Level(-7)).
		Return([]types.Level{2, 1}, nil)

	s.Require().NoError(s.handler.Uptime(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var uptime responses.ValidatorUptime
	err := json.NewDecoder(rec.Body).Decode(&uptime)
	s.Require().NoError(err)
	s.Require().EqualValues("0.6667", uptime.Uptime)
	s.Require().Len(uptime.Blocks, 3)

	block := uptime.Blocks[0]
	s.Require().False(block.Signed)
	s.Require().EqualValues(3, block.Height)
}
