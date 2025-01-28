// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"slices"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func (s *StorageTestSuite) TestCelestialsById() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	item, err := s.storage.Celestials.ById(ctx, "name 3")
	s.Require().NoError(err)
	s.Require().EqualValues("", item.ImageUrl)
	s.Require().EqualValues("name 3", item.Id)
	s.Require().EqualValues(3, item.ChangeId)
	s.Require().EqualValues(2, item.AddressId)
}

func (s *StorageTestSuite) TestCelestialsByAddressId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Celestials.ByAddressId(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(items, 2)

	item := items[0]
	s.Require().EqualValues("", item.ImageUrl)
	s.Require().EqualValues("name 2", item.Id)
	s.Require().EqualValues(2, item.ChangeId)
	s.Require().EqualValues(1, item.AddressId)
}

func (s *StorageTestSuite) TestCelestialsTransaction() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginCelestialTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	state, err := s.storage.CelestialState.ByName(ctx, "indexer")
	s.Require().NoError(err)

	celIds := []storage.Celestial{
		{
			Id:        "name 3",
			AddressId: 1,
			ChangeId:  4,
			ImageUrl:  "image_url",
		}, {
			Id:        "name 4",
			AddressId: 3,
			ChangeId:  5,
			ImageUrl:  "image_url2",
		},
	}
	state.ChangeId = celIds[1].ChangeId

	err = tx.SaveCelestials(ctx, slices.Values(celIds))
	s.Require().NoError(err)

	err = tx.UpdateState(ctx, &state)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	state1, err := s.storage.CelestialState.ByName(ctx, "indexer")
	s.Require().NoError(err)
	s.Require().EqualValues(celIds[1].ChangeId, state1.ChangeId)

	item, err := s.storage.Celestials.ById(ctx, "name 3")
	s.Require().NoError(err)
	s.Require().EqualValues("image_url", item.ImageUrl)
	s.Require().EqualValues("name 3", item.Id)
	s.Require().EqualValues(4, item.ChangeId)
	s.Require().EqualValues(1, item.AddressId)

	item2, err := s.storage.Celestials.ById(ctx, "name 4")
	s.Require().NoError(err)
	s.Require().EqualValues("image_url2", item2.ImageUrl)
	s.Require().EqualValues("name 4", item2.Id)
	s.Require().EqualValues(5, item2.ChangeId)
	s.Require().EqualValues(3, item2.AddressId)
}

func (s *StorageTestSuite) TestCelestialsStateSave() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	err := s.storage.CelestialState.Save(ctx, &storage.CelestialState{
		Name:     "new",
		ChangeId: 10,
	})
	s.Require().NoError(err)

	state, err := s.storage.CelestialState.ByName(ctx, "new")
	s.Require().NoError(err)
	s.Require().EqualValues(10, state.ChangeId)
	s.Require().EqualValues("new", state.Name)
}

func (s *StorageTestSuite) TestCelestialsStateByName() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	state, err := s.storage.CelestialState.ByName(ctx, "indexer")
	s.Require().NoError(err)
	s.Require().EqualValues(3, state.ChangeId)
	s.Require().EqualValues("indexer", state.Name)
}
