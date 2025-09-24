// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestHyperlaneGasPaymentList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.HLGasPayment.List(ctx, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	payment := items[0]
	s.Require().EqualValues(1, payment.Id)
	s.Require().EqualValues(1488, payment.Height)
	s.Require().EqualValues("111", payment.Amount.String())
	s.Require().EqualValues("11", payment.GasAmount.String())
	s.Require().EqualValues(1, payment.IgpId)
	s.Require().EqualValues(1, payment.TransferId)
}
