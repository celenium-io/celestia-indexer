// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestHyperlaneIgpByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	igp, err := s.storage.HLIGP.ByHash(ctx, []byte("igp_1"))
	s.Require().NoError(err)

	s.Require().EqualValues(1, igp.Id)
	s.Require().EqualValues(1488, igp.Height)
	s.Require().EqualValues([]byte("igp_1"), igp.IgpId)

	s.Require().Len(igp.Configs, 2)
	config := igp.Configs[0]
	s.Require().EqualValues(1234, config.RemoteDomain)
	s.Require().EqualValues("4321", config.TokenExchangeRate)
	s.Require().EqualValues("100000", config.GasOverhead.String())
	s.Require().EqualValues("1", config.GasPrice.String())
}

func (s *StorageTestSuite) TestHyperlaneIgpList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.HLIGP.List(ctx, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(items, 2)

	igp := items[1]
	s.Require().EqualValues(2, igp.Id)
	s.Require().EqualValues(1489, igp.Height)
	s.Require().EqualValues([]byte("igp_2"), igp.IgpId)

	s.Require().Len(igp.Configs, 0)
}
