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
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// ValidatorTestSuite -
type ValidatorTestSuite struct {
	suite.Suite
	validators      *mock.MockIValidator
	blocks          *mock.MockIBlock
	blockSignatures *mock.MockIBlockSignature
	delegations     *mock.MockIDelegation
	jails           *mock.MockIJail
	constants       *mock.MockIConstant
	votes           *mock.MockIVote
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
	s.delegations = mock.NewMockIDelegation(s.ctrl)
	s.constants = mock.NewMockIConstant(s.ctrl)
	s.jails = mock.NewMockIJail(s.ctrl)
	s.votes = mock.NewMockIVote(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.handler = NewValidatorHandler(s.validators, s.blocks, s.blockSignatures, s.delegations, s.constants, s.jails, s.votes, s.state, testIndexerName)
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
		ListByPower(gomock.Any(), storage.ValidatorFilters{
			Limit: 10,
		}).
		Return([]storage.Validator{
			testValidator,
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

func (s *ValidatorTestSuite) TestDelegators() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("show_zero", "true")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validators/:id/delegators")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.delegations.EXPECT().
		ByValidator(gomock.Any(), uint64(1), 10, 0, true).
		Return([]storage.Delegation{
			{
				AddressId:   1,
				ValidatorId: 1,
				Amount:      decimal.RequireFromString("100"),
				Validator:   &testValidator,
				Address: &storage.Address{
					Address: testAddress,
					Id:      1,
				},
			},
		}, nil)

	s.Require().NoError(s.handler.Delegators(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var delegations []responses.Delegation
	err := json.NewDecoder(rec.Body).Decode(&delegations)
	s.Require().NoError(err)
	s.Require().Len(delegations, 1)

	d := delegations[0]
	s.Require().Equal("100", d.Amount)
	s.Require().Equal(testAddress, d.Delegator.Hash)
}

func (s *ValidatorTestSuite) TestJails() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validators/:id/jails")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.jails.EXPECT().
		ByValidator(gomock.Any(), uint64(1), 10, 0).
		Return([]storage.Jail{
			{
				Burned:      decimal.RequireFromString("100"),
				Reason:      "double_sign",
				Height:      100,
				Time:        testTime,
				Id:          1,
				ValidatorId: 1,
			},
		}, nil)

	s.Require().NoError(s.handler.Jails(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var jail []responses.Jail
	err := json.NewDecoder(rec.Body).Decode(&jail)
	s.Require().NoError(err)
	s.Require().Len(jail, 1)

	j := jail[0]
	s.Require().Equal("100", j.Burned)
	s.Require().Equal("double_sign", j.Reason)
}

func (s *ValidatorTestSuite) TestCount() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validators/count")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(storage.State{
			LastHeight:      4,
			TotalValidators: 10,
		}, nil).
		Times(1)

	s.validators.EXPECT().
		JailedCount(gomock.Any()).
		Return(2, nil).
		Times(1)

	s.constants.EXPECT().
		Get(gomock.Any(), storageTypes.ModuleNameStaking, "max_validators").
		Return(storage.Constant{
			Value: "6",
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Count(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count responses.ValidatorCount
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)

	s.Require().EqualValues(10, count.Total)
	s.Require().EqualValues(2, count.Jailed)
	s.Require().EqualValues(6, count.Active)
	s.Require().EqualValues(2, count.Inactive)
}

func (s *ValidatorTestSuite) TestVotes() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validators/:id/votes")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.votes.EXPECT().
		ByValidatorId(gomock.Any(), uint64(1), storage.VoteFilters{
			Limit:  10,
			Offset: 0,
		}).
		Return([]storage.Vote{
			{
				Id:          2,
				Height:      1000,
				Weight:      decimal.NewFromFloat(1),
				Option:      storageTypes.VoteOptionNoWithVeto,
				ValidatorId: testsuite.Ptr(uint64(1)),
				Validator:   &testValidator,
			},
		}, nil)

	s.Require().NoError(s.handler.Votes(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var votes []responses.Vote
	err := json.NewDecoder(rec.Body).Decode(&votes)
	s.Require().NoError(err)
	s.Require().Len(votes, 1)
	s.Require().EqualValues(2, votes[0].Id)
	s.Require().EqualValues(1000, votes[0].Height)
	s.Require().EqualValues(decimal.NewFromFloat(1), votes[0].Weight)
	s.Require().EqualValues(storageTypes.VoteOptionNoWithVeto, votes[0].Option)
	s.Require().EqualValues(1, votes[0].Validator.Id)
	s.Require().EqualValues("moniker", votes[0].Validator.Moniker)
	s.Require().EqualValues("012345", votes[0].Validator.ConsAddress)
}

func (s *ValidatorTestSuite) TestMessages() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validators/:id/messages")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.validators.EXPECT().
		Messages(gomock.Any(), uint64(1), storage.ValidatorMessagesFilters{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("asc"),
		}).
		Return([]storage.MsgValidator{
			{
				MsgId:       2,
				ValidatorId: 1,
				Height:      1000,
				Time:        testTime,
				Msg: &storage.Message{
					Id:       2,
					Type:     storageTypes.MsgCreateValidator,
					Size:     12,
					Position: 2,
				},
				Validator: &testValidator,
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Messages(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var messages []responses.Message
	err := json.NewDecoder(rec.Body).Decode(&messages)
	s.Require().NoError(err)
	s.Require().Len(messages, 1)

	s.Require().EqualValues(2, messages[0].Id)
	s.Require().EqualValues(1000, messages[0].Height)
	s.Require().Equal(testTime.String(), messages[0].Time.String())
	s.Require().EqualValues(12, messages[0].Size)
	s.Require().EqualValues(2, messages[0].Position)
}
