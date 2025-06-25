// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestIbcConnectionById() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	conn, err := s.storage.IbcConnections.ById(ctx, "connection-1")
	s.Require().NoError(err)

	s.Require().EqualValues("connection-1", conn.ConnectionId)
	s.Require().EqualValues("client-1", conn.ClientId)
	s.Require().EqualValues("counterparty-client-1", conn.CounterpartyClientId)
	s.Require().EqualValues("counterparty-1", conn.CounterpartyConnectionId)
	s.Require().EqualValues(1000, conn.Height)
	s.Require().EqualValues(1, conn.CreateTxId)
	s.Require().EqualValues(2, conn.ConnectionTxId)
	s.Require().EqualValues(1, conn.ChannelsCount)
	s.Require().NotNil(conn.ConnectionTx)
	s.Require().NotNil(conn.CreateTx)
	s.Require().NotNil(conn.Client)

	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, conn.CreateTx.Hash)

	txHash2, err := hex.DecodeString("652452A670011D629CC116E510BA88C1CABE061336661B1F3D206D248BD55811")
	s.Require().NoError(err)
	s.Require().Equal(txHash2, conn.ConnectionTx.Hash)

	s.Require().EqualValues("osmosis-1", conn.Client.ChainId)
	s.Require().EqualValues("client", conn.Client.Type)
	s.Require().EqualValues(1, conn.Client.ConnectionCount)
}

func (s *StorageTestSuite) TestIbcConnectionList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, fltrs := range []storage.ListConnectionFilters{
		{
			Limit:  10,
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
		}, {
			Limit:    1,
			Offset:   0,
			Sort:     sdk.SortOrderDesc,
			ClientId: "client-1",
		},
	} {

		conns, err := s.storage.IbcConnections.List(ctx, fltrs)
		s.Require().NoError(err)
		s.Require().Len(conns, 1)

		conn := conns[0]
		s.Require().EqualValues("connection-1", conn.ConnectionId)
		s.Require().EqualValues("client-1", conn.ClientId)
		s.Require().EqualValues("counterparty-client-1", conn.CounterpartyClientId)
		s.Require().EqualValues("counterparty-1", conn.CounterpartyConnectionId)
		s.Require().EqualValues(1000, conn.Height)
		s.Require().EqualValues(1, conn.CreateTxId)
		s.Require().EqualValues(2, conn.ConnectionTxId)
		s.Require().EqualValues(1, conn.ChannelsCount)
		s.Require().NotNil(conn.ConnectionTx)
		s.Require().NotNil(conn.CreateTx)
		s.Require().NotNil(conn.Client)

		txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
		s.Require().NoError(err)
		s.Require().Equal(txHash, conn.CreateTx.Hash)

		txHash2, err := hex.DecodeString("652452A670011D629CC116E510BA88C1CABE061336661B1F3D206D248BD55811")
		s.Require().NoError(err)
		s.Require().Equal(txHash2, conn.ConnectionTx.Hash)

		s.Require().EqualValues("osmosis-1", conn.Client.ChainId)
		s.Require().EqualValues("client", conn.Client.Type)
		s.Require().EqualValues(1, conn.Client.ConnectionCount)
	}
}
