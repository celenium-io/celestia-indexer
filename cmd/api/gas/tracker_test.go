// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package gas

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTracker_computeMetrics(t *testing.T) {
	t.Run("compute metrics", func(t *testing.T) {
		tracker := NewTracker(nil, nil, nil, nil)

		tracker.q.Push(info{
			Height: 1,
			Percentiles: []decimal.Decimal{
				decimal.RequireFromString("1"),
				decimal.RequireFromString("2"),
				decimal.RequireFromString("3"),
			},
		})
		tracker.q.Push(info{
			Height: 2,
			Percentiles: []decimal.Decimal{
				decimal.RequireFromString("2"),
				decimal.RequireFromString("3"),
				decimal.RequireFromString("4"),
			},
		})
		tracker.q.Push(info{
			Height: 3,
			Percentiles: []decimal.Decimal{
				decimal.RequireFromString("3"),
				decimal.RequireFromString("4"),
				decimal.RequireFromString("5"),
			},
		})

		err := tracker.computeMetrics()
		require.NoError(t, err)
		state := tracker.State()
		require.Equal(t, "2.000000", state.Slow)
		require.Equal(t, "3.000000", state.Median)
		require.Equal(t, "4.000000", state.Fast)
	})

	t.Run("compute metrics: less than default gas price", func(t *testing.T) {
		tracker := NewTracker(nil, nil, nil, nil)

		tracker.q.Push(info{
			Height:      1,
			Percentiles: []decimal.Decimal{},
		})
		tracker.q.Push(info{
			Height:      2,
			Percentiles: []decimal.Decimal{},
		})
		tracker.q.Push(info{
			Height:      3,
			Percentiles: []decimal.Decimal{},
		})

		err := tracker.computeMetrics()
		require.NoError(t, err)
		state := tracker.State()
		require.Equal(t, "0.002000", state.Slow)
		require.Equal(t, "0.002000", state.Median)
		require.Equal(t, "0.002000", state.Fast)
	})
}

func TestTracker_processBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockITx(ctrl)
	ts := time.Now()

	t.Run("empty block", func(t *testing.T) {
		tracker := NewTracker(nil, nil, tx, nil)
		blockStats := storage.BlockStats{
			Time:         ts,
			Height:       1,
			TxCount:      0,
			GasLimit:     0,
			GasUsed:      0,
			Fee:          decimal.New(0, 1),
			BytesInBlock: 0,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err := tracker.processBlock(ctx, blockStats)
		require.NoError(t, err)
		require.Len(t, tracker.q.data, 1)
	})

	t.Run("block with transaction", func(t *testing.T) {

		tx.EXPECT().
			Gas(gomock.Any(), types.Level(1), gomock.Any()).
			Return([]storage.Gas{
				{
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("2000"),
					GasPrice:  decimal.RequireFromString("2"),
				},
			}, nil).
			Times(1)

		tracker := NewTracker(nil, nil, tx, nil)
		blockStats := storage.BlockStats{
			Time:         ts,
			Height:       1,
			TxCount:      1,
			GasLimit:     1000,
			GasUsed:      500,
			Fee:          decimal.RequireFromString("2000"),
			BytesInBlock: maxBlockSize - 1000,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err := tracker.processBlock(ctx, blockStats)
		require.NoError(t, err)
		require.Len(t, tracker.q.data, 1)

		item := tracker.q.data[0]
		require.EqualValues(t, "0.50", item.GasUsedRatio.StringFixed(2))
		require.EqualValues(t, 1, item.TxCount)
		require.Len(t, item.Percentiles, 3)
		require.EqualValues(t, "2", item.Percentiles[0].StringFixed(0))
		require.EqualValues(t, "2", item.Percentiles[1].StringFixed(0))
		require.EqualValues(t, "2", item.Percentiles[2].StringFixed(0))
	})

	t.Run("block with 3 transaction", func(t *testing.T) {

		tx.EXPECT().
			Gas(gomock.Any(), types.Level(2), gomock.Any()).
			Return([]storage.Gas{
				{
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("1000"),
					GasPrice:  decimal.RequireFromString("1"),
				}, {
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("2000"),
					GasPrice:  decimal.RequireFromString("2"),
				}, {
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("3000"),
					GasPrice:  decimal.RequireFromString("3"),
				},
			}, nil).
			Times(1)

		tracker := NewTracker(nil, nil, tx, nil)
		blockStats := storage.BlockStats{
			Time:         ts,
			Height:       2,
			TxCount:      3,
			GasLimit:     3000,
			GasUsed:      1500,
			Fee:          decimal.RequireFromString("6000"),
			BytesInBlock: maxBlockSize - 1000,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err := tracker.processBlock(ctx, blockStats)
		require.NoError(t, err)
		require.Len(t, tracker.q.data, 1)

		item := tracker.q.data[0]
		require.EqualValues(t, "0.50", item.GasUsedRatio.StringFixed(2))
		require.EqualValues(t, 3, item.TxCount)
		require.Len(t, item.Percentiles, 3)
		require.EqualValues(t, "1", item.Percentiles[0].StringFixed(0))
		require.EqualValues(t, "2", item.Percentiles[1].StringFixed(0))
		require.EqualValues(t, "3", item.Percentiles[2].StringFixed(0))
	})

	t.Run("empty block with 3 transaction", func(t *testing.T) {

		tx.EXPECT().
			Gas(gomock.Any(), types.Level(2), gomock.Any()).
			Return([]storage.Gas{
				{
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("1000"),
					GasPrice:  decimal.RequireFromString("1"),
				}, {
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("2000"),
					GasPrice:  decimal.RequireFromString("2"),
				}, {
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("3000"),
					GasPrice:  decimal.RequireFromString("3"),
				},
			}, nil).
			Times(1)

		tracker := NewTracker(nil, nil, tx, nil)
		blockStats := storage.BlockStats{
			Height:       2,
			TxCount:      3,
			GasLimit:     3000,
			GasUsed:      1500,
			Fee:          decimal.RequireFromString("6000"),
			BytesInBlock: 1000,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err := tracker.processBlock(ctx, blockStats)
		require.NoError(t, err)
		require.Len(t, tracker.q.data, 1)

		item := tracker.q.data[0]
		require.EqualValues(t, "0.50", item.GasUsedRatio.StringFixed(2))
		require.EqualValues(t, 3, item.TxCount)
		require.Len(t, item.Percentiles, 3)
		require.EqualValues(t, "1", item.Percentiles[0].StringFixed(0))
		require.EqualValues(t, "1", item.Percentiles[1].StringFixed(0))
		require.EqualValues(t, "1", item.Percentiles[2].StringFixed(0))
	})
}

func TestTracker_compute(t *testing.T) {
	tests := []struct {
		name     string
		txs      []storage.Gas
		gasLimit int64
		want     []string
	}{
		{
			name: "test 1",
			txs: []storage.Gas{
				{
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("1000"),
					GasPrice:  decimal.RequireFromString("1"),
				}, {
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("2000"),
					GasPrice:  decimal.RequireFromString("2"),
				}, {
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("3000"),
					GasPrice:  decimal.RequireFromString("3"),
				},
			},
			gasLimit: 3000,
			want:     []string{"1", "2", "3"},
		}, {
			name: "test 2",
			txs: []storage.Gas{
				{
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("1000"),
					GasPrice:  decimal.RequireFromString("1"),
				}, {
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("1000"),
					GasPrice:  decimal.RequireFromString("1"),
				}, {
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("1000"),
					GasPrice:  decimal.RequireFromString("1"),
				}, {
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("1000"),
					GasPrice:  decimal.RequireFromString("1"),
				}, {
					GasWanted: 1000,
					GasUsed:   500,
					Fee:       decimal.RequireFromString("1000"),
					GasPrice:  decimal.RequireFromString("1"),
				},
			},
			gasLimit: 8000000,
			want:     []string{"1", "1", "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewTracker(nil, nil, nil, nil)
			data := info{
				GasUsedRatio: decimal.New(0, 1),
				Percentiles:  make([]decimal.Decimal, 0),
			}
			for range percentiles {
				data.Percentiles = append(data.Percentiles, decimal.New(0, 1))
			}

			tracker.compute(tt.txs, tt.gasLimit, &data)

			require.Len(t, data.Percentiles, 3)
			require.EqualValues(t, tt.want[0], data.Percentiles[0].StringFixed(0))
			require.EqualValues(t, tt.want[1], data.Percentiles[1].StringFixed(0))
			require.EqualValues(t, tt.want[2], data.Percentiles[2].StringFixed(0))
		})
	}
}
