// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"
)

func (s *StorageTestSuite) TestSearchText() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.storage.Search.SearchText(ctx, "con")
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues(1, result.Id)
	s.Require().EqualValues("validator", result.Type)
	s.Require().EqualValues("Conqueror", result.Value)
}

func (s *StorageTestSuite) TestSearchTextNamespace() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.storage.Search.SearchText(ctx, "5F7A8DDFE6136FE76B65B9066D4F816D707F")
	s.Require().NoError(err)
	s.Require().Len(results, 2)

	result := results[0]
	s.Require().EqualValues(1, result.Id)
	s.Require().EqualValues("namespace", result.Type)
	s.Require().EqualValues("5f7a8ddfe6136fe76b65b9066d4f816d707f", result.Value)
}

func (s *StorageTestSuite) TestSearch() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("5F7A8DDFE6136FE76B65B9066D4F816D707F28C05B3362D66084664C5B39BA98")
	s.Require().NoError(err)

	results, err := s.storage.Search.Search(ctx, hash)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues(1, result.Id)
	s.Require().EqualValues("block", result.Type)
	s.Require().EqualValues("5f7a8ddfe6136fe76b65b9066d4f816d707f28c05b3362d66084664c5b39ba98", result.Value)
}

func (s *StorageTestSuite) TestSearchByDataHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("3D96B7D238E7E0456F6AF8E7CDF0A67BD6CF9C2089ECB559C659DCAA1F880353")
	s.Require().NoError(err)

	results, err := s.storage.Search.Search(ctx, hash)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues(1, result.Id)
	s.Require().EqualValues("block", result.Type)
	s.Require().EqualValues("3d96b7d238e7e0456f6af8e7cdf0a67bd6cf9c2089ecb559c659dcaa1f880353", result.Value)
}
