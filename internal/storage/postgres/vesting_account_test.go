// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestVestingAccountByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	vestings, err := s.storage.VestingAccounts.ByAddress(ctx, 1, 1, 0, true)
	s.Require().NoError(err)
	s.Require().Len(vestings, 1)

	vesting := vestings[0]
	s.Require().Equal("100000", vesting.Amount.String())
	s.Require().Equal("delayed", vesting.Type.String())
}
