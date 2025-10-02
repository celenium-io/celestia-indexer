// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func (s *StorageTestSuite) TestValidatorByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	validator, err := s.storage.Validator.ByAddress(ctx, "celestiavaloper17vmk8m246t648hpmde2q7kp4ft9uwrayy09dmw")
	s.Require().NoError(err)

	s.Require().Equal("celestiavaloper17vmk8m246t648hpmde2q7kp4ft9uwrayy09dmw", validator.Address)
	s.Require().Equal("celestia17vmk8m246t648hpmde2q7kp4ft9uwrayps85dg", validator.Delegator)
	s.Require().Equal("Conqueror", validator.Moniker)
	s.Require().Equal("https://github.com/DasRasyo", validator.Website)
	s.Require().Equal("EAD22B173DE57E6A", validator.Identity)
	s.Require().Equal("https://t.me/DasRasyo || conqueror.prime", validator.Contacts)
	s.Require().Equal("1", validator.MinSelfDelegation.String())
	s.Require().Equal("0.2", validator.MaxRate.String())
}

func (s *StorageTestSuite) TestTotalPower() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	power, err := s.storage.Validator.TotalVotingPower(ctx)
	s.Require().NoError(err)
	s.Require().Equal("2", power.String())
}

func (s *StorageTestSuite) TestListByPower() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	validators, err := s.storage.Validator.ListByPower(ctx, storage.ValidatorFilters{
		Limit: 10,
	})
	s.Require().NoError(err)
	s.Require().Len(validators, 2)
}

func (s *StorageTestSuite) TestJailedCount() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	count, err := s.storage.Validator.JailedCount(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(0, count)
}

func (s *StorageTestSuite) TestMessages() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	msgs, err := s.storage.Validator.Messages(ctx, 1, storage.ValidatorMessagesFilters{})
	s.Require().NoError(err)
	s.Require().Len(msgs, 1)
	s.Require().NotNil(msgs[0].Msg)
}

func (s *StorageTestSuite) TestMetrics() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW validator_metrics;")
	s.Require().NoError(err)

	metrics, err := s.storage.Validator.Metrics(ctx, 1)
	s.Require().NoError(err)
	s.Require().NotEmpty(metrics.Id)
	s.Require().NotEmpty(metrics.BlockMissedMetric.String())
	s.Require().NotEmpty(metrics.VotesMetric.String())
	s.Require().NotEmpty(metrics.OperationTimeMetric.String())
	s.Require().NotEmpty(metrics.CommissionMetric.String())
	s.Require().NotEmpty(metrics.SelfDelegationMetric.String())
}

func (s *StorageTestSuite) TestTopNMetrics() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW validator_metrics;")
	s.Require().NoError(err)

	metrics, err := s.storage.Validator.TopNMetrics(ctx, 10)
	s.Require().NoError(err)
	s.Require().NotEmpty(metrics.BlockMissedMetric.String())
	s.Require().NotEmpty(metrics.VotesMetric.String())
	s.Require().NotEmpty(metrics.OperationTimeMetric.String())
	s.Require().NotEmpty(metrics.CommissionMetric.String())
	s.Require().NotEmpty(metrics.SelfDelegationMetric.String())
}
