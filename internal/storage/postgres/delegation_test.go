// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestDelegationByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	delegations, err := s.storage.Delegation.ByAddress(ctx, 1, 10, 0, false)
	s.Require().NoError(err)
	s.Require().Len(delegations, 1)

	d := delegations[0]
	s.Require().EqualValues(1, d.Id)
	s.Require().EqualValues(1, d.AddressId)
	s.Require().EqualValues(1, d.ValidatorId)
	s.Require().EqualValues("10000", d.Amount.String())
	s.Require().NotNil(d.Validator)
	s.Require().Equal("Conqueror", d.Validator.Moniker)
	s.Require().Nil(d.Address)
}

func (s *StorageTestSuite) TestDelegationByValidator() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	delegations, err := s.storage.Delegation.ByValidator(ctx, 1, 10, 0, true)
	s.Require().NoError(err)
	s.Require().Len(delegations, 2)

	d := delegations[0]
	s.Require().EqualValues(1, d.Id)
	s.Require().EqualValues(1, d.AddressId)
	s.Require().EqualValues(1, d.ValidatorId)
	s.Require().EqualValues("10000", d.Amount.String())
	s.Require().Nil(d.Validator)
	s.Require().NotNil(d.Address)
	s.Require().NotNil(d.Address.Celestials)
}
