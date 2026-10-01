// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSaveIgps_RenounceOwnership(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	repo := mock.NewMockIHLIGP(ctrl)

	igpId := util.CreateMockHexAddress("igp", 1)
	dCtx := decodeContext.NewContext()
	dCtx.AddIgp(igpId.String(), &storage.HLIGP{IgpId: igpId.Bytes(), Owner: nil})

	tx.EXPECT().
		SaveHyperlaneIgps(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, igps ...*storage.HLIGP) error {
			require.Len(t, igps, 1)
			require.Zero(t, igps[0].OwnerId)
			return nil
		}).
		Times(1)

	// addrToId has no "" key: renounce must not look the owner up
	err := saveIgps(t.Context(), tx, repo, dCtx, map[string]uint64{})
	require.NoError(t, err)
}

func TestSaveIgps_Configs(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	repo := mock.NewMockIHLIGP(ctrl)

	var (
		newIgp      = util.CreateMockHexAddress("igp", 1)
		existingIgp = util.CreateMockHexAddress("igp", 2)
		owner       = "celestia1lg0e9n4pt29lpq2k4ptue4ckw09dx0aujlpe4j"
	)

	dCtx := decodeContext.NewContext()
	dCtx.AddIgp(newIgp.String(), &storage.HLIGP{
		IgpId: newIgp.Bytes(),
		Owner: &storage.Address{Address: owner},
	})
	for _, domain := range []uint64{1, 42161} {
		dCtx.AddIgpConfig(newIgp.String(), &storage.HLIGPConfig{RemoteDomain: domain, GasPrice: types.NumericFromInt64(1)})
		dCtx.AddIgpConfig(existingIgp.String(), &storage.HLIGPConfig{RemoteDomain: domain, GasPrice: types.NumericFromInt64(2)})
	}

	tx.EXPECT().
		SaveHyperlaneIgps(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, igps ...*storage.HLIGP) error {
			require.Len(t, igps, 1)
			require.EqualValues(t, 10, igps[0].OwnerId)
			igps[0].Id = 7 // filled by RETURNING
			return nil
		}).
		Times(1)

	// the IGP saved in this block is resolved without a DB lookup, the other one once for both domains
	repo.EXPECT().
		IdByHash(gomock.Any(), existingIgp.Bytes()).
		Return(uint64(3), nil).
		Times(1)

	tx.EXPECT().
		SaveHyperlaneIgpConfigs(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, configs ...storage.HLIGPConfig) error {
			require.Len(t, configs, 4)

			got := make(map[uint64][]uint64)
			for i := range configs {
				got[configs[i].Id] = append(got[configs[i].Id], configs[i].RemoteDomain)
			}
			require.ElementsMatch(t, []uint64{1, 42161}, got[7])
			require.ElementsMatch(t, []uint64{1, 42161}, got[3])
			return nil
		}).
		Times(1)

	err := saveIgps(t.Context(), tx, repo, dCtx, map[string]uint64{owner: 10})
	require.NoError(t, err)
}
