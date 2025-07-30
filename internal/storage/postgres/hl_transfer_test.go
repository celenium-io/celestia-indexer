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

func (s *StorageTestSuite) TestHyperlaneTransferList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, fltrs := range []storage.ListHyperlaneTransferFilters{
		{
			Limit:  1,
			Offset: 0,
			Sort:   sdk.SortOrderAsc,
		}, {
			Limit:     1,
			Offset:    0,
			Sort:      sdk.SortOrderDesc,
			AddressId: 1,
		}, {
			Limit:     1,
			Offset:    0,
			Sort:      sdk.SortOrderDesc,
			RelayerId: 2,
		}, {
			Limit:     1,
			Offset:    0,
			Sort:      sdk.SortOrderDesc,
			MailboxId: 1,
		}, {
			Limit:   1,
			Offset:  0,
			Sort:    sdk.SortOrderDesc,
			TokenId: 1,
		},
		{
			Limit:  1,
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
			Type:   types.HLTransferTypeValues(),
		},
		{
			Limit:  1,
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
			Domain: 1234,
		},
	} {

		transfers, err := s.storage.HLTransfer.List(ctx, fltrs)
		s.Require().NoError(err)
		s.Require().Len(transfers, 1)

		transfer := transfers[0]
		s.Require().EqualValues(1, transfer.Id)
		s.Require().EqualValues(1000, transfer.Height)
		s.Require().EqualValues(1234, transfer.Counterparty)
		s.Require().EqualValues("utia", transfer.Denom)
		s.Require().EqualValues(1, transfer.TokenId)
		s.Require().EqualValues(1, transfer.MailboxId)
		s.Require().EqualValues(1, transfer.AddressId)
		s.Require().EqualValues(2, transfer.RelayerId)
		s.Require().EqualValues(1, transfer.Version)
		s.Require().EqualValues(1, transfer.Nonce)
		s.Require().EqualValues("1000", transfer.Amount.String())
		s.Require().EqualValues("1234567890abcdef", transfer.CounterpartyAddress)
		s.Require().EqualValues(types.HLTransferTypeSend, transfer.Type)
		s.Require().NotNil(transfer.Body)
		s.Require().NotNil(transfer.Metadata)

		s.Require().NotNil(transfer.Tx)
		txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
		s.Require().NoError(err)
		s.Require().Equal(txHash, transfer.Tx.Hash)

		s.Require().NotNil(transfer.Address)
		s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", transfer.Address.Address)
		s.Require().NotNil(transfer.Address.Celestials)

		s.Require().NotNil(transfer.Relayer)
		s.Require().EqualValues("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", transfer.Relayer.Address)
		s.Require().NotNil(transfer.Relayer.Celestials)

		s.Require().NotNil(transfer.Mailbox)
		s.Require().Equal([]byte("mailbox"), transfer.Mailbox.Mailbox)

		s.Require().NotNil(transfer.Token)
		s.Require().Equal([]byte("token"), transfer.Token.TokenId)
	}
}

func (s *StorageTestSuite) TestHyperlaneTransferById() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	transfer, err := s.storage.HLTransfer.ById(ctx, 1)
	s.Require().NoError(err)

	s.Require().EqualValues(1, transfer.Id)
	s.Require().EqualValues("utia", transfer.Denom)
	s.Require().EqualValues(1000, transfer.Height)
	s.Require().EqualValues(1234, transfer.Counterparty)

	s.Require().NotNil(transfer.Tx)
}

func (s *StorageTestSuite) TestHyperlaneTransferByIdNotFound() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.HLTransfer.ById(ctx, 100000)
	s.Require().Error(err)
}
