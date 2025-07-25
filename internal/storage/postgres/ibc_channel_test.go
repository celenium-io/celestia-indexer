// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestIbcChannelById() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	channel, err := s.storage.IbcChannels.ById(ctx, "channel-1")
	s.Require().NoError(err)

	s.Require().EqualValues("connection-1", channel.ConnectionId)
	s.Require().EqualValues("client-1", channel.ClientId)
	s.Require().EqualValues("channel-10", channel.CounterpartyChannelId)
	s.Require().EqualValues("transfer", channel.CounterpartyPortId)
	s.Require().EqualValues("transfer", channel.PortId)
	s.Require().EqualValues("ics20-1", channel.Version)
	s.Require().EqualValues(1000, channel.Height)
	s.Require().EqualValues(1, channel.CreateTxId)
	s.Require().EqualValues(2, channel.ConfirmationTxId)
	s.Require().EqualValues(1001, channel.ConfirmationHeight)
	s.Require().Equal(types.IbcChannelStatusOpened, channel.Status)
	s.Require().NotNil(channel.ConfirmationTx)
	s.Require().NotNil(channel.CreateTx)
	s.Require().NotNil(channel.Client)
	s.Require().NotNil(channel.Creator)

	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, channel.CreateTx.Hash)

	txHash2, err := hex.DecodeString("652452A670011D629CC116E510BA88C1CABE061336661B1F3D206D248BD55811")
	s.Require().NoError(err)
	s.Require().Equal(txHash2, channel.ConfirmationTx.Hash)

	s.Require().EqualValues("osmosis-1", channel.Client.ChainId)
	s.Require().EqualValues("client", channel.Client.Type)
	s.Require().EqualValues(1, channel.Client.ConnectionCount)

	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", channel.Creator.Address)
}

func (s *StorageTestSuite) TestIbcChannelList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, fltrs := range []storage.ListChannelFilters{
		{
			Limit:  10,
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
		}, {
			Limit:    1,
			Offset:   0,
			Sort:     sdk.SortOrderDesc,
			ClientId: "client-1",
		}, {
			Limit:        1,
			Offset:       0,
			Sort:         sdk.SortOrderDesc,
			ConnectionId: "connection-1",
		}, {
			Limit:        1,
			Offset:       0,
			Sort:         sdk.SortOrderDesc,
			ClientId:     "client-1",
			ConnectionId: "connection-1",
		}, {
			Limit:  1,
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
			Status: types.IbcChannelStatusOpened,
		},
	} {

		channels, err := s.storage.IbcChannels.List(ctx, fltrs)
		s.Require().NoError(err)
		s.Require().Len(channels, 1)

		channel := channels[0]
		s.Require().EqualValues("connection-1", channel.ConnectionId)
		s.Require().EqualValues("client-1", channel.ClientId)
		s.Require().EqualValues("channel-10", channel.CounterpartyChannelId)
		s.Require().EqualValues("transfer", channel.CounterpartyPortId)
		s.Require().EqualValues("transfer", channel.PortId)
		s.Require().EqualValues("ics20-1", channel.Version)
		s.Require().EqualValues(1000, channel.Height)
		s.Require().EqualValues(1, channel.CreateTxId)
		s.Require().EqualValues(2, channel.ConfirmationTxId)
		s.Require().EqualValues(1001, channel.ConfirmationHeight)
		s.Require().Equal(types.IbcChannelStatusOpened, channel.Status)
		s.Require().NotNil(channel.ConfirmationTx)
		s.Require().NotNil(channel.CreateTx)
		s.Require().NotNil(channel.Client)
		s.Require().NotNil(channel.Creator)

		txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
		s.Require().NoError(err)
		s.Require().Equal(txHash, channel.CreateTx.Hash)

		txHash2, err := hex.DecodeString("652452A670011D629CC116E510BA88C1CABE061336661B1F3D206D248BD55811")
		s.Require().NoError(err)
		s.Require().Equal(txHash2, channel.ConfirmationTx.Hash)

		s.Require().EqualValues("osmosis-1", channel.Client.ChainId)
		s.Require().EqualValues("client", channel.Client.Type)
		s.Require().EqualValues(1, channel.Client.ConnectionCount)

		s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", channel.Creator.Address)
	}
}

func (s *StorageTestSuite) TestIbcChannelStatsByChainId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	stats, err := s.storage.IbcChannels.StatsByChain(ctx, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(stats, 1)

	s.Require().Equal("osmosis-1", stats[0].Chain)
	s.Require().Equal("100", stats[0].Sent.String())
	s.Require().Equal("100", stats[0].Received.String())
}

func (s *StorageTestSuite) TestIbcBusiestChannel1m() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	channel, err := s.storage.IbcChannels.BusiestChannel1m(ctx)
	s.Require().NoError(err)

	s.Require().Equal("channel-2", channel.ChannelId)
	s.Require().Equal(int64(2), channel.TransfersCount)
}
