// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestIbcClientById() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	client, err := s.storage.IbcClients.ById(ctx, "client-1")
	s.Require().NoError(err)

	s.Require().EqualValues("client-1", client.Id)
	s.Require().EqualValues(1000, client.Height)
	s.Require().EqualValues(1, client.TrustLevelDenominator)
	s.Require().EqualValues("osmosis-1", client.ChainId)
	s.Require().NotNil(client.Tx)
	s.Require().NotNil(client.Creator)

	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, client.Tx.Hash)

	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", client.Creator.Address)
}

func (s *StorageTestSuite) TestIbcClientList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	clients, err := s.storage.IbcClients.List(ctx, 1, 0, sdk.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(clients, 1)

	client := clients[0]
	s.Require().EqualValues("client-1", client.Id)
	s.Require().EqualValues(1000, client.Height)
	s.Require().EqualValues(1, client.TrustLevelDenominator)
	s.Require().EqualValues("osmosis-1", client.ChainId)
	s.Require().NotNil(client.Tx)
	s.Require().NotNil(client.Creator)

	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, client.Tx.Hash)

	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", client.Creator.Address)
}
