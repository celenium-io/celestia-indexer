// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestNewBlock(t *testing.T) {
	type args struct {
		block     storage.Block
		withStats bool
	}
	tests := []struct {
		name string
		args args
		want Block
	}{
		{
			name: "without stats",
			args: args{
				block: storage.Block{
					Id:           1000,
					Height:       1000,
					Time:         time.Time{},
					VersionBlock: 10,
					VersionApp:   11,
					MessageTypes: storageTypes.MsgTypeBits{
						Bits: storageTypes.NewBits(1),
					},
					Hash:               []byte{0x01},
					ParentHash:         []byte{0x02},
					LastCommitHash:     []byte{0x03},
					DataHash:           []byte{0x04},
					ValidatorsHash:     []byte{0x05},
					NextValidatorsHash: []byte{0x06},
					ConsensusHash:      []byte{0x07},
					AppHash:            []byte{0x08},
					LastResultsHash:    []byte{0x09},
					EvidenceHash:       []byte{0x10},
					ChainId:            "dipdup",
					Proposer: storage.Validator{
						Id:          1,
						Moniker:     "moniker",
						ConsAddress: "001122",
					},
				},
				withStats: false,
			},
			want: Block{
				Id:                 1000,
				Height:             1000,
				Time:               time.Time{},
				VersionBlock:       "10",
				VersionApp:         "11",
				Hash:               []byte{0x01},
				ParentHash:         []byte{0x02},
				LastCommitHash:     []byte{0x03},
				DataHash:           []byte{0x04},
				ValidatorsHash:     []byte{0x05},
				NextValidatorsHash: []byte{0x06},
				ConsensusHash:      []byte{0x07},
				AppHash:            []byte{0x08},
				LastResultsHash:    []byte{0x09},
				EvidenceHash:       []byte{0x10},
				MessageTypes:       []storageTypes.MsgType{storageTypes.MsgUnknown},
				Proposer: &ShortValidator{
					Id:          1,
					Moniker:     "moniker",
					ConsAddress: "001122",
				},
				Stats: nil,
			},
		},

		{
			name: "with stats",
			args: args{
				block: storage.Block{
					Id:           1000,
					Height:       1000,
					Time:         time.Time{},
					VersionBlock: 10,
					VersionApp:   11,
					MessageTypes: storageTypes.MsgTypeBits{
						Bits: storageTypes.NewBits(1),
					},
					Hash:               []byte{0x01},
					ParentHash:         []byte{0x02},
					LastCommitHash:     []byte{0x03},
					DataHash:           []byte{0x04},
					ValidatorsHash:     []byte{0x05},
					NextValidatorsHash: []byte{0x06},
					ConsensusHash:      []byte{0x07},
					AppHash:            []byte{0x08},
					LastResultsHash:    []byte{0x09},
					EvidenceHash:       []byte{0x10},
					ChainId:            "dipdup",
					Proposer: storage.Validator{
						Id:          1,
						Moniker:     "moniker",
						ConsAddress: "001122",
					},
					Stats: storage.BlockStats{
						Id:            1000,
						Height:        1000,
						Time:          time.Time{},
						TxCount:       6,
						EventsCount:   10,
						BlobsSize:     1234,
						BlockTime:     11000,
						SupplyChange:  decimal.NewFromInt(123),
						InflationRate: decimal.NewFromFloat(0.08),
						Fee:           decimal.NewFromInt(125),
						BytesInBlock:  10000,
					},
				},
				withStats: true,
			},
			want: Block{
				Id:                 1000,
				Height:             1000,
				Time:               time.Time{},
				VersionBlock:       "10",
				VersionApp:         "11",
				Hash:               []byte{0x01},
				ParentHash:         []byte{0x02},
				LastCommitHash:     []byte{0x03},
				DataHash:           []byte{0x04},
				ValidatorsHash:     []byte{0x05},
				NextValidatorsHash: []byte{0x06},
				ConsensusHash:      []byte{0x07},
				AppHash:            []byte{0x08},
				LastResultsHash:    []byte{0x09},
				EvidenceHash:       []byte{0x10},
				MessageTypes:       []storageTypes.MsgType{storageTypes.MsgUnknown},
				Proposer: &ShortValidator{
					Id:          1,
					Moniker:     "moniker",
					ConsAddress: "001122",
				},
				Stats: &BlockStats{
					TxCount:       6,
					EventsCount:   10,
					BlobsSize:     1234,
					Fee:           "125",
					SupplyChange:  "123",
					InflationRate: "0.08",
					Rewards:       "0",
					Commissions:   "0",
					BlockTime:     11000,
					BytesInBlock:  10000,
					FillRate:      "0.0051",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBlock(tt.args.block, tt.args.withStats)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNewBlockStats(t *testing.T) {
	type args struct {
		stats storage.BlockStats
	}
	tests := []struct {
		name string
		args args
		want *BlockStats
	}{
		{
			name: "with msgs counts",
			args: args{
				stats: storage.BlockStats{
					Id:            100,
					Height:        100,
					Time:          time.Time{},
					TxCount:       6,
					EventsCount:   10,
					BlobsSize:     1234,
					BlockTime:     11000,
					SupplyChange:  decimal.NewFromInt(123),
					InflationRate: decimal.NewFromFloat(0.08),
					Fee:           decimal.NewFromInt(125),
					BytesInBlock:  10000,
				},
			},
			want: &BlockStats{
				TxCount:       6,
				EventsCount:   10,
				BlobsSize:     1234,
				Fee:           "125",
				SupplyChange:  "123",
				InflationRate: "0.08",
				Rewards:       "0",
				Commissions:   "0",
				BlockTime:     11000,
				BytesInBlock:  10000,
				FillRate:      "0.0051",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBlockStats(tt.args.stats)
			require.Equal(t, tt.want, got)
		})
	}
}
