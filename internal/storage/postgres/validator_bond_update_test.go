// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestValidatorBondUpdateListByValidator() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	updates, err := s.storage.ValidatorBondUpdates.ListByValidator(ctx, 1, storage.FilterBondUpdatesListByValidator{
		Limit: 10,
		Sort:  sdk.SortOrderDesc,
	})
	s.Require().NoError(err)
	s.Require().Len(updates, 2)
	s.Require().EqualValues(2, updates[0].Id)
	s.Require().EqualValues(1000, updates[0].Height)
	s.Require().Equal("1", updates[0].Power.String())
	s.Require().EqualValues(1, updates[1].Id)
	s.Require().Equal("2", updates[1].Power.String())

	for i := range updates {
		s.Require().EqualValues(1, updates[i].ValidatorId)
	}
}

func (s *StorageTestSuite) TestValidatorBondUpdateListByValidatorPagination() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	// asc is the default order
	updates, err := s.storage.ValidatorBondUpdates.ListByValidator(ctx, 1, storage.FilterBondUpdatesListByValidator{
		Limit:  1,
		Offset: 1,
	})
	s.Require().NoError(err)
	s.Require().Len(updates, 1)
	s.Require().EqualValues(2, updates[0].Id)

	updates, err = s.storage.ValidatorBondUpdates.ListByValidator(ctx, 100, storage.FilterBondUpdatesListByValidator{Limit: 10})
	s.Require().NoError(err)
	s.Require().Empty(updates)
}
