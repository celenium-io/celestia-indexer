// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test_updateState(t *testing.T) {
	type args struct {
		block           *storage.Block
		totalAccounts   int64
		totalNamespaces int64
		totalProposals  int64
		ibcClientsCount int64
		totalValidators int
		state           *storage.State
	}

	now := time.Now()
	after := time.Now().Add(time.Minute)

	tests := []struct {
		name string
		args args
		want storage.State
	}{
		{
			name: "test 1",
			args: args{
				block: &storage.Block{
					Height:  101,
					Time:    after,
					ChainId: "test",
					Stats: storage.BlockStats{
						TxCount:      10,
						EventsCount:  300,
						BlobsSize:    100,
						SupplyChange: decimal.RequireFromString("100"),
						Fee:          decimal.RequireFromString("10"),
					},
				},
				totalAccounts:   10,
				totalNamespaces: 12,
				totalProposals:  8,
				totalValidators: 3,
				ibcClientsCount: 4,
				state: &storage.State{
					Id:              1,
					Name:            "test",
					LastHeight:      100,
					LastTime:        now,
					ChainId:         "chain_id",
					TotalTx:         10,
					TotalAccounts:   2,
					TotalNamespaces: 2,
					TotalBlobsSize:  1,
					TotalSupply:     decimal.RequireFromString("1000"),
					TotalFee:        decimal.RequireFromString("10"),
					TotalValidators: 1,
					TotalIbcClients: 3,
				},
			},
			want: storage.State{
				Id:              1,
				Name:            "test",
				LastHeight:      101,
				LastTime:        after,
				ChainId:         "chain_id",
				TotalTx:         20,
				TotalAccounts:   12,
				TotalNamespaces: 14,
				TotalBlobsSize:  101,
				TotalSupply:     decimal.RequireFromString("1100"),
				TotalFee:        decimal.RequireFromString("20"),
				TotalValidators: 4,
				TotalProposals:  8,
				TotalIbcClients: 7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := updateState(tt.args.block, tt.args.totalAccounts, tt.args.totalNamespaces, tt.args.totalProposals, tt.args.ibcClientsCount, tt.args.totalValidators, decimal.Zero, tt.args.state)
			require.NoError(t, err)
		})
	}
}
