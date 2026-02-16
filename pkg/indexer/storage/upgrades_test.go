package storage

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	indexerCfg "github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpgradeV7(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validators := mock.NewMockIValidator(ctrl)
	validators.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*storage.Validator{
			{
				Id:      1,
				Rate:    decimal.RequireFromString("0.150000000000000000"),
				MaxRate: decimal.RequireFromString("0.500000000000000000"),
			},
			{
				Id:      2,
				Rate:    decimal.RequireFromString("0.250000000000000000"),
				MaxRate: decimal.RequireFromString("0.700000000000000000"),
			},
		}, nil).
		Times(1)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	module := NewModule(nil, nil, validators, nil, indexerCfg.Indexer{Name: testIndexerName})
	dCtx := decodeContext.NewContext()

	err := module.upgradeV7(ctx, dCtx, 7)
	require.NoError(t, err)

	minCommissionRate := decimal.RequireFromString("0.200000000000000000")
	maxCommissionRate := decimal.RequireFromString("0.600000000000000000")

	err = dCtx.Validators.Range(func(key string, value *storage.Validator) (error, bool) {
		require.True(t, value.Rate.GreaterThanOrEqual(minCommissionRate))
		require.True(t, value.MaxRate.LessThanOrEqual(maxCommissionRate))
		return nil, true
	})
	require.NoError(t, err)

	err = dCtx.Constants.Range(func(_ string, value *storage.Constant) (error, bool) {
		if value.Name == "min_commission_rate" {
			require.Equal(t, "0.200000000000000000", value.Value)
		}
		if value.Name == "max_commission_rate" {
			require.Equal(t, "0.600000000000000000", value.Value)
		}
		return nil, true
	})
	require.NoError(t, err)
}
