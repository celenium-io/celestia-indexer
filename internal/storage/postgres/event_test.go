package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
)

func (s *StorageTestSuite) TestEventByTxId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	events, err := s.storage.Event.ByTxId(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(events, 1)
	s.Require().EqualValues(2, events[0].Id)
	s.Require().EqualValues(1000, events[0].Height)
	s.Require().EqualValues(1, events[0].Position)
	s.Require().Equal(types.EventTypeMint, events[0].Type)
}

func (s *StorageTestSuite) TestEventByBlock() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	events, err := s.storage.Event.ByBlock(ctx, 1000, 2, 0)
	s.Require().NoError(err)
	s.Require().Len(events, 1)
	s.Require().EqualValues(1, events[0].Id)
	s.Require().EqualValues(1000, events[0].Height)
	s.Require().EqualValues(0, events[0].Position)
	s.Require().Equal(types.EventTypeBurn, events[0].Type)
}
