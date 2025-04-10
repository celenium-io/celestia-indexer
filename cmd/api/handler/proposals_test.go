// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var testProposal = storage.Proposal{
	Id: 1,
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
