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

func (s *StorageTestSuite) TestJailedCOunt() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	count, err := s.storage.Validator.JailedCount(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(0, count)
}
