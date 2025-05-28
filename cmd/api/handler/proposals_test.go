// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var testProposal = storage.Proposal{
	Id:          1,
	Height:      55555,
	Title:       "test proposal",
	Description: "test description",
	Deposit:     decimal.NewFromFloat(1000000),
	Status:      types.ProposalStatusActive,
	VotesCount:  123,
	Yes:         123,
	No:          0,
	Abstain:     0,
	NoWithVeto:  0,
}

// ProposalTestSuite -
type ProposalTestSuite struct {
	suite.Suite
	proposal *mock.MockIProposal
	votes    *mock.MockIVote
	address  *mock.MockIAddress
	echo     *echo.Echo
	handler  ProposalsHandler
	ctrl     *gomock.Controller
}

// SetupSuite -
func (s *ProposalTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.proposal = mock.NewMockIProposal(s.ctrl)
	s.votes = mock.NewMockIVote(s.ctrl)
	s.address = mock.NewMockIAddress(s.ctrl)
	s.handler = NewProposalsHandler(
		s.proposal,
		s.votes,
		s.address,
	)
}

// TearDownSuite -
func (s *ProposalTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(s.T().Context()))
}

func TestSuiteProposal_Run(t *testing.T) {
	suite.Run(t, new(ProposalTestSuite))
}

func (s *ProposalTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/proposal")

	s.proposal.EXPECT().
		ListWithFilters(gomock.Any(), storage.ListProposalFilters{
			Limit:      10,
			Offset:     0,
			ProposerId: 0,
			Status:     make([]types.ProposalStatus, 0),
			Type:       make([]types.ProposalType, 0),
			Sort:       sdk.SortOrderDesc,
		}).
		Return([]storage.Proposal{
			testProposal,
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var proposals []responses.Proposal
	err := json.NewDecoder(rec.Body).Decode(&proposals)
	s.Require().NoError(err)
	s.Require().Len(proposals, 1)
	s.Require().EqualValues(1, proposals[0].Id)
}

func (s *ProposalTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/proposal/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.proposal.EXPECT().
		ById(gomock.Any(), uint64(1)).
		Return(testProposal, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var proposal responses.Proposal
	err := json.NewDecoder(rec.Body).Decode(&proposal)
	s.Require().NoError(err)
	s.Require().EqualValues(1, proposal.Id)
	s.Require().EqualValues("test proposal", proposal.Title)
	s.Require().EqualValues("test description", proposal.Description)
	s.Require().EqualValues(55555, proposal.Height)
	s.Require().EqualValues(types.ProposalStatusActive, proposal.Status)
	s.Require().EqualValues(123, proposal.VotesCount)
	s.Require().EqualValues(123, proposal.Yes)
	s.Require().EqualValues(0, proposal.No)
	s.Require().EqualValues(0, proposal.Abstain)
	s.Require().EqualValues(0, proposal.NoWithVeto)
}

func (s *ProposalTestSuite) TestVotes() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/proposal/:id/votes")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.votes.EXPECT().
		ByProposalId(gomock.Any(), uint64(1), storage.VoteFilters{
			Limit:  10,
			Offset: 0,
		}).
		Return([]storage.Vote{
			{
				Id:          1,
				Height:      66666,
				Weight:      decimal.NewFromFloat(1),
				Option:      types.VoteOptionYes,
				ValidatorId: 1,
				Validator:   &testValidator,
			},
			{
				Id:          2,
				Height:      66666,
				Weight:      decimal.NewFromFloat(1),
				Option:      types.VoteOptionYes,
				ValidatorId: 0,
				Voter: &storage.Address{
					Id:         111,
					Hash:       testHashAddress,
					Address:    testAddress,
					Height:     333,
					LastHeight: 333,
					Balance: storage.Balance{
						Currency:  "utia",
						Spendable: decimal.RequireFromString("100"),
						Delegated: decimal.RequireFromString("1"),
						Unbonding: decimal.RequireFromString("2"),
					},
					Celestials: &celestials.Celestial{
						Id:       "name",
						ImageUrl: "image",
					},
				},
			},
		}, nil)

	s.Require().NoError(s.handler.Votes(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var votes []responses.Vote
	err := json.NewDecoder(rec.Body).Decode(&votes)
	s.Require().NoError(err)
	s.Require().Len(votes, 2)
	s.Require().EqualValues(1, votes[0].Id)
	s.Require().EqualValues(types.VoteOptionYes, votes[0].Option)
	s.Require().NotNil(votes[0].Validator)
	s.Require().Nil(votes[0].Voter)
	s.Require().EqualValues("moniker", votes[0].Validator.Moniker)
	s.Require().EqualValues(1, votes[0].Validator.Id)
	s.Require().EqualValues(66666, votes[1].Height)
	s.Require().EqualValues(2, votes[1].Id)
	s.Require().EqualValues(types.VoteOptionYes, votes[1].Option)
	s.Require().Nil(votes[1].Validator)
	s.Require().NotNil(votes[1].Voter)
	s.Require().NotNil(votes[1].Voter.Celestials)
	s.Require().EqualValues("image", votes[1].Voter.Celestials.ImageUrl)
}
