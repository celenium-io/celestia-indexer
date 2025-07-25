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
)

func (s *StorageTestSuite) TestIbcTransferList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, fltrs := range []storage.ListIbcTransferFilters{
		{
			Limit:  1,
			Offset: 0,
			Sort:   sdk.SortOrderAsc,
		}, {
			Limit:     1,
			Offset:    0,
			Sort:      sdk.SortOrderDesc,
			ChannelId: "channel-1",
		}, {
			Limit:    1,
			Offset:   0,
			Sort:     sdk.SortOrderDesc,
			SenderId: testsuite.Ptr(uint64(1)),
		}, {
			Limit:     1,
			Offset:    0,
			Sort:      sdk.SortOrderDesc,
			AddressId: testsuite.Ptr(uint64(1)),
		}, {
			Limit:         1,
			Offset:        0,
			Sort:          sdk.SortOrderDesc,
			ConnectionIds: []string{"connection-1"},
		},
	} {

		transfers, err := s.storage.IbcTransfers.List(ctx, fltrs)
		s.Require().NoError(err)
		s.Require().Len(transfers, 1)

		transfer := transfers[0]
		s.Require().EqualValues("connection-1", transfer.ConnectionId)
		s.Require().EqualValues("channel-1", transfer.ChannelId)
		s.Require().EqualValues(1000, transfer.Height)
		s.Require().EqualValues("123456", transfer.Amount.String())
		s.Require().EqualValues("utia", transfer.Denom)
		s.Require().Zero(transfer.HeightTimeout)
		s.Require().Nil(transfer.Timeout)
		s.Require().EqualValues(321654, transfer.Sequence)
		s.Require().NotNil(transfer.Tx)
		s.Require().NotNil(transfer.Sender)

		txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
		s.Require().NoError(err)
		s.Require().Equal(txHash, transfer.Tx.Hash)

		s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", transfer.Sender.Address)

		s.Require().NotNil(transfer.Connection)
		s.Require().NotNil(transfer.Connection.Client)
		s.Require().Equal("osmosis-1", transfer.Connection.Client.ChainId)
	}
}

func (s *StorageTestSuite) TestIbcTransferSeries() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	type args struct {
		tf     storage.Timeframe
		column string
		req    storage.SeriesRequest

		wantCount int
	}

	for _, fltrs := range []args{
		{
			tf:        storage.TimeframeHour,
			column:    "count",
			req:       storage.NewSeriesRequest(0, 0),
			wantCount: 1,
		}, {
			tf:        storage.TimeframeDay,
			column:    "count",
			req:       storage.NewSeriesRequest(0, 0),
			wantCount: 1,
		}, {
			tf:        storage.TimeframeMonth,
			column:    "count",
			req:       storage.NewSeriesRequest(0, 0),
			wantCount: 1,
		}, {
			tf:        storage.TimeframeHour,
			column:    "amount",
			req:       storage.NewSeriesRequest(0, 0),
			wantCount: 1,
		}, {
			tf:        storage.TimeframeDay,
			column:    "amount",
			req:       storage.NewSeriesRequest(0, 0),
			wantCount: 1,
		}, {
			tf:        storage.TimeframeMonth,
			column:    "amount",
			req:       storage.NewSeriesRequest(0, 0),
			wantCount: 1,
		}, {
			tf:        storage.TimeframeHour,
			column:    "amount",
			req:       storage.NewSeriesRequest(1715942016, 0),
			wantCount: 0,
		}, {
			tf:        storage.TimeframeHour,
			column:    "amount",
			req:       storage.NewSeriesRequest(1652783616, 0),
			wantCount: 1,
		}, {
			tf:        storage.TimeframeHour,
			column:    "amount",
			req:       storage.NewSeriesRequest(0, 1715942016),
			wantCount: 1,
		}, {
			tf:        storage.TimeframeHour,
			column:    "amount",
			req:       storage.NewSeriesRequest(0, 1652783616),
			wantCount: 0,
		},
	} {
		series, err := s.storage.IbcTransfers.Series(ctx, "channel-1", fltrs.tf, fltrs.column, fltrs.req)
		s.Require().NoError(err)
		s.Require().Len(series, fltrs.wantCount)
	}
}

func (s *StorageTestSuite) TestLargestTransfer24h() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	transfer, err := s.storage.IbcTransfers.LargestTransfer24h(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(transfer)
	s.Require().EqualValues("connection-2", transfer.ConnectionId)
	s.Require().EqualValues("channel-2", transfer.ChannelId)
	s.Require().EqualValues(1002, transfer.Height)
	s.Require().EqualValues("12345678", transfer.Amount.String())
	s.Require().EqualValues("utia", transfer.Denom)
	s.Require().Zero(transfer.HeightTimeout)
	s.Require().Nil(transfer.Timeout)
	s.Require().EqualValues(321656, transfer.Sequence)
	s.Require().NotNil(transfer.Tx)
	s.Require().NotNil(transfer.Sender)
}
