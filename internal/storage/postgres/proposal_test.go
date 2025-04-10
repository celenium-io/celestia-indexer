// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
)

func (s *StorageTestSuite) TestProposalByProposer() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	proposals, err := s.storage.Proposals.ListWithFilters(ctx, storage.ListProposalFilters{
		Limit:      10,
		Offset:     0,
		ProposerId: 1,
	})
	s.Require().NoError(err)
	s.Require().Len(proposals, 1)

	proposal := proposals[0]
	s.Require().EqualValues(1, proposal.Id)
	s.Require().EqualValues(1, proposal.ProposerId)
	s.Require().EqualValues(1000, proposal.Height)
	s.Require().EqualValues("Description", proposal.Description)
	s.Require().NotNil(proposal.Proposer)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", proposal.Proposer.String())
}

func (s *StorageTestSuite) TestProposalByStatus() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	proposals, err := s.storage.Proposals.ListWithFilters(ctx, storage.ListProposalFilters{
		Limit:  10,
		Offset: 0,
		Status: []types.ProposalStatus{
			types.ProposalStatusInactive,
		},
	})
	s.Require().NoError(err)
	s.Require().Len(proposals, 1)

	proposal := proposals[0]
	s.Require().EqualValues(1, proposal.Id)
	s.Require().EqualValues(1, proposal.ProposerId)
	s.Require().EqualValues(1000, proposal.Height)
	s.Require().EqualValues("Description", proposal.Description)
	s.Require().NotNil(proposal.Proposer)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", proposal.Proposer.String())
}

func (s *StorageTestSuite) TestProposalByType() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	proposals, err := s.storage.Proposals.ListWithFilters(ctx, storage.ListProposalFilters{
		Limit:  10,
		Offset: 0,
		Type: []types.ProposalType{
			types.ProposalTypeText,
		},
	})
	s.Require().NoError(err)
	s.Require().Len(proposals, 1)

	proposal := proposals[0]
	s.Require().EqualValues(1, proposal.Id)
	s.Require().EqualValues(1, proposal.ProposerId)
	s.Require().EqualValues(1000, proposal.Height)
	s.Require().EqualValues("Description", proposal.Description)
	s.Require().NotNil(proposal.Proposer)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", proposal.Proposer.String())
}
