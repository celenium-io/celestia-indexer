// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestJailByValidator() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	jails, err := s.storage.Jails.ByValidator(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(jails, 1)

	j := jails[0]
	s.Require().EqualValues(1, j.Id)
	s.Require().EqualValues(1, j.ValidatorId)
	s.Require().EqualValues("double_sign", j.Reason)
	s.Require().EqualValues("10000", j.Burned.String())
}
