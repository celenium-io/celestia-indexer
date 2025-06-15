// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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

func (s *StorageTestSuite) TestHyperlaneTokenByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	token, err := s.storage.HLToken.ByHash(ctx, []byte("token"))
	s.Require().NoError(err)

	s.Require().EqualValues(1, token.Id)
	s.Require().EqualValues(1000, token.Height)
	s.Require().EqualValues(12, token.SentTransfers)
	s.Require().EqualValues(21, token.ReceiveTransfers)
	s.Require().EqualValues("100000", token.Sent.String())
	s.Require().EqualValues("200000", token.Received.String())
	s.Require().EqualValues("utia", token.Denom)
	s.Require().EqualValues([]byte("token"), token.TokenId)
	s.Require().EqualValues(types.HLTokenTypeCollateral, token.Type)

	s.Require().NotNil(token.Tx)
	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, token.Tx.Hash)

	s.Require().NotNil(token.Owner)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", token.Owner.Address)
	s.Require().NotNil(token.Owner.Celestials)

	s.Require().NotNil(token.Mailbox)
	s.Require().Equal([]byte("mailbox"), token.Mailbox.Mailbox)
}

func (s *StorageTestSuite) TestHyperlaneTokenList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tokens, err := s.storage.HLToken.List(ctx, storage.ListHyperlaneTokens{
		Limit:     1,
		Offset:    0,
		Sort:      sdk.SortOrderAsc,
		MailboxId: 1,
		OwnerId:   1,
		Type:      []types.HLTokenType{types.HLTokenTypeCollateral},
	})
	s.Require().NoError(err)
	s.Require().Len(tokens, 1)

	token := tokens[0]
	s.Require().EqualValues(1, token.Id)
	s.Require().EqualValues(1000, token.Height)
	s.Require().EqualValues(12, token.SentTransfers)
	s.Require().EqualValues(21, token.ReceiveTransfers)
	s.Require().EqualValues("100000", token.Sent.String())
	s.Require().EqualValues("200000", token.Received.String())
	s.Require().EqualValues("utia", token.Denom)
	s.Require().EqualValues([]byte("token"), token.TokenId)
	s.Require().EqualValues(types.HLTokenTypeCollateral, token.Type)

	s.Require().NotNil(token.Tx)
	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, token.Tx.Hash)

	s.Require().NotNil(token.Owner)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", token.Owner.Address)
	s.Require().NotNil(token.Owner.Celestials)

	s.Require().NotNil(token.Mailbox)
	s.Require().Equal([]byte("mailbox"), token.Mailbox.Mailbox)
}
