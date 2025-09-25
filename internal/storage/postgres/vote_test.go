// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
)

func (s *StorageTestSuite) TestVoteByProposalId() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	votes, err := s.storage.Votes.ByProposalId(ctx, 1, storage.VoteFilters{
		Limit:  10,
		Offset: 0,
	})
	s.Require().NoError(err)
	s.Require().Len(votes, 2)

	vote := votes[0]
	s.Require().EqualValues(1, vote.Id)
	s.Require().EqualValues(1000, vote.Height)
	s.Require().EqualValues(types.VoteOptionYes, vote.Option)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", vote.Voter.Address)
}

func (s *StorageTestSuite) TestVoteByProposalIdWithOption() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	votes, err := s.storage.Votes.ByProposalId(ctx, 1, storage.VoteFilters{
		Limit:  10,
		Offset: 0,
		Option: []types.VoteOption{types.VoteOptionNo},
	})
	s.Require().NoError(err)
	s.Require().Len(votes, 1)

	vote := votes[0]
	s.Require().EqualValues(2, vote.Id)
	s.Require().EqualValues(1000, vote.Height)
	s.Require().EqualValues(types.VoteOptionNo, vote.Option)
	s.Require().EqualValues("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", vote.Voter.Address)
}

func (s *StorageTestSuite) TestVoteByProposalIdWithVoterType() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	votes, err := s.storage.Votes.ByProposalId(ctx, 2, storage.VoteFilters{
		Limit:     10,
		Offset:    0,
		VoterType: types.VoterTypeValidator,
	})
	s.Require().NoError(err)
	s.Require().Len(votes, 1)

	vote := votes[0]
	s.Require().EqualValues(3, vote.Id)
	s.Require().EqualValues(1000, vote.Height)
	s.Require().EqualValues(types.VoteOptionAbstain, vote.Option)
	s.Require().EqualValues("81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8", vote.Validator.ConsAddress)
	s.Require().EqualValues("Conqueror", vote.Validator.Moniker)
}

func (s *StorageTestSuite) TestVoteByProposalIdWithVoter() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	votes, err := s.storage.Votes.ByProposalId(ctx, 1, storage.VoteFilters{
		Limit:     10,
		Offset:    0,
		AddressId: testsuite.Ptr(uint64(2)),
	})
	s.Require().NoError(err)
	s.Require().Len(votes, 1)

	vote := votes[0]
	s.Require().EqualValues(2, vote.Id)
	s.Require().EqualValues(1000, vote.Height)
	s.Require().EqualValues(2, vote.VoterId)
	s.Require().NotNil(vote.Voter)
	s.Require().EqualValues("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", vote.Voter.Address)
}

func (s *StorageTestSuite) TestVoteByProposalIdWithValidator() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	votes, err := s.storage.Votes.ByProposalId(ctx, 2, storage.VoteFilters{
		Limit:       10,
		Offset:      0,
		ValidatorId: testsuite.Ptr(uint64(1)),
	})
	s.Require().NoError(err)
	s.Require().Len(votes, 1)

	vote := votes[0]
	s.Require().EqualValues(3, vote.Id)
	s.Require().EqualValues(1000, vote.Height)
	s.Require().NotNil(vote.ValidatorId)
	s.Require().EqualValues(1, *vote.ValidatorId)
	s.Require().NotNil(vote.Validator)
	s.Require().EqualValues("Conqueror", vote.Validator.Moniker)
	s.Require().EqualValues("81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8", vote.Validator.ConsAddress)
}

func (s *StorageTestSuite) TestVoteByVoterId() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	votes, err := s.storage.Votes.ByVoterId(ctx, 2, storage.VoteFilters{
		Limit:  10,
		Offset: 0,
	})
	s.Require().NoError(err)
	s.Require().Len(votes, 2)

	vote := votes[0]
	s.Require().EqualValues(2, vote.Id)
	s.Require().EqualValues(1000, vote.Height)
	s.Require().EqualValues(types.VoteOptionNo, vote.Option)
	s.Require().Nil(vote.Validator)
	s.Require().NotNil(vote.Proposal)
	s.Require().EqualValues(1, vote.Proposal.Id)
	s.Require().EqualValues("Title", vote.Proposal.Title)
	s.Require().EqualValues("Description", vote.Proposal.Description)
	s.Require().EqualValues(types.ProposalStatusInactive, vote.Proposal.Status)
}

func (s *StorageTestSuite) TestVoteByValidatorId() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	votes, err := s.storage.Votes.ByValidatorId(ctx, 1, storage.VoteFilters{
		Limit:  10,
		Offset: 0,
	})
	s.Require().NoError(err)
	s.Require().Len(votes, 1)

	vote := votes[0]
	s.Require().EqualValues(3, vote.Id)
	s.Require().EqualValues(1000, vote.Height)
	s.Require().EqualValues(types.VoteOptionAbstain, vote.Option)
}

func (s *StorageTestSuite) TestVoteByValidatorIdNoResult() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	votes, err := s.storage.Votes.ByValidatorId(ctx, 4, storage.VoteFilters{
		Limit:  10,
		Offset: 0,
	})
	s.Require().NoError(err)
	s.Require().Len(votes, 0)
}
