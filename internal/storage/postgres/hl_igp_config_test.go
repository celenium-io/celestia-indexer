// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestHyperlaneIgpConfigList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.HLIGPConfig.List(ctx, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	config := items[0]
	s.Require().EqualValues(1, config.Id)
	s.Require().EqualValues(1488, config.Height)
	s.Require().EqualValues("100000", config.GasOverhead.String())
	s.Require().EqualValues("1", config.GasPrice.String())
	s.Require().EqualValues(1234, config.RemoteDomain)
	s.Require().EqualValues("4321", config.TokenExchangeRate)
}
