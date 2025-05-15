// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
			Limit:  10,
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
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
	}
}
