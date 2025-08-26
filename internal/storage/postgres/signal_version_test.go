// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"github.com/shopspring/decimal"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestSignalVersionList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, fltrs := range []storage.ListSignalsFilter{
		{
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
			TxId:   testsuite.Ptr(uint64(1)),
		}, {
			Offset:  0,
			Sort:    sdk.SortOrderDesc,
			Version: 1488,
		}, {
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
			From:   time.Date(2025, 8, 7, 0, 0, 0, 0, time.UTC),
			To:     time.Date(2025, 8, 8, 0, 0, 0, 0, time.UTC),
		}, {
			Limit:  1,
			Offset: 0,
			Sort:   sdk.SortOrderAsc,
		},
	} {

		signals, err := s.storage.SignalVersion.List(ctx, fltrs)
		s.Require().NoError(err)
		s.Require().Len(signals, 1)

		signal := signals[0]
		s.Require().EqualValues(1, signal.Id)
		s.Require().EqualValues(101, signal.Height)
		s.Require().EqualValues(1488, signal.Version)
		s.Require().EqualValues(decimal.RequireFromString("123"), signal.VotingPower)
		s.Require().EqualValues(1, signal.MsgId)
		s.Require().EqualValues(1, signal.TxId)
		s.Require().NotNil(signal.Validator)
		s.Require().EqualValues("celestiavaloper17vmk8m246t648hpmde2q7kp4ft9uwrayy09dmw", signal.Validator.Address)
	}
}
