// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

func (s *StorageTestSuite) TestZkISMList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.ZkISM.List(ctx, storage.ZkISMFilter{
		Limit: 10,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 2)

	// Default sort is DESC by time,id — id=2 (newer time) comes first.
	ism := items[0]
	s.Require().EqualValues(2, ism.Id)
	s.Require().EqualValues(1000, ism.Height)
	s.Require().EqualValues(200, ism.ExternalId)
	s.Require().NotEmpty(ism.State)
	s.Require().NotEmpty(ism.StateRoot)
	s.Require().NotEmpty(ism.MerkleTreeAddress)
	s.Require().NotEmpty(ism.Groth16VKey)
	s.Require().NotEmpty(ism.StateTransitionVKey)
	s.Require().NotEmpty(ism.StateMembershipVKey)

	s.Require().NotNil(ism.Tx)
	txHash, err := hex.DecodeString("652452A670011D629CC116E510BA88C1CABE061336661B1F3D206D248BD55811")
	s.Require().NoError(err)
	s.Require().Equal(txHash, ism.Tx.Hash)

	s.Require().NotNil(ism.Creator)
	s.Require().EqualValues("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", ism.Creator.Address)
}

func (s *StorageTestSuite) TestZkISMListByCreator() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	creatorId := uint64(1)
	items, err := s.storage.ZkISM.List(ctx, storage.ZkISMFilter{
		Limit:     10,
		CreatorId: &creatorId,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	ism := items[0]
	s.Require().EqualValues(1, ism.Id)
	s.Require().EqualValues(100, ism.ExternalId)

	s.Require().NotNil(ism.Creator)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", ism.Creator.Address)
}

func (s *StorageTestSuite) TestZkISMListByTx() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txId := uint64(2)
	items, err := s.storage.ZkISM.List(ctx, storage.ZkISMFilter{
		Limit: 10,
		TxId:  &txId,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	ism := items[0]
	s.Require().EqualValues(2, ism.Id)
	s.Require().EqualValues(200, ism.ExternalId)

	s.Require().NotNil(ism.Tx)
	txHash, err := hex.DecodeString("652452A670011D629CC116E510BA88C1CABE061336661B1F3D206D248BD55811")
	s.Require().NoError(err)
	s.Require().Equal(txHash, ism.Tx.Hash)
}

func (s *StorageTestSuite) TestZkISMListOffset() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.ZkISM.List(ctx, storage.ZkISMFilter{
		Limit:  10,
		Offset: 1,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	// With default DESC sort and offset=1, the second (older) item is returned.
	s.Require().EqualValues(1, items[0].Id)
}

func (s *StorageTestSuite) TestZkISMListAsc() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.ZkISM.List(ctx, storage.ZkISMFilter{
		Limit: 10,
		Sort:  sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 2)

	// ASC sort by time,id — id=1 (older time) comes first.
	s.Require().EqualValues(1, items[0].Id)
	s.Require().EqualValues(2, items[1].Id)
}

// ---------------------------------------------------------------------------
// ById
// ---------------------------------------------------------------------------

func (s *StorageTestSuite) TestZkISMById() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	ism, err := s.storage.ZkISM.ById(ctx, 1)
	s.Require().NoError(err)

	s.Require().EqualValues(1, ism.Id)
	s.Require().EqualValues(1000, ism.Height)
	s.Require().EqualValues(100, ism.ExternalId)
	s.Require().NotEmpty(ism.State)
	s.Require().NotEmpty(ism.StateRoot)
	s.Require().NotEmpty(ism.MerkleTreeAddress)
	s.Require().NotEmpty(ism.Groth16VKey)
	s.Require().NotEmpty(ism.StateTransitionVKey)
	s.Require().NotEmpty(ism.StateMembershipVKey)

	s.Require().NotNil(ism.Tx)
	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, ism.Tx.Hash)

	s.Require().NotNil(ism.Creator)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", ism.Creator.Address)
}

func (s *StorageTestSuite) TestZkISMByIdNotFound() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.ZkISM.ById(ctx, 999)
	s.Require().Error(err)
}

// ---------------------------------------------------------------------------
// Updates
// ---------------------------------------------------------------------------

func (s *StorageTestSuite) TestZkISMUpdates() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.ZkISM.Updates(ctx, 1, storage.ZkISMUpdatesFilter{
		Limit: 10,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 2)

	// Default DESC sort — id=2 (newer time) comes first.
	u := items[0]
	s.Require().EqualValues(2, u.Id)
	s.Require().EqualValues(1, u.ZkISMId)
	s.Require().EqualValues(1002, u.Height)
	s.Require().NotEmpty(u.NewState)
	s.Require().NotEmpty(u.NewStateRoot)

	s.Require().NotNil(u.Tx)
	txHash, err := hex.DecodeString("652452A670011D629CC116E510BA88C1CABE061336661B1F3D206D248BD55811")
	s.Require().NoError(err)
	s.Require().Equal(txHash, u.Tx.Hash)

	s.Require().NotNil(u.Signer)
	s.Require().EqualValues("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", u.Signer.Address)
}

func (s *StorageTestSuite) TestZkISMUpdatesBySigner() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	signerId := uint64(1)
	items, err := s.storage.ZkISM.Updates(ctx, 1, storage.ZkISMUpdatesFilter{
		Limit:    10,
		SignerId: &signerId,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	u := items[0]
	s.Require().EqualValues(1, u.Id)
	s.Require().EqualValues(1001, u.Height)

	s.Require().NotNil(u.Signer)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", u.Signer.Address)
}

func (s *StorageTestSuite) TestZkISMUpdatesByTx() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txId := uint64(2)
	items, err := s.storage.ZkISM.Updates(ctx, 1, storage.ZkISMUpdatesFilter{
		Limit: 10,
		TxId:  &txId,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	u := items[0]
	s.Require().EqualValues(2, u.Id)

	s.Require().NotNil(u.Tx)
	txHash, err := hex.DecodeString("652452A670011D629CC116E510BA88C1CABE061336661B1F3D206D248BD55811")
	s.Require().NoError(err)
	s.Require().Equal(txHash, u.Tx.Hash)
}

func (s *StorageTestSuite) TestZkISMUpdatesFrom() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	// Row 1 is at 04:11:57, Row 2 is at 05:12:57.
	// From 05:00:00 should include only Row 2.
	from := time.Date(2023, 7, 4, 5, 0, 0, 0, time.UTC)
	items, err := s.storage.ZkISM.Updates(ctx, 1, storage.ZkISMUpdatesFilter{
		Limit: 10,
		From:  from,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)
	s.Require().EqualValues(2, items[0].Id)
}

func (s *StorageTestSuite) TestZkISMUpdatesTo() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	// To 05:00:00 should include only Row 1 (04:11:57 < 05:00:00).
	to := time.Date(2023, 7, 4, 5, 0, 0, 0, time.UTC)
	items, err := s.storage.ZkISM.Updates(ctx, 1, storage.ZkISMUpdatesFilter{
		Limit: 10,
		To:    to,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)
	s.Require().EqualValues(1, items[0].Id)
}

func (s *StorageTestSuite) TestZkISMUpdatesUnknownISM() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.ZkISM.Updates(ctx, 999, storage.ZkISMUpdatesFilter{
		Limit: 10,
	})
	s.Require().NoError(err)
	s.Require().Empty(items)
}

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------

func (s *StorageTestSuite) TestZkISMMessages() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.ZkISM.Messages(ctx, 1, storage.ZkISMUpdatesFilter{
		Limit: 10,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 2)

	// Default DESC sort — id=2 (newer time) comes first.
	m := items[0]
	s.Require().EqualValues(2, m.Id)
	s.Require().EqualValues(1, m.ZkISMId)
	s.Require().EqualValues(1002, m.Height)
	s.Require().NotEmpty(m.StateRoot)
	s.Require().NotEmpty(m.MessageId)

	s.Require().NotNil(m.Tx)
	txHash, err := hex.DecodeString("652452A670011D629CC116E510BA88C1CABE061336661B1F3D206D248BD55811")
	s.Require().NoError(err)
	s.Require().Equal(txHash, m.Tx.Hash)

	s.Require().NotNil(m.Signer)
	s.Require().EqualValues("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", m.Signer.Address)
}

func (s *StorageTestSuite) TestZkISMMessagesBySigner() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	signerId := uint64(1)
	items, err := s.storage.ZkISM.Messages(ctx, 1, storage.ZkISMUpdatesFilter{
		Limit:    10,
		SignerId: &signerId,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	m := items[0]
	s.Require().EqualValues(1, m.Id)
	s.Require().EqualValues(1001, m.Height)

	s.Require().NotNil(m.Signer)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", m.Signer.Address)
}

func (s *StorageTestSuite) TestZkISMMessagesByTx() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txId := uint64(1)
	items, err := s.storage.ZkISM.Messages(ctx, 1, storage.ZkISMUpdatesFilter{
		Limit: 10,
		TxId:  &txId,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	m := items[0]
	s.Require().EqualValues(1, m.Id)

	s.Require().NotNil(m.Tx)
	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, m.Tx.Hash)
}

func (s *StorageTestSuite) TestZkISMMessagesFrom() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	from := time.Date(2023, 7, 4, 5, 0, 0, 0, time.UTC)
	items, err := s.storage.ZkISM.Messages(ctx, 1, storage.ZkISMUpdatesFilter{
		Limit: 10,
		From:  from,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)
	s.Require().EqualValues(2, items[0].Id)
}

func (s *StorageTestSuite) TestZkISMMessagesUnknownISM() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.ZkISM.Messages(ctx, 999, storage.ZkISMUpdatesFilter{
		Limit: 10,
	})
	s.Require().NoError(err)
	s.Require().Empty(items)
}
