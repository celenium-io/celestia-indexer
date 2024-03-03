// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestUndelegationByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	undelegations, err := s.storage.Undelegation.ByAddress(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(undelegations, 1)

	d := undelegations[0]
	s.Require().EqualValues(1, d.Id)
	s.Require().EqualValues(1, d.AddressId)
	s.Require().EqualValues(1, d.ValidatorId)
	s.Require().EqualValues("1000", d.Amount.String())
	s.Require().NotNil(d.Validator)
	s.Require().Equal("Conqueror", d.Validator.Moniker)
	s.Require().Nil(d.Address)
}
