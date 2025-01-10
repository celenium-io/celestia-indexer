package postgres

import (
	"context"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"time"
)

func (s *StorageTestSuite) TestLastSyncTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tvl, err := s.storage.Tvl.LastSyncTime(ctx)
	s.Require().NoError(err)

	s.Require().EqualValues(time.Date(2024, 12, 25, 3, 0, 0, 0, time.UTC), tvl.UTC())
}

func (s *TransactionTestSuite) TestSaveBulk() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	var saveTime = time.Now().UTC()
	var newTvl = &storage.Tvl{
		Value:    12345,
		RollupId: 1,
		Time:     saveTime}

	err := s.storage.Tvl.SaveBulk(ctx, newTvl)
	s.Require().NoError(err)

	tvl, err := s.storage.Tvl.LastSyncTime(ctx)
	s.Require().NoError(err)
	s.Require().True(saveTime.Equal(tvl.UTC()))
}
