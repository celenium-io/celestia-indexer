// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestBlobLogsByNamespace() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logs, err := s.storage.BlobLogs.ByNamespace(ctx, 2, storage.BlobLogFilters{
		Limit:  2,
		Offset: 0,
		Sort:   sdk.SortOrderAsc,
		SortBy: "size",
		Joins:  true,
	})
	s.Require().NoError(err)
	s.Require().Len(logs, 2)

	log := logs[0]
	s.Require().EqualValues(1, log.Id)
	s.Require().EqualValues(0, log.Height)
	s.Require().EqualValues("RWW7eaKKXasSGK/DS8PlpErARbl5iFs1vQIycYEAlk0=", log.Commitment)
	s.Require().EqualValues(10, log.Size)
	s.Require().EqualValues(2, log.NamespaceId)
	s.Require().EqualValues(1, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(4, log.TxId)

	s.Require().NotNil(log.Signer)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", log.Signer.Address)
	s.Require().NotNil(log.Signer.Celestials)

	s.Require().NotNil(log.Tx)
	s.Require().EqualValues(4, log.Tx.Id)

	s.Require().NotNil(log.Rollup)
	s.Require().Equal("Rollup 1", log.Rollup.Name)
}

func (s *StorageTestSuite) TestBlobLogsByNamespaceWithoutJoins() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logs, err := s.storage.BlobLogs.ByNamespace(ctx, 2, storage.BlobLogFilters{
		Limit:  2,
		Offset: 0,
		Sort:   sdk.SortOrderAsc,
		SortBy: "size",
		Joins:  false,
	})
	s.Require().NoError(err)
	s.Require().Len(logs, 2)

	log := logs[0]
	s.Require().EqualValues(1, log.Id)
	s.Require().EqualValues(0, log.Height)
	s.Require().EqualValues("RWW7eaKKXasSGK/DS8PlpErARbl5iFs1vQIycYEAlk0=", log.Commitment)
	s.Require().EqualValues(10, log.Size)
	s.Require().EqualValues(2, log.NamespaceId)
	s.Require().EqualValues(1, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(4, log.TxId)

	s.Require().Nil(log.Signer)
	s.Require().Nil(log.Tx)
	s.Require().Nil(log.Rollup)
}

func (s *StorageTestSuite) TestBlobLogsByNamespaceAndCommitment() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logs, err := s.storage.BlobLogs.ByNamespace(ctx, 2, storage.BlobLogFilters{
		Limit:      2,
		Offset:     0,
		Sort:       sdk.SortOrderAsc,
		SortBy:     "size",
		Commitment: "0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw=",
		Joins:      true,
	})
	s.Require().NoError(err)
	s.Require().Len(logs, 1)

	log := logs[0]
	s.Require().EqualValues(3, log.Id)
	s.Require().EqualValues(1000, log.Height)
	s.Require().EqualValues("0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw=", log.Commitment)
	s.Require().EqualValues(12, log.Size)
	s.Require().EqualValues(2, log.NamespaceId)
	s.Require().EqualValues(2, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(2, log.TxId)

	s.Require().NotNil(log.Signer)
	s.Require().EqualValues("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", log.Signer.Address)

	s.Require().NotNil(log.Tx)
}

func (s *StorageTestSuite) TestBlobLogsByNamespaceAndTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logs, err := s.storage.BlobLogs.ByNamespace(ctx, 2, storage.BlobLogFilters{
		Limit:  1,
		Offset: 0,
		Sort:   sdk.SortOrderAsc,
		SortBy: "time",
		From:   time.Date(2023, 7, 3, 11, 0, 0, 0, time.UTC),
		Joins:  true,
	})
	s.Require().NoError(err)
	s.Require().Len(logs, 1)

	log := logs[0]
	s.Require().EqualValues(3, log.Id)
	s.Require().EqualValues(1000, log.Height)
	s.Require().EqualValues("0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw=", log.Commitment)
	s.Require().EqualValues(12, log.Size)
	s.Require().EqualValues(2, log.NamespaceId)
	s.Require().EqualValues(2, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(2, log.TxId)

	s.Require().NotNil(log.Signer)
	s.Require().EqualValues("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", log.Signer.Address)

	s.Require().NotNil(log.Tx)
}

func (s *StorageTestSuite) TestBlobLogsByNamespaceAndToTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logs, err := s.storage.BlobLogs.ByNamespace(ctx, 2, storage.BlobLogFilters{
		Limit:  1,
		Offset: 0,
		Sort:   sdk.SortOrderDesc,
		SortBy: "time",
		To:     time.Date(2023, 7, 4, 13, 0, 0, 0, time.UTC),
		Joins:  true,
	})
	s.Require().NoError(err)
	s.Require().Len(logs, 1)

	log := logs[0]
	s.Require().EqualValues(3, log.Id)
	s.Require().EqualValues(1000, log.Height)
	s.Require().EqualValues("0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw=", log.Commitment)
	s.Require().EqualValues(12, log.Size)
	s.Require().EqualValues(2, log.NamespaceId)
	s.Require().EqualValues(2, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(2, log.TxId)

	s.Require().NotNil(log.Signer)
	s.Require().EqualValues("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", log.Signer.Address)

	s.Require().NotNil(log.Tx)
}

func (s *StorageTestSuite) TestBlobLogsSigner() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logs, err := s.storage.BlobLogs.BySigner(ctx, 1, storage.BlobLogFilters{
		Limit:  2,
		Offset: 0,
		Sort:   sdk.SortOrderAsc,
		SortBy: "size",
		Joins:  true,
	})
	s.Require().NoError(err)
	s.Require().Len(logs, 2)

	log := logs[0]
	s.Require().EqualValues(1, log.Id)
	s.Require().EqualValues(0, log.Height)
	s.Require().EqualValues("RWW7eaKKXasSGK/DS8PlpErARbl5iFs1vQIycYEAlk0=", log.Commitment)
	s.Require().EqualValues(10, log.Size)
	s.Require().EqualValues(2, log.NamespaceId)
	s.Require().EqualValues(1, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(4, log.TxId)
	s.Require().NotNil(log.Namespace)
	s.Require().NotNil(log.Tx)
}

func (s *StorageTestSuite) TestBlobLogsSignerWithoutJoins() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logs, err := s.storage.BlobLogs.BySigner(ctx, 1, storage.BlobLogFilters{
		Limit:  2,
		Offset: 0,
		Sort:   sdk.SortOrderAsc,
		SortBy: "size",
		Joins:  false,
	})
	s.Require().NoError(err)
	s.Require().Len(logs, 2)

	log := logs[0]
	s.Require().EqualValues(1, log.Id)
	s.Require().EqualValues(0, log.Height)
	s.Require().EqualValues("RWW7eaKKXasSGK/DS8PlpErARbl5iFs1vQIycYEAlk0=", log.Commitment)
	s.Require().EqualValues(10, log.Size)
	s.Require().EqualValues(2, log.NamespaceId)
	s.Require().EqualValues(1, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(4, log.TxId)
	s.Require().NotNil(log.Namespace)
	s.Require().Nil(log.Tx)
}

func (s *StorageTestSuite) TestBlobLogsTx() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logs, err := s.storage.BlobLogs.ByTxId(ctx, 4, storage.BlobLogFilters{
		Limit:  2,
		Offset: 0,
		Sort:   sdk.SortOrderAsc,
		SortBy: "size",
	})
	s.Require().NoError(err)
	s.Require().Len(logs, 2)

	log := logs[0]
	s.Require().EqualValues(1, log.Id)
	s.Require().EqualValues(0, log.Height)
	s.Require().EqualValues("RWW7eaKKXasSGK/DS8PlpErARbl5iFs1vQIycYEAlk0=", log.Commitment)
	s.Require().EqualValues(10, log.Size)
	s.Require().EqualValues(2, log.NamespaceId)
	s.Require().EqualValues(1, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(4, log.TxId)
	s.Require().NotNil(log.Namespace)
	s.Require().NotNil(log.TxId)
	s.Require().NotNil(log.Signer)

	s.Require().NotNil(log.Rollup)
	s.Require().Equal("Rollup 1", log.Rollup.Name)
}

func (s *StorageTestSuite) TestCountBlobLogsTx() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	count, err := s.storage.BlobLogs.CountByTxId(ctx, 4)
	s.Require().NoError(err)
	s.Require().EqualValues(count, 2)
}

func (s *StorageTestSuite) TestBlobLogsByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logs, err := s.storage.BlobLogs.ByHeight(ctx, 1000, storage.BlobLogFilters{
		Limit:  2,
		Offset: 0,
		Sort:   sdk.SortOrderDesc,
		SortBy: "size",
	})
	s.Require().NoError(err)
	s.Require().Len(logs, 2)

	log := logs[0]
	s.Require().EqualValues(2, log.Id)
	s.Require().EqualValues(1000, log.Height)
	s.Require().EqualValues("RWW7eaKKXasSGK/DS8PlpErARbl5iFs1vQIycYEAlk0=", log.Commitment)
	s.Require().EqualValues(20, log.Size)
	s.Require().EqualValues(1, log.NamespaceId)
	s.Require().EqualValues(1, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(3, log.TxId)
	s.Require().NotNil(log.Namespace)
	s.Require().NotNil(log.Tx)
	s.Require().NotNil(log.Signer)

	s.Require().NotNil(log.Rollup)
	s.Require().Equal("Rollup 1", log.Rollup.Name)
}

func (s *StorageTestSuite) TestBlobLogsByProviders() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, sortBy := range []string{"", "time", "size"} {
		logs, err := s.storage.BlobLogs.ByProviders(ctx, []storage.RollupProvider{
			{
				AddressId:   1,
				NamespaceId: 1,
			},
		}, storage.BlobLogFilters{
			Limit:  10,
			SortBy: sortBy,
			Joins:  true,
		})
		s.Require().NoError(err)
		s.Require().Len(logs, 1)

		log := logs[0]
		s.Require().NotNil(log.Tx)
		s.Require().NotNil(log.Namespace)
		s.Require().NotNil(log.Signer)
	}
}

func (s *StorageTestSuite) TestBlobLogsByProvidersWithoutJoins() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, sortBy := range []string{"", "time", "size"} {
		logs, err := s.storage.BlobLogs.ByProviders(ctx, []storage.RollupProvider{
			{
				AddressId:   1,
				NamespaceId: 1,
			},
		}, storage.BlobLogFilters{
			Limit:  10,
			SortBy: sortBy,
			Joins:  false,
		})
		s.Require().NoError(err)
		s.Require().Len(logs, 1)

		log := logs[0]
		s.Require().Nil(log.Tx)
		s.Require().NotNil(log.Namespace)
		s.Require().Nil(log.Signer)
	}
}

func (s *StorageTestSuite) TestBlobLogsExportByProviders() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	buf := new(bytes.Buffer)

	from := time.Date(2023, 7, 1, 3, 10, 0, 0, time.UTC)
	to := time.Date(2023, 7, 5, 3, 10, 0, 0, time.UTC)
	err := s.storage.BlobLogs.ExportByProviders(ctx, []storage.RollupProvider{
		{
			AddressId:   1,
			NamespaceId: 1,
		},
	}, from, to, buf)
	s.Require().NoError(err)

	reader := csv.NewReader(buf)

	var count int
	for columns, err := reader.Read(); err != io.EOF; columns, err = reader.Read() {
		s.Require().NoError(err)
		s.Require().Len(columns, 9)
		count++
	}
	s.Require().EqualValues(2, count)
}

func (s *StorageTestSuite) TestBlob() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	log, err := s.storage.BlobLogs.Blob(ctx, 0, 2, "RWW7eaKKXasSGK/DS8PlpErARbl5iFs1vQIycYEAlk0=")
	s.Require().NoError(err)
	s.Require().EqualValues(1, log.Id)
	s.Require().EqualValues(0, log.Height)
	s.Require().EqualValues("RWW7eaKKXasSGK/DS8PlpErARbl5iFs1vQIycYEAlk0=", log.Commitment)
	s.Require().EqualValues(10, log.Size)
	s.Require().EqualValues(2, log.NamespaceId)
	s.Require().EqualValues(1, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(4, log.TxId)

	s.Require().NotNil(log.Signer)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", log.Signer.Address)
	s.Require().NotNil(log.Signer.Celestials)
	s.Require().EqualValues("name 1", log.Signer.Celestials.Id)

	s.Require().NotNil(log.Tx)
	s.Require().EqualValues(4, log.Tx.Id)
}

func (s *StorageTestSuite) TestListBlob() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logs, err := s.storage.BlobLogs.ListBlobs(ctx, storage.ListBlobLogFilters{
		Namespaces: []uint64{2},
		Signers:    []uint64{1},
		SortBy:     "time",
		Sort:       sdk.SortOrderAsc,
		Limit:      1,
	})
	s.Require().NoError(err)
	s.Require().Len(logs, 1)

	log := logs[0]
	s.Require().EqualValues(1, log.Id)
	s.Require().EqualValues(0, log.Height)
	s.Require().EqualValues("RWW7eaKKXasSGK/DS8PlpErARbl5iFs1vQIycYEAlk0=", log.Commitment)
	s.Require().EqualValues(10, log.Size)
	s.Require().EqualValues(2, log.NamespaceId)
	s.Require().EqualValues(1, log.SignerId)
	s.Require().EqualValues(1, log.MsgId)
	s.Require().EqualValues(4, log.TxId)

	s.Require().NotNil(log.Signer)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", log.Signer.Address)

	s.Require().NotNil(log.Tx)
	s.Require().NotNil(log.Namespace)
}
