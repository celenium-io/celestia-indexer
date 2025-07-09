// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestRedelegationByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	redelegations, err := s.storage.Redelegation.ByAddress(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(redelegations, 1)

	d := redelegations[0]
	s.Require().EqualValues(1, d.Id)
	s.Require().EqualValues(1, d.AddressId)
	s.Require().EqualValues(1, d.SrcId)
	s.Require().EqualValues(2, d.DestId)
	s.Require().EqualValues("1000", d.Amount.String())
	s.Require().NotNil(d.Source)
	s.Require().Equal("Conqueror", d.Source.Moniker)
	s.Require().NotNil(d.Destination)
	s.Require().Equal("Witval", d.Destination.Moniker)
	s.Require().Nil(d.Address)
}
