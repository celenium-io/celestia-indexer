// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestApiKeyValid() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	key, err := s.storage.ApiKeys.Get(ctx, "test_key")
	s.Require().NoError(err)
	s.Require().EqualValues("test_key", key.Key)
	s.Require().EqualValues("valid key", key.Description)
}

func (s *StorageTestSuite) TestApiKeyInvalid() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.ApiKeys.Get(ctx, "invalid")
	s.Require().Error(err)
}
