// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

func (s *StorageTestSuite) TestSignalVersionList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, fltrs := range []storage.ListSignalsFilter{
		{
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
			TxId:   testsuite.Ptr(uint64(3)),
		}, {
			Offset:  0,
			Sort:    sdk.SortOrderDesc,
			Version: 1488,
		}, {
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
			From:   time.Date(2025, 8, 9, 3, 11, 0, 0, time.UTC),
			To:     time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		}, {
			Limit:  1,
			Offset: 1,
			Sort:   sdk.SortOrderAsc,
		}, {
			Sort:        sdk.SortOrderDesc,
			ValidatorId: 1,
		},
	} {

		signals, err := s.storage.SignalVersion.List(ctx, fltrs)
		s.Require().NoError(err)
		s.Require().Len(signals, 1)

		signal := signals[0]
		s.Require().EqualValues(3, signal.Id)
		s.Require().EqualValues(103, signal.Height)
		s.Require().EqualValues(1488, signal.Version)
		s.Require().EqualValues(decimal.RequireFromString("8"), signal.VotingPower)
		s.Require().EqualValues(3, signal.MsgId)
		s.Require().EqualValues(3, signal.TxId)
		s.Require().NotNil(signal.Validator)
		s.Require().EqualValues("81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8", signal.Validator.ConsAddress)
		s.Require().EqualValues(1, signal.Validator.Id)
		s.Require().EqualValues("Conqueror", signal.Validator.Moniker)

		txHash, err := hex.DecodeString("BA37478C3E9A804697271ACC474D484E9160899C86E551D737EEA819FCC75003")
		s.Require().NoError(err)
		s.Require().NotNil(signal.Tx)
		s.Require().EqualValues(txHash, signal.Tx.Hash)
	}
}
